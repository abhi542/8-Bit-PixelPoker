package poker

import (
	"fmt"
)

// Evaluate7Card finds the best 5-card hand from 7 cards.
// It iterates through all 21 combinations (7 choose 5).
func Evaluate7Card(cards []Card) (EvaluatedHand, error) {
	if len(cards) != 7 {
		return EvaluatedHand{}, fmt.Errorf("expected 7 cards, got %d", len(cards))
	}

	var bestHand EvaluatedHand
	first := true

	// 7 choose 5 is equivalent to 7 choose 2 (exclude 2 cards).
	// We iterate through all pairs of indices to exclude.
	for i := 0; i < 7; i++ {
		for j := i + 1; j < 7; j++ {
			// Construct 5-card hand by skipping i and j
			hand := make([]Card, 0, 5)
			for k := 0; k < 7; k++ {
				if k == i || k == j {
					continue
				}
				hand = append(hand, cards[k])
			}

			evaluated, err := Evaluate5Card(hand)
			if err != nil {
				return EvaluatedHand{}, err
			}

			if first {
				bestHand = evaluated
				first = false
			} else {
				if evaluated.Compare(bestHand) > 0 {
					bestHand = evaluated
				}
			}
		}
	}

	return bestHand, nil
}
