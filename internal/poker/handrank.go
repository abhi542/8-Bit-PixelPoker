package poker

type HandRank int

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (r HandRank) String() string {
	switch r {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

// EvaluatedHand represents the result of a hand evaluation.
type EvaluatedHand struct {
	Rank HandRank
	// TieBreakers contains ranks used to break ties within the same HandRank.
	// Order matters: most significant first.
	// e.g. Full House (Ks over 9s): [King, Nine]
	// e.g. Pair (Jacks, kickers A, 8, 2): [Jack, Ace, Eight, Two]
	TieBreakers []Rank
}

// Compare returns:
// 1 if h > other
// -1 if h < other
// 0 if h == other
func (h EvaluatedHand) Compare(other EvaluatedHand) int {
	if h.Rank > other.Rank {
		return 1
	}
	if h.Rank < other.Rank {
		return -1
	}

	// Determine max length to compare (should be same, but safe code)
	n := len(h.TieBreakers)
	if len(other.TieBreakers) < n {
		n = len(other.TieBreakers)
	}

	for i := 0; i < n; i++ {
		if h.TieBreakers[i] > other.TieBreakers[i] {
			return 1
		}
		if h.TieBreakers[i] < other.TieBreakers[i] {
			return -1
		}
	}
	return 0
}
