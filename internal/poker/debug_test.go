package poker

import (
	"fmt"
	"testing"
)

func TestDebugFlow(t *testing.T) {
	table := NewTable("test-table", TableConfig{SmallBlind: 10, BigBlind: 20})
	table.JoinSeat(0, "p1", "Alice", 1000)
	table.JoinSeat(1, "p2", "Bob", 1000)
	table.JoinSeat(2, "p3", "Charlie", 1000)
	table.StartHand()

	fmt.Printf("Start: ActionIndex=%d, HighBet=%d\n", table.Betting.ActionIndex, table.Betting.CurrentHighBet)
	if table.Betting.ActionIndex != 0 {
		t.Errorf("Expected Action 0, got %d", table.Betting.ActionIndex)
	}

	// Alice Call
	err := table.HandleAction(0, Action{Type: ActionCall, Amount: 20})
	if err != nil {
		t.Fatalf("Alice Call failed: %v", err)
	}

	fmt.Printf("After Alice: ActionIndex=%d, HighBet=%d\n", table.Betting.ActionIndex, table.Betting.CurrentHighBet)

	// Check Bob
	bob := table.Betting.Players[1] // Should be Bob
	fmt.Printf("Bob (Seat %d): Status=%d, Bet=%d\n", bob.SeatIndex, bob.Status, bob.BetThisRound)

	if table.Betting.ActionIndex != 1 {
		t.Errorf("Expected Action 1 (Bob), got %d", table.Betting.ActionIndex)
	}
}
