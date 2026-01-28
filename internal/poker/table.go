package poker

import (
	"errors"
	"fmt"
	"sync"
)

// Seat represents a player at the table.
type Seat struct {
	Index           int
	PlayerID        string
	Username        string
	Chips           int64           // Stack on table
	HoleCards       []Card          // Private cards
	Status          PlayerStatus    // Active, Folded, SittingOut
	BetState        *PlayerBetState // Linked state for betting engine
	CurrentHandRank HandRank        `json:"hand_rank"`
	CurrentHandDesc string          `json:"hand_desc"`
}

// Table represents a single poker game instance.
// It aggregates the StateMachine, BettingEngine, and Deck.
type Table struct {
	mu        sync.RWMutex
	ID        string
	State     *TableStateMachine
	Betting   *BettingEngine
	Deck      Deck
	Seats     [9]*Seat // Fixed 9 seats
	Community []Card
	Winners   []int `json:"winners"` // Indices of winners
	Config    TableConfig
}

type TableConfig struct {
	SmallBlind int64
	BigBlind   int64
}

func NewTable(id string, config TableConfig) *Table {
	t := &Table{
		ID:     id,
		Config: config,
		State:  NewTableStateMachine(),
		Deck:   NewDeck(),
	}
	// Init empty seats
	for i := 0; i < 9; i++ {
		t.Seats[i] = nil
	}
	return t
}

// JoinSeat attempts to sit a player at a specific seat.
func (t *Table) JoinSeat(seatIdx int, playerID, username string, buyIn int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if seatIdx < 0 || seatIdx >= 9 {
		return errors.New("invalid seat index")
	}
	if t.Seats[seatIdx] != nil {
		return errors.New("seat occupied")
	}

	// Create Seat
	seat := &Seat{
		Index:    seatIdx,
		PlayerID: playerID,
		Username: username,
		Chips:    buyIn,
		Status:   StatusOut, // Sit out until next hand starts
	}
	t.Seats[seatIdx] = seat
	return nil
}

// StartHand initializes a new hand if possible.
func (t *Table) StartHand() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 1. Check player count (need 2 active)
	active := 0
	for _, s := range t.Seats {
		if s != nil && s.Chips > 0 {
			active++
			s.Status = StatusActive
			s.HoleCards = nil
		}
	}
	if active < 2 {
		return errors.New("not enough players")
	}

	// 2. Reset State
	t.State.AdvanceState(StateDealPrivate)
	t.Deck = NewDeck()
	t.Deck.Shuffle(CryptoRNG{}) // Production Shuffle
	t.Community = nil

	// 3. Init Betting Engine
	// Collect active players for this hand
	var playersForBet []*PlayerBetState
	for _, s := range t.Seats {
		if s != nil && s.Status == StatusActive {
			// Create/Reset Betting State
			pbs := &PlayerBetState{
				SeatIndex: s.Index,
				Chips:     s.Chips,
				Status:    StatusActive,
			}
			s.BetState = pbs
			playersForBet = append(playersForBet, pbs)
		}
	}

	// Simple Dealer rotation needed here (mocked to 0 for now)
	dealerIdx := 0
	t.Betting = NewBettingEngine(playersForBet, dealerIdx, t.Config.BigBlind)

	// 4. Deal Cards
	// Deal 2 to each active player
	for _, s := range t.Seats {
		if s != nil && s.Status == StatusActive {
			cards, _ := t.Deck.Draw(2)
			s.HoleCards = cards
		}
	}

	// 5. Post Blinds
	// Handle Heads-up vs Multi-way
	n := len(playersForBet)
	sbIdx := (dealerIdx + 1) % n
	bbIdx := (dealerIdx + 2) % n
	actionIdx := (dealerIdx + 3) % n

	if n == 2 {
		// Heads-up: Dealer is SB, Other is BB. Action on SB (Dealer).
		sbIdx = dealerIdx
		bbIdx = (dealerIdx + 1) % n
		actionIdx = sbIdx
	}

	// Apply SB
	sbPlayer := playersForBet[sbIdx]
	sbAmt := t.Config.SmallBlind
	if sbAmt > sbPlayer.Chips {
		sbAmt = sbPlayer.Chips
	}
	sbPlayer.Chips -= sbAmt
	sbPlayer.BetThisRound += sbAmt
	sbPlayer.TotalCommitted += sbAmt

	// Apply BB
	bbPlayer := playersForBet[bbIdx]
	bbAmt := t.Config.BigBlind
	if bbAmt > bbPlayer.Chips {
		bbAmt = bbPlayer.Chips
	}
	bbPlayer.Chips -= bbAmt
	bbPlayer.BetThisRound += bbAmt
	bbPlayer.TotalCommitted += bbAmt

	// Update Engine State
	t.Betting.CurrentHighBet = bbAmt // Assuming BB > SB
	if sbAmt > bbAmt {
		t.Betting.CurrentHighBet = sbAmt
	} // Rare
	t.Betting.MinRaise = t.Config.BigBlind
	t.Betting.ActionIndex = actionIdx
	t.Betting.AggressorIndex = actionIdx // UTG/Stopper for PreFlop round loop

	// Ensure the ActionIndex player is valid (not all-in from blind? loop handled by engine usually, but here we set directly)
	// If Action player is all-in? Unlikely at start unless short stack.
	// We strictly trust the index for now.

	t.State.AdvanceState(StateBettingPreFlop)

	// Initial Rank check (e.g. pockets)
	t.updateHandRanks()

	return nil
}

// updateHandRanks calculates and stores the current best hand for each active seat.
func (t *Table) updateHandRanks() {
	// caller holds lock
	for _, s := range t.Seats {
		if s == nil || s.Status == StatusFolded || s.HoleCards == nil {
			continue
		}

		// Combine Hole + Community
		var allCards []Card
		allCards = append(allCards, s.HoleCards...)
		allCards = append(allCards, t.Community...)

		best := EvaluateBestHand(allCards)
		s.CurrentHandRank = best.Rank
		s.CurrentHandDesc = best.Rank.String()
	}
}

// HandleAction delegates to BettingEngine and StateMachine.
func (t *Table) HandleAction(seatIdx int, action Action) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Betting == nil {
		return errors.New("game not started")
	}

	// Delegate to Betting Engine
	if err := t.Betting.ApplyAction(seatIdx, action); err != nil {
		return err
	}

	// Sync Chips back to Seat (if needed for persistence, though using BetState is cleaner)
	s := t.Seats[seatIdx]
	if s != nil && s.BetState != nil {
		s.Chips = s.BetState.Chips
		// Sync Status
		if s.BetState.Status == StatusFolded {
			s.Status = StatusFolded
		}
	}

	// Check for Round End
	if t.Betting.IsRoundClosed {
		t.nextStreet()
	}

	return nil
}

func (t *Table) nextStreet() {
	// Logic to deal cards and move state phases
	// E.g. PreFlop -> DealFlop -> BettingFlop
	current := t.State.State

	switch current {
	case StateBettingPreFlop:
		t.State.AdvanceState(StateDealFlop)
		t.Deck.Burn()
		cards, _ := t.Deck.Draw(3)
		t.Community = append(t.Community, cards...)
		t.updateHandRanks()
		t.State.AdvanceState(StateBettingFlop)

		// Reset Betting for new round
		if t.Betting != nil {
			t.Betting.StartNewRound()
		}

	case StateBettingFlop:
		t.State.AdvanceState(StateDealTurn)
		t.Deck.Burn()
		cards, _ := t.Deck.Draw(1)
		t.Community = append(t.Community, cards...)
		t.updateHandRanks()
		t.State.AdvanceState(StateBettingTurn)

		if t.Betting != nil {
			t.Betting.StartNewRound()
		}

	case StateBettingTurn:
		t.State.AdvanceState(StateDealRiver)
		t.Deck.Burn()
		cards, _ := t.Deck.Draw(1)
		t.Community = append(t.Community, cards...)
		t.updateHandRanks()
		t.State.AdvanceState(StateBettingRiver)
		// Reset betting for River
		if t.Betting != nil {
			t.Betting.StartNewRound()
		}

	case StateBettingRiver:
		// Critical: Ensure betting is closed before showing down
		t.processShowdown()
	}
}

func (t *Table) processShowdown() {
	// Caller (HandleAction) holds t.mu Lock. Do NOT lock here.

	fmt.Println("=== START SHOWDOWN ===")
	t.State.AdvanceState(StateShowdown)
	t.updateHandRanks() // Final check

	// Solve for winners using Water Levels (Side Pots)
	pots, err := t.Betting.ResolvePots()
	if err != nil {
		fmt.Printf("CRITICAL ERROR Resolving Pots: %v\n", err)
		return
	}

	if len(pots) > 0 {
		t.Betting.MainPot = pots[0].Amount
		if len(pots) > 1 {
			t.Betting.SidePots = pots[1:]
		} else {
			t.Betting.SidePots = nil
		}
	} else {
		t.Betting.MainPot = 0
		t.Betting.SidePots = nil
	}

	// DEBUG: Print Pots
	for i, p := range pots {
		fmt.Printf("Pot %d: %d chips (Eligible: %v)\n", i, p.Amount, p.EligibleSeats)
	}

	// Track all unique winners for the frontend
	uniqueWinners := make(map[int]bool)

	// For each pot (main + sides)
	for _, pot := range pots {
		if pot.Amount == 0 {
			continue
		}

		// Find winner(s) among eligible
		var winners []int
		var bestHand EvaluatedHand
		bestHand.Rank = -1 // Starts lower than HighCard (0)

		// Evaluate all eligible to find Best
		for _, seatIdx := range pot.EligibleSeats {
			s := t.Seats[seatIdx]
			if s == nil {
				continue
			}

			// Re-eval just to be safe (or trust updateHandRanks)
			var all []Card
			all = append(all, s.HoleCards...)
			all = append(all, t.Community...)
			hand := EvaluateBestHand(all)

			fmt.Printf("Seat %d (%s) Hand: %s\n", seatIdx, s.Username, hand.Rank.String())

			// Compare
			if bestHand.Rank == -1 {
				bestHand = hand
				winners = []int{seatIdx}
			} else {
				cmp := hand.Compare(bestHand)
				if cmp > 0 {
					bestHand = hand
					winners = []int{seatIdx}
				} else if cmp == 0 {
					winners = append(winners, seatIdx)
				}
			}
		}

		// Payout
		if len(winners) > 0 {
			share := pot.Amount / int64(len(winners))
			remainder := pot.Amount % int64(len(winners))

			fmt.Printf("Winners for Pot %d: %v. Share: %d\n", pot.Amount, winners, share)

			for i, wIdx := range winners {
				uniqueWinners[wIdx] = true // Track for UI
				t.Seats[wIdx].Chips += share
				if i == 0 {
					t.Seats[wIdx].Chips += remainder // Odd chip to first
				}
				// Sync BetState too, just in case
				if t.Seats[wIdx].BetState != nil {
					t.Seats[wIdx].BetState.Chips = t.Seats[wIdx].Chips
				}
			}
		}
	}

	// Assign to Table struct
	t.Winners = nil
	for wIdx := range uniqueWinners {
		t.Winners = append(t.Winners, wIdx)
	}

	fmt.Println("=== END SHOWDOWN ===")
}

// GetViewFor returns a sanitized copy of the table for a specific player.
// It hides hole cards of other players unless the game is in Showdown.
func (t *Table) GetViewFor(playerID string) *Table {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 1. Shallow Copy Table
	view := *t

	// 2. Deep Copy Seats (to modify them without touching original)
	// We only need to clone the array and the structs we modify.
	var newSeats [9]*Seat
	for i, s := range t.Seats {
		if s == nil {
			continue
		}
		// Copy Seat
		sCopy := *s
		newSeats[i] = &sCopy

		// 3. Logic: Hide Cards?
		// Show if: It's ME OR GameState is Showdown OR (maybe) player is All-In/Folded? (Folded usually muck).
		// Standard: Show only mine. Showdown shows all active.
		shouldShow := false
		if s.PlayerID == playerID {
			shouldShow = true
		} else if t.State.State == StateShowdown {
			shouldShow = true
		}

		if !shouldShow {
			newSeats[i].HoleCards = nil // Hide
			// We might want to indicate "2 cards" exist?
			// The frontend knows "active" implies 2 cards.
			// Ideally we send "backed" cards or just nil and frontend infers.
			// Frontend: if active && cards==nil -> draw backs.
		}
	}
	view.Seats = newSeats
	return &view
}
