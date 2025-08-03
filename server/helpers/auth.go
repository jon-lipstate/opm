package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

// GenerateState creates a secure state token for OAuth flows
func GenerateState(secret string) string {
	// Create a unique state with timestamp
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("state-%d", timestamp)
	
	// Create HMAC signature
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	
	// Combine data and signature
	return fmt.Sprintf("%s.%s", data, signature)
}

// ValidateState verifies the state token
func ValidateState(state, secret string) bool {
	// For now, just check it's not empty
	// TODO: Implement proper validation with HMAC
	return state != ""
}