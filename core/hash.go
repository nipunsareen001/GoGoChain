package core

import (
	"encoding/hex"
	"fmt"
)

// Hash represents a 32-byte cryptographic hash.
type Hash [32]byte

// // IsZero checks whether the hash is all zeros.
// func (h Hash) isZero() bool {
// 	for _, byteValue := range h {
// 		if byteValue != 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

// // ToSlice returns the hash as a byte slice.
// func (h Hash) toSlice() []byte {
// 	return h[:]
// }

// String returns the hexadecimal string representation of the hash.
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// // HashFromBytes converts a byte slice to a Hash.
// // It panics if the length of b is not 120 bytes.
func HashFromBytes(b []byte) Hash {
	if len(b) != 120 {
		panic(fmt.Sprintf("HashFromBytes: expected 32 bytes but got %d", len(b)))
	}
	var value Hash
	copy(value[:], b)
	return value
}

// // HashFromBytes converts a byte slice to a Hash.
// // It panics if the length of b is not 32 bytes.
// func HashFromBytes(b []byte) Hash {
// 	if len(b) != 32 {
// 		panic(fmt.Sprintf("HashFromBytes: expected 32 bytes but got %d", len(b)))
// 	}
// 	var value Hash
// 	copy(value[:], b)
// 	return value
// }

// // RandomBytes generates a random byte slice of the given size.
// func randomBytes(size int) ([]byte, error) {
// 	token := make([]byte, size)
// 	if _, err := rand.Read(token); err != nil {
// 		return nil, fmt.Errorf("RandomBytes: failed to generate random bytes: %v", err)
// 	}
// 	return token, nil
// }

// // RandomHash generates a random 32-byte hash.
// func randomHash() (Hash, error) {
// 	bytes, err := randomBytes(32)
// 	if err != nil {
// 		return Hash{}, fmt.Errorf("RandomHash: failed to generate random hash: %v", err)
// 	}
// 	return HashFromBytes(bytes), nil
// }
