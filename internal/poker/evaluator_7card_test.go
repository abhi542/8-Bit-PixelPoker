package poker

import (
	"testing"
)

func TestEvaluate7Card(t *testing.T) {
	tests := []struct {
		name         string
		cards        []Card
		expectedRank HandRank
		// We could check tie-breakers too, but verifying Rank is the primary correctness check for the combo logic
	}{
		{
			// Board: Ah, Kh, Qh, Jh, 2c
			// Hand:  Th, 3c
			// Best:  Ah, Kh, Qh, Jh, Th (Royal Flush)
			"Royal Flush",
			[]Card{
				{Hearts, Ace}, {Hearts, King}, {Hearts, Queen}, {Hearts, Jack}, {Clubs, Two},
				{Hearts, Ten}, {Clubs, Three},
			},
			RoyalFlush,
		},
		{
			// Board: 2h, 3h, 4h, 5h, 9d
			// Hand:  6h, Ah
			// Best:  2h, 3h, 4h, 5h, 6h (Straight Flush) - Higher than Ace high flush?
			// Wait, A,2,3,4,5,6 of hearts.
			// 6-high straight flush (2,3,4,5,6)
			// Ace-high flush (A,6,5,4,3)
			// SF wins.
			"Straight Flush over Flush",
			[]Card{
				{Hearts, Two}, {Hearts, Three}, {Hearts, Four}, {Hearts, Five}, {Diamonds, Nine},
				{Hearts, Six}, {Hearts, Ace},
			},
			StraightFlush,
		},
		{
			// Quads
			// Board: 8c, 8d, 8h, 2s, 3s
			// Hand:  8s, Ac
			// Best:  8,8,8,8,A
			"Four of a Kind",
			[]Card{
				{Clubs, Eight}, {Diamonds, Eight}, {Hearts, Eight}, {Spades, Two}, {Spades, Three},
				{Spades, Eight}, {Clubs, Ace},
			},
			FourOfAKind,
		},
		{
			// Full House
			// Hand: KK
			// Board: 888 2 3
			// Best: 888KK (Full House) > 888A2 (Trips)
			"Full House with Pocket Pair",
			[]Card{
				{Clubs, King}, {Diamonds, King}, // Pocket
				{Hearts, Eight}, {Spades, Eight}, {Clubs, Eight}, {Diamonds, Two}, {Diamonds, Three}, // Board
			},
			FullHouse,
		},
		{
			// Flush (6 card flush, must pick best 5)
			// Board: Ah, Kh, 9h, 2h, 2c
			// Hand:  5h, 4h
			// Available: A, K, 9, 2, 5, 4 (Hearts).
			// Best 5: A, K, 9, 5, 4
			"Flush (Select Best 5)",
			[]Card{
				{Hearts, Ace}, {Hearts, King}, {Hearts, Nine}, {Hearts, Two}, {Clubs, Two},
				{Hearts, Five}, {Hearts, Four},
			},
			Flush,
		},
		{
			// Straight (Counterfeit protection)
			// Hand: 9, T
			// Board: J, Q, K, A, A
			// Best: T, J, Q, K, A (Straight) - Pair of Aces doesn't matter if Straight is better
			"Broadway Straight ignoring Pair",
			[]Card{
				{Clubs, Nine}, {Diamonds, Ten},
				{Hearts, Jack}, {Spades, Queen}, {Clubs, King}, {Diamonds, Ace}, {Spades, Ace},
			},
			Straight,
		},
		{
			// Two Pair (3 pairs on board/hand)
			// Hand: KK
			// Board: QQ JJ 2
			// Best: KK QQ A (Two Pair with K)
			"Two Pair (Best of 3 Pairs)",
			[]Card{
				{Clubs, King}, {Diamonds, King},
				{Hearts, Queen}, {Spades, Queen}, {Clubs, Jack}, {Diamonds, Jack}, {Spades, Two},
			},
			TwoPair,
		},
		{
			// High Card
			// Hand: 2, 3 (offsuit)
			// Board: 5, 7, 9, J, K (all diff suits)
			// Best: K, J, 9, 7, 5
			"High Card",
			[]Card{
				{Clubs, Two}, {Diamonds, Three},
				{Hearts, Five}, {Spades, Seven}, {Clubs, Nine}, {Diamonds, Jack}, {Spades, King},
			},
			HighCard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate7Card(tt.cards)
			if err != nil {
				t.Fatalf("Evaluate7Card() error = %v", err)
			}
			if got.Rank != tt.expectedRank {
				t.Errorf("Rank = %v, want %v. TieBreakers: %v", got.Rank, tt.expectedRank, got.TieBreakers)
			}
		})
	}
}
