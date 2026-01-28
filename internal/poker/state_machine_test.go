package poker

import (
	"testing"
)

func TestTableStateMachine_Transitions(t *testing.T) {
	sm := NewTableStateMachine()

	// 1. Initial State
	if sm.State != StateWaiting {
		t.Fatalf("Expected WAITING, got %s", sm.State)
	}

	// 2. Valid Flow: Waiting -> DealPrivate -> BettingPreFlop
	if err := sm.AdvanceState(StateDealPrivate); err != nil {
		t.Errorf("Failed transition to DealPrivate: %v", err)
	}
	if err := sm.AdvanceState(StateBettingPreFlop); err != nil {
		t.Errorf("Failed transition to BettingPreFlop: %v", err)
	}

	// 3. Invalid Jump: BettingPreFlop -> River (Skip Flop/Turn)
	if err := sm.AdvanceState(StateDealRiver); err == nil {
		t.Error("Expected error transitioning PreFlop -> River, got nil")
	}

	// 4. Valid Flow: PreFlop -> Flop -> BettingFlop
	if err := sm.AdvanceState(StateDealFlop); err != nil {
		t.Errorf("Failed transition to Flop: %v", err)
	}
	if err := sm.AdvanceState(StateBettingFlop); err != nil {
		t.Errorf("Failed transition to BettingFlop: %v", err)
	}

	// 5. Early Win Scenario (Everyone folds)
	// Reset to PreFlop
	sm.ForceState(StateBettingPreFlop)
	if err := sm.AdvanceState(StatePayout); err != nil {
		t.Errorf("Failed transition to Payout (Early Win): %v", err)
	}
}

func TestGameState_String(t *testing.T) {
	if StateWaiting.String() != "WAITING" {
		t.Error("String() failed for WAITING")
	}
	if StateDealRiver.String() != "DEAL_RIVER" {
		t.Error("String() failed for DEAL_RIVER")
	}
}

func TestGameState_IsBettingRound(t *testing.T) {
	if !StateBettingPreFlop.IsBettingRound() {
		t.Error("PreFlop should be betting round")
	}
	if StateDealFlop.IsBettingRound() {
		t.Error("DealFlop should NOT be betting round")
	}
}
