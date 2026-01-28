package poker

import (
	"testing"
)

func TestBettingEngine_SimpleFlow(t *testing.T) {
	// Setup 3 players. Blinds 10/20.
	p1 := &PlayerBetState{SeatIndex: 0, Chips: 1000, Status: StatusActive}
	p2 := &PlayerBetState{SeatIndex: 1, Chips: 1000, Status: StatusActive}
	p3 := &PlayerBetState{SeatIndex: 2, Chips: 1000, Status: StatusActive}

	players := []*PlayerBetState{p1, p2, p3}
	be := NewBettingEngine(players, 2, 20) // Dealer is p3. Action starts p1 (SB? No, usually UTG. Let's assume PreFlop manual setup)

	// Simulating Pre-Flop:
	// P1 (SB) posts 10. P2 (BB) posts 20. Action on P3 (UTG).
	p1.BetThisRound = 10
	p1.Chips -= 10
	p1.TotalCommitted = 10
	p2.BetThisRound = 20
	p2.Chips -= 20
	p2.TotalCommitted = 20
	be.CurrentHighBet = 20
	be.ActionIndex = 2    // P3 to act
	be.AggressorIndex = 2 // UTG is the loop stopper for PreFlop (unless raised)

	// P3 Calls 20
	err := be.ApplyAction(2, Action{Type: ActionCall})
	if err != nil {
		t.Fatalf("P3 Call failed: %v", err)
	}
	if p3.BetThisRound != 20 {
		t.Errorf("P3 bet %d, want 20", p3.BetThisRound)
	}

	// Action moves to P1 (SB)
	if be.ActionIndex != 0 {
		t.Errorf("Action should be P1, got P%d", be.ActionIndex+1)
	}

	// P1 Calls (adds 10 more)
	err = be.ApplyAction(0, Action{Type: ActionCall})
	if err != nil {
		t.Fatalf("P1 Call failed: %v", err)
	}
	if p1.BetThisRound != 20 {
		t.Errorf("P1 bet %d, want 20", p1.BetThisRound)
	}

	// Action moves to P2 (BB) - Option to check
	if be.ActionIndex != 1 {
		t.Errorf("Action should be P2, got P%d", be.ActionIndex+1)
	}

	// P2 Checks
	err = be.ApplyAction(1, Action{Type: ActionCheck})
	if err != nil {
		t.Fatalf("P2 Check failed: %v", err)
	}

	// Round should be close
	if !be.IsRoundClosed {
		t.Error("Round should be closed after BB Check")
	}
}

func TestBettingEngine_SidePots(t *testing.T) {
	// A: 100 (All in)
	// B: 500 (Bets 500)
	// C: 1000 (Calls 500)
	// Expected:
	// Main Pot (from 100 level): 100*3 = 300. Eligible: A, B, C
	// Side Pot 1 (from 400 level): 400*2 = 800. Eligible: B, C

	pA := &PlayerBetState{SeatIndex: 0, Chips: 100, Status: StatusActive, TotalCommitted: 100} // Started with 100
	pB := &PlayerBetState{SeatIndex: 1, Chips: 0, Status: StatusActive, TotalCommitted: 500}   // Started with 500
	pC := &PlayerBetState{SeatIndex: 2, Chips: 500, Status: StatusActive, TotalCommitted: 500} // Started with 1000

	be := &BettingEngine{
		Players: []*PlayerBetState{pA, pB, pC},
	}

	pots, err := be.ResolvePots()
	if err != nil {
		t.Fatalf("ResolvePots error: %v", err)
	}

	if len(pots) != 2 {
		t.Fatalf("Expected 2 pots, got %d", len(pots))
	}

	// Verify Main Pot
	if pots[0].Amount != 300 {
		t.Errorf("Main Pot amount %d, want 300", pots[0].Amount)
	}
	if len(pots[0].EligibleSeats) != 3 {
		t.Errorf("Main Pot eligible %v, want 3 players", pots[0].EligibleSeats)
	}

	// Verify Side Pot
	if pots[1].Amount != 800 {
		t.Errorf("Side Pot amount %d, want 800", pots[1].Amount)
	}
	if len(pots[1].EligibleSeats) != 2 {
		t.Errorf("Side Pot eligible %v, want 2 players", pots[1].EligibleSeats)
	}
}

func TestBettingEngine_RaiseLogic(t *testing.T) {
	// MinRaise logic
	p1 := &PlayerBetState{SeatIndex: 0, Chips: 1000, BetThisRound: 100}
	p2 := &PlayerBetState{SeatIndex: 1, Chips: 1000, BetThisRound: 0}

	be := &BettingEngine{
		Players:        []*PlayerBetState{p1, p2},
		CurrentHighBet: 100,
		MinRaise:       50,
		ActionIndex:    1,
	}

	// P2 tries to raise to 120 (Raise 20). Invalid, MinRaise 50 -> Needs 150 total.
	err := be.ApplyAction(1, Action{Type: ActionRaise, Amount: 120})
	if err == nil {
		t.Error("Expected error for small raise, got nil")
	}

	// P2 Raise to 150. Valid.
	err = be.ApplyAction(1, Action{Type: ActionRaise, Amount: 150})
	if err != nil {
		t.Errorf("Valid raise failed: %v", err)
	}

	// Verify State Update
	if be.CurrentHighBet != 150 {
		t.Errorf("HighBet %d, want 150", be.CurrentHighBet)
	}
	if be.MinRaise != 50 {
		t.Errorf("MinRaise %d, want 50 (50 on top of 100)", be.MinRaise)
	}
}
