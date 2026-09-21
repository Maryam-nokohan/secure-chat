package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GeneratePublicIDCandidate() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}