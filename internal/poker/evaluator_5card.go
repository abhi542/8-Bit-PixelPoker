package poker

import (
	"fmt"
	"sort"
)

// Evaluate5Card computes the rank of exactly 5 cards.
// It assumes the input slice has exactly 5 cards.
func Evaluate5Card(hand []Card) (EvaluatedHand, error) {
	if len(hand) != 5 {
		return EvaluatedHand{}, fmt.Errorf("expected 5 cards, got %d", len(hand))
	}

	// 1. Sort cards by Rank (Descending) for easier analysis
	// Copy to avoid mutating the original slice
	sorted := make([]Card, 5)
	copy(sorted, hand)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank > sorted[j].Rank
	})

	isFlush := checkFlush(sorted)
	isStraight, highRank := checkStraight(sorted)

	// Straight Flush / Royal Flush
	if isFlush && isStraight {
		if highRank == Ace {
			// Check if it's actually an A-5 straight flush (Ace is high rank, but 5 is top for SF unless it's TJQKA)
			// Wait, checkStraight returns the *effective* high rank of the straight.
			// If it's TJQKA, highRank is Ace.
			// If it's A2345, highRank is Five.
			return EvaluatedHand{Rank: RoyalFlush, TieBreakers: []Rank{Ace}}, nil
		}
		return EvaluatedHand{Rank: StraightFlush, TieBreakers: []Rank{highRank}}, nil
	}

	// Four of a Kind
	// Patterns: AAAA B or B AAAA
	quads, kickers := getMultiples(sorted, 4)
	if len(quads) > 0 {
		return EvaluatedHand{Rank: FourOfAKind, TieBreakers: []Rank{quads[0], kickers[0]}}, nil
	}

	// Full House
	// Patterns: AAA BB or AA BBB (sorted)
	// Let's implement specific frequency map logic for clarity.
	counts := make(map[Rank]int)
	for _, c := range sorted {
		counts[c.Rank]++
	}

	var tripRank, pairRank Rank
	hasTrips := false
	hasPair := false

	// Iterate 14 down to 2 to find highest trips/pairs
	for r := Ace; r >= Two; r-- {
		if counts[r] == 3 {
			tripRank = r
			hasTrips = true
		} else if counts[r] == 2 {
			// If we already have a pair (e.g. 2 pair), this loop finds the higher one first?
			// But wait, full house can only comprise 3+2. 5 cards total.
			// So there is at most one set of 3 and one set of 2.
			pairRank = r
			hasPair = true
		}
	}

	if hasTrips && hasPair {
		return EvaluatedHand{Rank: FullHouse, TieBreakers: []Rank{tripRank, pairRank}}, nil
	}

	// Flush
	if isFlush {
		// Tiebreakers are all 5 cards
		tb := make([]Rank, 5)
		for i, c := range sorted {
			tb[i] = c.Rank
		}
		return EvaluatedHand{Rank: Flush, TieBreakers: tb}, nil
	}

	// Straight
	if isStraight {
		return EvaluatedHand{Rank: Straight, TieBreakers: []Rank{highRank}}, nil
	}

	// Three of a Kind
	if hasTrips {
		// Tiebreakers: TripRank, then 2 kickers
		tb := []Rank{tripRank}
		for _, c := range sorted {
			if c.Rank != tripRank {
				tb = append(tb, c.Rank)
			}
		}
		return EvaluatedHand{Rank: ThreeOfAKind, TieBreakers: tb}, nil
	}

	// Two Pair
	// We need to scan carefully for 2 pairs.
	// Since 'sorted' is high-to-low, we can check adjacent cards or counts.
	// We used map above, but maps are unordered.
	// Let's re-scan counts for Two Pair logic.
	firstPair := Rank(0)
	secondPair := Rank(0)
	for r := Ace; r >= Two; r-- {
		if counts[r] == 2 {
			if firstPair == 0 {
				firstPair = r
			} else {
				secondPair = r
				break
			}
		}
	}

	if firstPair != 0 && secondPair != 0 {
		// Found two pair
		var kicker Rank
		for r := Ace; r >= Two; r-- {
			if counts[r] == 1 {
				kicker = r
				break
			}
		}
		return EvaluatedHand{Rank: TwoPair, TieBreakers: []Rank{firstPair, secondPair, kicker}}, nil
	}

	// Pair
	if firstPair != 0 {
		var kickers []Rank
		for _, c := range sorted {
			if c.Rank != firstPair {
				kickers = append(kickers, c.Rank)
			}
		}
		return EvaluatedHand{Rank: Pair, TieBreakers: append([]Rank{firstPair}, kickers...)}, nil
	}

	// High Card
	tb := make([]Rank, 5)
	for i, c := range sorted {
		tb[i] = c.Rank
	}
	return EvaluatedHand{Rank: HighCard, TieBreakers: tb}, nil
}

// Helpers

func checkFlush(sorted []Card) bool {
	s := sorted[0].Suit
	for i := 1; i < 5; i++ {
		if sorted[i].Suit != s {
			return false
		}
	}
	return true
}

func checkStraight(sorted []Card) (bool, Rank) {
	// Standard check: Descending ranks, difference of 1
	isStandard := true
	for i := 0; i < 4; i++ {
		if sorted[i].Rank != sorted[i+1].Rank+1 {
			isStandard = false
			break
		}
	}
	if isStandard {
		return true, sorted[0].Rank
	}

	// Wheel Check: A, 5, 4, 3, 2
	// Sorted: A, 5, 4, 3, 2
	if sorted[0].Rank == Ace &&
		sorted[1].Rank == Five &&
		sorted[2].Rank == Four &&
		sorted[3].Rank == Three &&
		sorted[4].Rank == Two {
		return true, Five // 5-high straight
	}

	return false, 0
}

// getMultiples finds 'count' of same rank. Returns ranks found and remaining kickers.
// Note: This is a helper, but for 5-card strict eval, the map approach in main function covers it.
func getMultiples(sorted []Card, count int) ([]Rank, []Rank) {
	// Not strictly used by main logic anymore since I switched to map,
	// but useful structure for finding Quads.
	counts := make(map[Rank]int)
	for _, c := range sorted {
		counts[c.Rank]++
	}

	var matches []Rank
	var kickers []Rank

	// Collect matches
	for r, c := range counts {
		if c == count {
			matches = append(matches, r)
		}
	}
	// Sort matches descending (e.g. if looking for pairs and we have two pairs)
	sort.Slice(matches, func(i, j int) bool { return matches[i] > matches[j] })

	// Collect kickers (cards not in matches)
	// We iterate sorted array to keep kickers sorted high-to-low
	for _, c := range sorted {
		isMatch := false
		for _, m := range matches {
			if c.Rank == m {
				isMatch = true
				break
			}
		}
		if !isMatch {
			kickers = append(kickers, c.Rank)
		}
	}

	return matches, kickers
}
