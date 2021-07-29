package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func NewSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
