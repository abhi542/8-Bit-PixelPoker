package lobby

import (
	"crypto/rand"
	"math/big"
)

// Charset excludes ambiguous characters (0/O, 1/I).
const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func GenerateCode() string {
	b := make([]byte, 4)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
