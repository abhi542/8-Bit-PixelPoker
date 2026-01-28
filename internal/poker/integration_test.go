package poker

import (
	"testing"
)

// TestFullHandFlow simulates a complete hand to verify Table + Betting + Evaluator integration.
func TestFullHandFlow(t *testing.T) {
	// Setup
	table := NewTable("test-table", TableConfig{SmallBlind: 10, BigBlind: 20})

	// Join 3 Players
	table.JoinSeat(0, "p1", "Alice", 1000)
	table.JoinSeat(1, "p2", "Bob", 1000)
	table.JoinSeat(2, "p3", "Charlie", 1000)

	// Start Hand
	if err := table.StartHand(); err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	// Verify State: PreFlop
	if table.State.State != StateBettingPreFlop {
		t.Errorf("Expected PreFlop, got %s", table.State.State)
	}

	// Verify Hands Dealt
	if len(table.Seats[0].HoleCards) != 2 {
		t.Error("Alice not dealt cards")
	}

	// Verify Initial Hand Strength (High Card or Pair)
	if table.Seats[0].CurrentHandRank == 0 { // 0 is HighCard in iota? No, HighCard is 0.
		// Wait, HighCard is 0. If it wasn't calculated, it would be 0.
		// But calculateHandStrengths is called.
		// We can't distinguish 0 (default) from 0 (HighCard) easily here without checking desc.
		if table.Seats[0].CurrentHandDesc == "" {
			t.Error("Hand description empty")
		}
	}

	// Betting: PreFlop
	// Dealer 0 -> SB 1, BB 2, UTG 0.
	// Action on 0 (Alice).
	// Alice Calls 20.
	// Bob (SB 10) completes to 20.
	// Charlie (BB 20) Checks.

	if err := table.HandleAction(0, Action{Type: ActionCall, Amount: 20}); err != nil {
		t.Fatalf("Alice Call failed: %v", err)
	}
	if err := table.HandleAction(1, Action{Type: ActionCall, Amount: 20}); err != nil { // SB calls extra 10
		t.Fatalf("Bob Call failed: %v", err)
	}
	if err := table.HandleAction(2, Action{Type: ActionCheck, Amount: 0}); err != nil {
		t.Fatalf("Charlie Check failed: %v", err)
	}

	// Should advance to Flop
	if table.State.State != StateBettingFlop {
		t.Fatalf("Expected Flop, got %s", table.State.State)
	}
	if len(table.Community) != 3 {
		t.Errorf("Expected 3 community cards, got %d", len(table.Community))
	}

	// Verify Hand Strength Update
	// Just check it didn't crash and fields are populated.
	descFlop := table.Seats[0].CurrentHandDesc
	if descFlop == "" {
		t.Error("Hand desc empty on Flop")
	}

	// Betting: Flop
	// Bob Bets 50, Charlie Folds, Alice Calls.
	// SB (1) acts first post-flop.

	if err := table.HandleAction(1, Action{Type: ActionBet, Amount: 50}); err != nil {
		t.Fatalf("Bob Bet failed: %v", err)
	}
	if err := table.HandleAction(2, Action{Type: ActionFold, Amount: 0}); err != nil {
		t.Fatalf("Charlie Fold failed: %v", err)
	}
	if err := table.HandleAction(0, Action{Type: ActionCall, Amount: 50}); err != nil {
		t.Fatalf("Alice Call failed: %v", err)
	}

	// Turn
	if table.State.State != StateBettingTurn {
		t.Fatalf("Expected Turn, got %s", table.State.State)
	}
	if len(table.Community) != 4 {
		t.Errorf("Expected 4 community cards")
	}

	// Check/Check
	table.HandleAction(1, Action{Type: ActionCheck})
	table.HandleAction(0, Action{Type: ActionCheck})

	// River
	if table.State.State != StateBettingRiver {
		t.Fatalf("Expected River, got %s", table.State.State)
	}

	// Check/Check
	table.HandleAction(1, Action{Type: ActionCheck})
	table.HandleAction(0, Action{Type: ActionCheck})

	// Showdown -> HandEnd -> Payout
	// The StateMachine transitions Showdown -> Payout logic?
	// Currently Table.nextStreet() just advances state.
	// It doesn't auto-resolve payout in MVP nextStreet implementation probably?
	// Checking table.go:
	// case StateBettingRiver: ... t.State.AdvanceState(StateShowdown) ...
	// Wait, nextStreet logic for River needs to be checked.
}
