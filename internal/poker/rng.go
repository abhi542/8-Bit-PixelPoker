package poker

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
)

// RNG defines the source of randomness for shuffling.
type RNG interface {
	// Intn returns a non-negative random number in [0, n).
	Intn(n int) int
}

// CryptoRNG uses crypto/rand for secure number generation.
type CryptoRNG struct{}

func (c CryptoRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		// In a poker server, RNG failure is a critical panic.
		panic("crypto/rand failure: " + err.Error())
	}
	return int(v.Int64())
}

// DeterministicRNG uses math/rand for repeatable tests.
type DeterministicRNG struct {
	Source *mathrand.Rand
}

func NewDeterministicRNG(seed int64) *DeterministicRNG {
	return &DeterministicRNG{
		Source: mathrand.New(mathrand.NewSource(seed)),
	}
}

func (d *DeterministicRNG) Intn(n int) int {
	return d.Source.Intn(n)
}
