package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateSessionToken generates a random session token
func GenerateSessionToken() string {
	b := make([]byte, 16) // 128-bit
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate session token: %v", err)
	}
	return hex.EncodeToString(b)
}
