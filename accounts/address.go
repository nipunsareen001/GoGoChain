package accounts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// generateAddress generates an Ethereum-style address from a public key.
func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[len(hash)-20:]) // Use the last 20 bytes
}

// generateFilename generates a filename in Ethereum keystore format.
func generateFilename(address string) string {
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z")
	return fmt.Sprintf("UTC--%s--%s.json", timestamp, address)
}
