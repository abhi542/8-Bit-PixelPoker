package poker

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	tests := []struct {
		name    string
		s       Suit
		r       Rank
		wantErr bool
	}{
		{"ValidAceSpades", Spades, Ace, false},
		{"ValidTwoClubs", Clubs, Two, false},
		{"InvalidSuit", Suit(10), Ace, true},
		{"InvalidRankLow", Spades, Rank(1), true},
		{"InvalidRankHigh", Spades, Rank(15), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCard(tt.s, tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCardString(t *testing.T) {
	c := Card{Suit: Hearts, Rank: Ace}
	if got := c.String(); got != "Ace of Hearts" {
		t.Errorf("String() = %v, want %v", got, "Ace of Hearts")
	}
	if got := c.Code(); got != "Ah" {
		t.Errorf("Code() = %v, want %v", got, "Ah")
	}

	c2 := Card{Suit: Clubs, Rank: Ten}
	if got := c2.Code(); got != "Tc" {
		t.Errorf("Code() = %v, want %v", got, "Tc")
	}
}
