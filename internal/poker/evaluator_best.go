package poker

// EvaluateBestHand returns the best 5-card hand from a set of 5, 6, or 7 cards.
// For < 5 cards (PreFlop), it returns High Card or Pair based on hole cards (if we want that).
// Actually, standard Hold'em "Current Hand" usually implies "Best 5 cards available".
// PreFlop (2 cards): Technically "High Card", but often UI shows "Pocket Pair" or "High Card".
// We will support 2..7 cards to be safe.
func EvaluateBestHand(cards []Card) EvaluatedHand {
	n := len(cards)
	if n < 5 {
		// Just evaluate these cards (2, 3, 4) roughly?
		// Or return specific "pre-flop" status.
		// Detailed evaluator expects 5.
		// Let's just return a placeholder or simple High Card logic if < 5.
		// Actually, Evaluate5Card crashes if != 5.
		// For MVP, if < 5 (PreFlop), we might just say "High Card" of the hole cards.
		// Let's handle n >= 5 properly.
		return evaluatePartial(cards)
	}

	if n == 5 {
		ev, _ := Evaluate5Card(cards)
		return ev
	}

	if n == 7 {
		ev, _ := Evaluate7Card(cards) // Already implemented
		return ev
	}

	// Case n == 6 (Turn)
	// Combinations of 6 choose 5 = 6 combinations.
	// We can reuse Evaluate7Card logic but for 6.
	// Or just write a quick helper here.
	return evaluateCombination(cards, 5)
}

func evaluateCombination(cards []Card, k int) EvaluatedHand {
	// Generate combinations of k cards from len(cards)
	// Naive approach for n=6, k=5 is trivial.

	var best EvaluatedHand

	best.Rank = -1 // Initialize to invalid to ensure first valid hand updates it

	indices := make([]int, k)
	for i := range indices {
		indices[i] = i
	}

	n := len(cards)

	for {
		// Build hand
		hand := make([]Card, k)
		for i, idx := range indices {
			hand[i] = cards[idx]
		}

		ev, _ := Evaluate5Card(hand) // internal helper, we trust 5 cards
		if ev.Compare(best) > 0 {
			best = ev
		}

		// Next combo
		i := k - 1
		for i >= 0 && indices[i] == i+n-k {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < k; j++ {
			indices[j] = indices[j-1] + 1
		}
	}

	return best
}

func evaluatePartial(cards []Card) EvaluatedHand {
	// For 2 cards (Hole cards), we can check for Pair.
	// Otherwise High Card.
	if len(cards) == 2 {
		if cards[0].Rank == cards[1].Rank {
			return EvaluatedHand{
				Rank:        Pair,
				TieBreakers: []Rank{cards[0].Rank}, // Simplified
			}
		}
		// High card
		h := cards[0].Rank
		if cards[1].Rank > h {
			h = cards[1].Rank
		}
		return EvaluatedHand{
			Rank:        HighCard,
			TieBreakers: []Rank{h},
		}
	}
	return EvaluatedHand{Rank: HighCard}
}
