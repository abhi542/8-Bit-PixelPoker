package poker

import (
	"reflect"
	"testing"
)

func TestEvaluate5Card(t *testing.T) {
	tests := []struct {
		name         string
		cards        []Card
		expectedRank HandRank
		expectedTie  []Rank
	}{
		{
			"Royal Flush",
			[]Card{
				{Spades, Ace}, {Spades, King}, {Spades, Queen}, {Spades, Jack}, {Spades, Ten},
			},
			RoyalFlush,
			[]Rank{Ace},
		},
		{
			"Straight Flush",
			[]Card{
				{Hearts, Nine}, {Hearts, Eight}, {Hearts, Seven}, {Hearts, Six}, {Hearts, Five},
			},
			StraightFlush,
			[]Rank{Nine},
		},
		{
			"Wheel Straight Flush",
			[]Card{
				{Diamonds, Five}, {Diamonds, Four}, {Diamonds, Three}, {Diamonds, Two}, {Diamonds, Ace},
			},
			StraightFlush,
			[]Rank{Five},
		},
		{
			"Four of a Kind",
			[]Card{
				{Clubs, Eight}, {Diamonds, Eight}, {Hearts, Eight}, {Spades, Eight}, {Clubs, Two},
			},
			FourOfAKind,
			[]Rank{Eight, Two},
		},
		{
			"Full House",
			[]Card{
				{Clubs, Ten}, {Diamonds, Ten}, {Hearts, Ten}, {Spades, Nine}, {Clubs, Nine},
			},
			FullHouse,
			[]Rank{Ten, Nine},
		},
		{
			"Flush",
			[]Card{
				{Spades, Ace}, {Spades, Jack}, {Spades, Eight}, {Spades, Six}, {Spades, Four},
			},
			Flush,
			[]Rank{Ace, Jack, Eight, Six, Four},
		},
		{
			"Straight",
			[]Card{
				{Clubs, Six}, {Diamonds, Five}, {Hearts, Four}, {Spades, Three}, {Clubs, Two},
			},
			Straight,
			[]Rank{Six},
		},
		{
			"Wheel Straight",
			[]Card{
				{Hearts, Five}, {Spades, Four}, {Clubs, Three}, {Diamonds, Two}, {Hearts, Ace},
			},
			Straight,
			[]Rank{Five},
		},
		{
			"Three of a Kind",
			[]Card{
				{Clubs, Seven}, {Diamonds, Seven}, {Hearts, Seven}, {Spades, King}, {Clubs, Two},
			},
			ThreeOfAKind,
			[]Rank{Seven, King, Two},
		},
		{
			"Two Pair",
			[]Card{
				{Clubs, Jack}, {Diamonds, Jack}, {Hearts, Nine}, {Spades, Nine}, {Clubs, Ace},
			},
			TwoPair,
			[]Rank{Jack, Nine, Ace},
		},
		{
			"Pair",
			[]Card{
				{Clubs, Ten}, {Diamonds, Ten}, {Hearts, King}, {Spades, Queen}, {Clubs, Eight},
			},
			Pair,
			[]Rank{Ten, King, Queen, Eight},
		},
		{
			"High Card",
			[]Card{
				{Clubs, King}, {Diamonds, Jack}, {Hearts, Nine}, {Spades, Seven}, {Clubs, Five},
			},
			HighCard,
			[]Rank{King, Jack, Nine, Seven, Five},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate5Card(tt.cards)
			if err != nil {
				t.Fatalf("Evaluate5Card() error = %v", err)
			}
			if got.Rank != tt.expectedRank {
				t.Errorf("Rank = %v, want %v", got.Rank, tt.expectedRank)
			}
			if !reflect.DeepEqual(got.TieBreakers, tt.expectedTie) {
				t.Errorf("TieBreakers = %v, want %v", got.TieBreakers, tt.expectedTie)
			}
		})
	}
}

func TestEvaluatedHand_Compare(t *testing.T) {
	// Full House vs Flush
	fh := EvaluatedHand{Rank: FullHouse, TieBreakers: []Rank{Ten, Nine}}
	fl := EvaluatedHand{Rank: Flush, TieBreakers: []Rank{Ace, King, Queen, Jack, Nine}}
	if fh.Compare(fl) != 1 {
		t.Error("FullHouse should beat Flush")
	}

	// Pair vs Pair (Higher Pair)
	p1 := EvaluatedHand{Rank: Pair, TieBreakers: []Rank{Jack, Ace, King, Queen}}
	p2 := EvaluatedHand{Rank: Pair, TieBreakers: []Rank{Ten, Ace, King, Queen}}
	if p1.Compare(p2) != 1 {
		t.Error("Pair of Jacks should beat Pair of Tens")
	}

	// Pair vs Pair (Same Pair, Higher Kicker)
	p3 := EvaluatedHand{Rank: Pair, TieBreakers: []Rank{Jack, Ace, King, Eight}}
	if p1.Compare(p3) != 1 { // p1 kicker is Queen (12), p3 kicker is Eight (8)?? No, p1 Kicker3 is Queen. p3 Kicker3 is Eight.
		// Wait, test data struct:
		// p1: [Jack, Ace, King, Queen]
		// p3: [Jack, Ace, King, Eight]
		// Queen > Eight.
		t.Error("Pair J (kickers AKQ) should beat Pair J (kickers AK8)")
	}

	// Tie
	p4 := EvaluatedHand{Rank: Pair, TieBreakers: []Rank{Jack, Ace, King, Queen}}
	if p1.Compare(p4) != 0 {
		t.Error("Identical hands should tie")
	}
}
