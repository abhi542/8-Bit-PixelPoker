package poker

// GameState represents the current phase of the poker hand.
type GameState int

const (
	StateWaiting        GameState = iota // Waiting for players
	StateDealPrivate                     // Dealing hole cards
	StateBettingPreFlop                  // Pre-flop betting round
	StateDealFlop                        // Dealing 3 community leads
	StateBettingFlop                     // Flop betting round
	StateDealTurn                        // Dealing turn card
	StateBettingTurn                     // Turn betting round
	StateDealRiver                       // Dealing river card
	StateBettingRiver                    // River betting round
	StateShowdown                        // Revealing hands
	StatePayout                          // Distributing pot
	StateHandEnd                         // Cleanup before next hand
)

func (s GameState) String() string {
	switch s {
	case StateWaiting:
		return "WAITING"
	case StateDealPrivate:
		return "DEAL_PRIVATE"
	case StateBettingPreFlop:
		return "BETTING_PRE_FLOP"
	case StateDealFlop:
		return "DEAL_FLOP"
	case StateBettingFlop:
		return "BETTING_FLOP"
	case StateDealTurn:
		return "DEAL_TURN"
	case StateBettingTurn:
		return "BETTING_TURN"
	case StateDealRiver:
		return "DEAL_RIVER"
	case StateBettingRiver:
		return "BETTING_RIVER"
	case StateShowdown:
		return "SHOWDOWN"
	case StatePayout:
		return "PAYOUT"
	case StateHandEnd:
		return "HAND_END"
	default:
		return "UNKNOWN"
	}
}

// IsBettingRound returns true if the state accepts player actions.
func (s GameState) IsBettingRound() bool {
	switch s {
	case StateBettingPreFlop, StateBettingFlop, StateBettingTurn, StateBettingRiver:
		return true
	default:
		return false
	}
}
