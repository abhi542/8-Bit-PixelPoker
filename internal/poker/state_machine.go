package poker

import (
	"fmt"
)

// TableStateMachine manages the flow of the game.
// It embeds the data required for state transitions.
type TableStateMachine struct {
	State GameState
	// We will add more fields (Deck, Seats, etc.) in later steps.
	// For now, this struct focuses SOLELY on state transitions.
}

// NewTableStateMachine creates a table in Waiting state.
func NewTableStateMachine() *TableStateMachine {
	return &TableStateMachine{
		State: StateWaiting,
	}
}

// AdvanceState moves the game to the next state if valid.
// This is the SINGLE authority on transitions.
func (t *TableStateMachine) AdvanceState(next GameState) error {
	if !t.isValidTransition(next) {
		return fmt.Errorf("invalid transition from %s to %s", t.State, next)
	}

	t.State = next
	// In a full implementation, side effects would trigger here
	// e.g. if next == StateDealFlop, t.Deck.BurnAndDraw(3)

	return nil
}

// isValidTransition defines the strict directed graph of poker states.
func (t *TableStateMachine) isValidTransition(next GameState) bool {
	current := t.State

	switch current {
	case StateWaiting:
		// Can only start dealing if we have enough players (checked by caller)
		return next == StateDealPrivate

	case StateDealPrivate:
		return next == StateBettingPreFlop

	case StateBettingPreFlop:
		// Normal flow -> Flop.
		// If everyone folds -> Payout (Hand ends early).
		return next == StateDealFlop || next == StatePayout

	case StateDealFlop:
		return next == StateBettingFlop

	case StateBettingFlop:
		return next == StateDealTurn || next == StatePayout

	case StateDealTurn:
		return next == StateBettingTurn

	case StateBettingTurn:
		return next == StateDealRiver || next == StatePayout

	case StateDealRiver:
		return next == StateBettingRiver

	case StateBettingRiver:
		return next == StateShowdown || next == StatePayout

	case StateShowdown:
		return next == StatePayout

	case StatePayout:
		return next == StateHandEnd

	case StateHandEnd:
		// Loop back to waiting (or dealing immediately)
		return next == StateWaiting || next == StateDealPrivate
	}

	return false
}

// ForceState allows admin/debug overrides, should only be used in tests or crash recovery.
func (t *TableStateMachine) ForceState(s GameState) {
	t.State = s
}
