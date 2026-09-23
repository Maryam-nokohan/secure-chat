package pkg

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateNumericCode(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	s := n.String()
	return strings.Repeat("0", digits-len(s)) + s, nil
}