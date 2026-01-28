package poker

import (
	"fmt"
)

type Suit int8
type Rank int8

const (
	// Suits
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

const (
	// Ranks (2-14)
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

// Card represents a single playing card.
type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// NewCard creates a validated card.
func NewCard(s Suit, r Rank) (Card, error) {
	if s < Clubs || s > Spades {
		return Card{}, fmt.Errorf("invalid suit: %d", s)
	}
	if r < Two || r > Ace {
		return Card{}, fmt.Errorf("invalid rank: %d", r)
	}
	return Card{Suit: s, Rank: r}, nil
}

// String implements the fmt.Stringer interface.
// Returns "Rank of Suit" (e.g., "Ace of Spades")
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
}

// Code returns the concise 2-char string (e.g., "As", "Td", "2c")
func (c Card) Code() string {
	return c.Rank.Code() + c.Suit.Code()
}

// Suit helpers
func (s Suit) String() string {
	switch s {
	case Clubs:
		return "Clubs"
	case Diamonds:
		return "Diamonds"
	case Hearts:
		return "Hearts"
	case Spades:
		return "Spades"
	default:
		return "?"
	}
}

func (s Suit) Code() string {
	switch s {
	case Clubs:
		return "c"
	case Diamonds:
		return "d"
	case Hearts:
		return "h"
	case Spades:
		return "s"
	default:
		return "?"
	}
}

// Rank helpers
func (r Rank) String() string {
	switch r {
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Ten:
		return "Ten"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	default:
		return "?"
	}
}

func (r Rank) Code() string {
	switch r {
	case Two, Three, Four, Five, Six, Seven, Eight, Nine:
		return fmt.Sprintf("%d", r)
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return "?"
	}
}
