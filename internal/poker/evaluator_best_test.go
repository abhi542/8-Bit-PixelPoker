package poker

import (
	"testing"
)

func TestEvaluateBest6Card(t *testing.T) {
	// 6 cards: Ah Kh Qh Jh Th 2s (Royal Flush + garbage)
	cards := []Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: King, Suit: Hearts},
		{Rank: Queen, Suit: Hearts},
		{Rank: Jack, Suit: Hearts},
		{Rank: Ten, Suit: Hearts},
		{Rank: Two, Suit: Spades},
	}

	best := EvaluateBestHand(cards)
	if best.Rank != RoyalFlush {
		t.Errorf("expected Royal Flush, got %v", best.Rank)
	}
}

func TestEvaluateBest6Card_HighCardUpdate(t *testing.T) {
	// Ensure that even a weak hand updates 'best' correctly (from -1 initialization)
	// 6 cards: 2s 3d 4c 5h 7s 8d. Best 5 is 8,7,5,4,3 (High Card)
	cards := []Card{
		{Rank: Two, Suit: Spades},
		{Rank: Three, Suit: Diamonds},
		{Rank: Four, Suit: Clubs},
		{Rank: Five, Suit: Hearts},
		{Rank: Seven, Suit: Spades}, // Gap (6 missing)
		{Rank: Eight, Suit: Diamonds},
	}

	best := EvaluateBestHand(cards)
	if best.Rank != HighCard {
		t.Errorf("expected High Card, got %v", best.Rank)
	}
	// Check tie breakers to ensure it's not empty/nil
	if len(best.TieBreakers) != 5 {
		t.Errorf("expected 5 tie breakers, got %d", len(best.TieBreakers))
	}
}
