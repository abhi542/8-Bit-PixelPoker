package poker

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrNotYourTurn   = errors.New("not your turn")
	ErrInvalidAction = errors.New("invalid action")
	ErrInsufficient  = errors.New("insufficient chips")
	ErrRaiseTooSmall = errors.New("raise must be at least min-raise over previous bet")
)

type Pot struct {
	Amount        int64
	EligibleSeats []int // Indices of players who can share this pot
}

// BettingEngine manages the math and flow of a single hand's betting.
type BettingEngine struct {
	Players []*PlayerBetState

	// Round State
	CurrentHighBet int64 // The amount to match (e.g. BB, or Raise)
	MinRaise       int64 // The minimum additional amount required to raise
	LastScaleRaise int64 // Amount of the last full raise (to calc next MinRaise)
	ActionIndex    int   // Index into Players slice (who is acting)

	// Pots
	MainPot  int64
	SidePots []Pot // Generated at end of rounds or hand

	// State Control
	AggressorIndex int  // Index of player who made the last aggressive action (for reopen checks)
	IsRoundClosed  bool // True when all betting is settled for this street
	DealerIndex    int  // needed to determine start of next round
}

func NewBettingEngine(players []*PlayerBetState, dealerIndex int, blinds int64) *BettingEngine {
	// Simple init. Real game passes specific structure.
	// We assume players slice is 0..N sorted by seat.
	return &BettingEngine{
		Players:        players,
		ActionIndex:    (dealerIndex + 1) % len(players), // Default, effectively SB. Caller should adjust for PreFlop (UTG)
		MinRaise:       blinds,                           // Default min raise is usually Big Blind
		LastScaleRaise: blinds,
		DealerIndex:    dealerIndex,
	}
}

// StartNewRound resets the betting state for a new street (Flop, Turn, River).
func (be *BettingEngine) StartNewRound() {
	be.CurrentHighBet = 0
	be.MinRaise = 20 // Hardcoded to BB for now

	be.IsRoundClosed = false
	be.LastScaleRaise = 0

	// Reset Player Round Bets and Sweep into MainPot
	for _, p := range be.Players {
		be.MainPot += p.BetThisRound
		p.BetThisRound = 0
	}

	// Calculate Next Action (First active after Dealer)
	n := len(be.Players)
	nextIdx := (be.DealerIndex + 1) % n

	// Loop to find active
	found := false
	start := nextIdx
	for {
		p := be.Players[nextIdx]
		if p.Status == StatusActive && !p.IsAllIn {
			found = true
			break
		}
		nextIdx = (nextIdx + 1) % n
		if nextIdx == start {
			break // Everyone is out/all-in?
		}
	}

	if !found {
		be.IsRoundClosed = true
		return
	}

	be.ActionIndex = nextIdx
	be.AggressorIndex = nextIdx // First to act is the stopper for checking around
}

// ApplyAction executes a player's move.
func (be *BettingEngine) ApplyAction(seatIdx int, act Action) error {
	if be.IsRoundClosed {
		return errors.New("betting round is closed")
	}

	// 1. Validate Initiator
	// We scan our linear Players list to find the one matching seatIdx.
	// In production, we might map seatIdx -> slice index, but Loop is fine for N=9.
	// Actually, be.Players usually *is* ordered by seat.
	// Let's assume be.Players[i].SeatIndex == i is NOT guaranteed if seats are empty.
	// So we find the player object.
	var p *PlayerBetState
	var pIdx int
	for i, pl := range be.Players {
		if pl.SeatIndex == seatIdx {
			p = pl
			pIdx = i
			break
		}
	}
	if p == nil {
		return errors.New("player not found")
	}

	// 2. Validate Turn
	// Turn logic: be.ActionIndex points to be.Players slice index.
	if pIdx != be.ActionIndex {
		return ErrNotYourTurn
	}

	// 3. Process Action
	switch act.Type {
	case ActionFold:
		p.Status = StatusFolded

	case ActionCheck:
		if be.CurrentHighBet > p.BetThisRound {
			return errors.New("cannot check, must call or fold")
		}
		// Valid check. No chip movement.

	case ActionCall:
		needed := be.CurrentHighBet - p.BetThisRound
		if needed > p.Chips {
			// All-In Call (Partial)
			if err := be.executeBet(p, pIdx, p.Chips, true); err != nil {
				return err
			}
		} else {
			if err := be.executeBet(p, pIdx, needed, false); err != nil {
				return err
			}
		}

	case ActionBet, ActionRaise, ActionAllIn:
		// Logic:
		// Action.Amount is target TotalBet (e.g. Raise TO 100).
		// Chips to add = Target - AlreadyBet.

		toAdd := act.Amount - p.BetThisRound
		if toAdd < 0 {
			return errors.New("bet amount less than already committed")
		}

		isAllIn := (toAdd == p.Chips)

		// If explicit AllIn action, verify amount matches stack
		if act.Type == ActionAllIn && !isAllIn {
			return errors.New("all-in action amount mismatch")
		}

		// Min Raise Check
		// Raise size = NewTotal - CurrentHighBet
		raiseSize := act.Amount - be.CurrentHighBet

		// Note: A "Bet" when HighBet=0 is effectively a Raise from 0.
		// If HighBet > 0, it's a Raise.

		if raiseSize < be.MinRaise && !isAllIn && act.Amount > be.CurrentHighBet {
			// Allow "Incomplete Raise" purely for All-Ins.
			// If not all-in, reject.
			return fmt.Errorf("%w: need %d more", ErrRaiseTooSmall, be.MinRaise-raiseSize)
		} else if act.Amount < be.CurrentHighBet {
			// This is an invalid under-bet? Or a call?
			// Usually UI handles "Call", so explicit Bet < High shouldn't happen unless malicious.
			return errors.New("bet falls short of current call amount")
		}

		if err := be.executeBet(p, pIdx, toAdd, isAllIn); err != nil {
			return err
		}
	}

	// 4. Advance
	be.advanceTurn()
	return nil
}

// executeBet updates chips and state limits.
func (be *BettingEngine) executeBet(p *PlayerBetState, pIdx int, amount int64, isAllIn bool) error {
	if amount > p.Chips {
		return ErrInsufficient
	}

	p.Chips -= amount
	p.BetThisRound += amount
	p.TotalCommitted += amount
	p.IsAllIn = isAllIn

	// Update High Water Mark
	if p.BetThisRound > be.CurrentHighBet {
		raiseDiff := p.BetThisRound - be.CurrentHighBet

		// Rule: Does this re-open betting?
		// Full Legal Raise: Diff >= MinRaise.
		if raiseDiff >= be.MinRaise {
			be.MinRaise = raiseDiff // Standard NLH: Min raise matches previous raise size
			be.LastScaleRaise = raiseDiff
			be.AggressorIndex = pIdx // Reset loop trigger
		}
		// If Incomplete Raise (All-In < Min), we update HighBet but NOT MinRaise logic usually?
		// Actually HighBet definitely updates.
		be.CurrentHighBet = p.BetThisRound
	}

	return nil
}

// advanceTurn moves ActionIndex to next active player.
func (be *BettingEngine) advanceTurn() {
	start := be.ActionIndex
	n := len(be.Players)

	for {
		be.ActionIndex = (be.ActionIndex + 1) % n

		// Check Round End Condition:
		// If we wrapped around to the Aggressor (or everyone acted once), AND everyone matches bets.
		if be.checkRoundComplete() {
			be.IsRoundClosed = true
			return
		}

		// Found active player?
		p := be.Players[be.ActionIndex]
		if p.Status == StatusActive && !p.IsAllIn {
			return // Found next actor
		}

		// Safety break loop if everyone is All-In/Folded
		if be.ActionIndex == start {
			// Everyone checked/folded/all-in. Round over.
			be.IsRoundClosed = true
			return
		}
	}
}

// checkRoundComplete is complex.
// Simplification: Round ends when everyone active is either All-In or has Matched HighBet,
// AND the Aggressor has been passed (everyone had a chance to respond).
func (be *BettingEngine) checkRoundComplete() bool {
	// 1. Are all active players matched?
	activeCount := 0
	matchedCount := 0

	for _, p := range be.Players {
		if p.Status == StatusFolded || p.Status == StatusOut {
			continue
		}
		activeCount++

		// Match condition: Bet == HighBet OR All-In
		if p.BetThisRound == be.CurrentHighBet || p.IsAllIn {
			matchedCount++
		}
	}

	// If only 1 remains, round is effectively over (handled by game loop, but here we close betting)
	if activeCount <= 1 {
		return true
	}

	// 2. Everyone has matched?
	if matchedCount < activeCount {
		return false
	}

	// 3. Does Next Actor have options?
	// The advance loop logic handles skipping all-ins.
	// If everyone matches, and we are back to a state where no one *needs* to act...
	// We need to ensure everyone had at least one chance.
	// Our "AggressorIndex" tracks the start of the current "Level".
	// If we advanced past the AggressorIndex without a Raise, we are done.

	// This makes "Advance" tricky. A better way is:
	// Is the Next Candidate == AggressorIndex?
	// If yes, and everything is matched, then yes.

	// If the next person to act IS the person who set the price (Aggressor),
	// and no one raised them, the round is done.
	if be.ActionIndex == be.AggressorIndex {
		return true
	}

	return false
}

// ResolvePots calculates Main and Side Pots at the end of a hand.
// This implements the "Water Level" algorithm.
func (be *BettingEngine) ResolvePots() ([]Pot, error) {
	// 1. Identify active contributors
	type Contributor struct {
		Seat     int
		Amount   int64
		IsFolded bool
	}
	var contributors []Contributor

	for _, p := range be.Players {
		if p.TotalCommitted > 0 {
			contributors = append(contributors, Contributor{
				Seat:     p.SeatIndex,
				Amount:   p.TotalCommitted,
				IsFolded: (p.Status == StatusFolded),
			})
		}
	}

	// Sort by committed amount (asc) to build pots bottom-up
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Amount < contributors[j].Amount
	})

	var pots []Pot
	currentLevel := int64(0)

	for i, c := range contributors {
		if c.Amount <= currentLevel {
			continue
		}

		// Create a pot for the slice (c.Amount - currentLevel)
		chunck := c.Amount - currentLevel
		potAmt := int64(0)
		var eligible []int

		// Collect from all players who put in at least this much
		for j := i; j < len(contributors); j++ {
			contrib := contributors[j]
			potAmt += chunck
			// Only non-folded players are eligible to WIN
			if !contrib.IsFolded {
				eligible = append(eligible, contrib.Seat)
			}
		}

		if potAmt > 0 {
			pots = append(pots, Pot{
				Amount:        potAmt,
				EligibleSeats: eligible,
			})
		}

		currentLevel = c.Amount
	}

	return pots, nil
}
