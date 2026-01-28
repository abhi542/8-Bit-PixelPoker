package poker

import (
	"errors"
)

var (
	ErrDeckEmpty = errors.New("deck is empty")
)

// Deck represents a collection of cards.
// We expose it as a slice alias for easy iteration,
// but provide methods for safe interaction.
type Deck []Card

// NewDeck creates a fresh 52-card deck in sorted order:
// Clubs (2-A), Diamonds (2-A), Hearts (2-A), Spades (2-A).
// It does NOT shuffle by default.
func NewDeck() Deck {
	d := make(Deck, 0, 52)
	// Order: Clubs, Diamonds, Hearts, Spades
	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	for _, s := range suits {
		for r := Two; r <= Ace; r++ {
			d = append(d, Card{Suit: s, Rank: r})
		}
	}
	return d
}

// Draw removes the top n cards from the deck and returns them.
// "Top" is defined as index 0.
func (d *Deck) Draw(n int) ([]Card, error) {
	if len(*d) < n {
		return nil, ErrDeckEmpty
	}
	drawn := (*d)[:n]
	*d = (*d)[n:]
	return drawn, nil
}

// Burn removes the top card and returns it.
// This is used to "burn" a card before dealing streets.
func (d *Deck) Burn() (Card, error) {
	if len(*d) < 1 {
		return Card{}, ErrDeckEmpty
	}
	burned := (*d)[0]
	*d = (*d)[1:]
	return burned, nil
}

// Shuffle implements the Fisher-Yates shuffle using the provided RNG.
func (d Deck) Shuffle(rng RNG) {
	n := len(d)
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
}

// Len returns remaining cards
func (d Deck) Len() int {
	return len(d)
}
