package poker

type ActionType int

const (
	ActionFold ActionType = iota
	ActionCheck
	ActionCall
	ActionBet
	ActionRaise
	ActionAllIn // Explicit All-In action often maps to Bet/Raise/Call internally, but having a distinct type helps UI/Logic
)

func (a ActionType) String() string {
	switch a {
	case ActionFold:
		return "FOLD"
	case ActionCheck:
		return "CHECK"
	case ActionCall:
		return "CALL"
	case ActionBet:
		return "BET"
	case ActionRaise:
		return "RAISE"
	case ActionAllIn:
		return "ALL_IN"
	default:
		return "UNKNOWN"
	}
}

// Action represents a player's move.
type Action struct {
	Type   ActionType
	Amount int64 // For Bet/Raise, this is the TOTAL amount the player wants to have in the pot for this round.
}

// PlayerState tracks the betting status of a single player in the current round.
// This is separate from the permanent Seat struct to keep logic isolated and testable.
type PlayerBetState struct {
	SeatIndex      int          `json:"seat_index"`
	Chips          int64        `json:"chips"`
	BetThisRound   int64        `json:"bet_this_round"`
	TotalCommitted int64        `json:"total_committed"`
	Status         PlayerStatus `json:"status"`
	IsAllIn        bool         `json:"is_all_in"`
}

type PlayerStatus int

const (
	StatusActive PlayerStatus = iota
	StatusFolded
	StatusOut          // Empty seat or sitting out
	StatusDisconnected // Player definitely disconnected but seat reserved
)
