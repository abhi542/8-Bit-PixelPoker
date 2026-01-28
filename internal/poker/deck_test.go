package poker

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := NewDeck()
	if len(d) != 52 {
		t.Fatalf("NewDeck() size = %d, want 52", len(d))
	}

	// Verify order: First is 2c, Last is As
	first := d[0]
	if first.Code() != "2c" {
		t.Errorf("First card = %v, want 2c", first.Code())
	}
	last := d[51]
	if last.Code() != "As" {
		t.Errorf("Last card = %v, want As", last.Code())
	}
}

func TestDeck_Draw(t *testing.T) {
	d := NewDeck()
	drawn, err := d.Draw(5)
	if err != nil {
		t.Fatalf("Draw(5) error = %v", err)
	}

	if len(drawn) != 5 {
		t.Errorf("Drawn count = %d, want 5", len(drawn))
	}
	if len(d) != 47 {
		t.Errorf("Remaining deck = %d, want 47", len(d))
	}

	// Check containment
	if drawn[0].Code() != "2c" {
		t.Errorf("First drawn = %v, want 2c (since no shuffle)", drawn[0].Code())
	}
}

func TestDeck_DrawOverLimit(t *testing.T) {
	d := NewDeck()
	_, err := d.Draw(53)
	if err == nil {
		t.Error("Draw(53) expected error, got nil")
	}
	if err != ErrDeckEmpty {
		t.Errorf("Error = %v, want ErrDeckEmpty", err)
	}
}

func TestDeck_Burn(t *testing.T) {
	d := NewDeck()
	// Top card is 2c (unshuffled)
	burned, err := d.Burn()
	if err != nil {
		t.Fatalf("Burn() error = %v", err)
	}
	if burned.Code() != "2c" {
		t.Errorf("Burned = %v, want 2c", burned.Code())
	}
	if len(d) != 51 {
		t.Errorf("Remaining = %d, want 51", len(d))
	}
}

func TestDeck_Shuffle_Deterministic(t *testing.T) {
	d1 := NewDeck()
	d2 := NewDeck()

	// Seed 42
	rng1 := NewDeterministicRNG(42)
	rng2 := NewDeterministicRNG(42)

	d1.Shuffle(rng1)
	d2.Shuffle(rng2)

	// Verify exact match
	for i := 0; i < 52; i++ {
		if d1[i] != d2[i] {
			t.Errorf("Mismatch at index %d: %v vs %v", i, d1[i], d2[i])
		}
	}

	// Verify it actually changed (probabilistic but seed 42 is known)
	// Just check first card is NOT 2c
	if d1[0].Code() == "2c" {
		t.Errorf("Deck did not appear to shuffle (first card is still 2c)")
	}
}
