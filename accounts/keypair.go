package accounts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// KeyPair holds a public/private key pair.
type keyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// EncryptKey encrypts the private key using AES.
func encryptKey(privateKey []byte, password string) ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nonce, nonce, privateKey, nil)
	return append(salt, encrypted...), nil
}

// DecryptKey decrypts the private key using the user's password.
func decryptKey(encryptedPrivateKey []byte, password string) (*ecdsa.PrivateKey, error) {
	salt := encryptedPrivateKey[:16]
	encryptedData := encryptedPrivateKey[16:]

	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	nonce, ciphertext := encryptedData[:gcm.NonceSize()], encryptedData[gcm.NonceSize():]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256() // Use the appropriate elliptic curve
	privateKey.D = new(big.Int).SetBytes(decryptedData)
	privateKey.PublicKey.Curve = privateKey.Curve
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(decryptedData)

	return privateKey, nil
}

// generateKeyPair generates a new ECDSA key pair.
func generateKeyPair() (*keyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := privateKey.D.Bytes()
	pubKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &keyPair{
		PrivateKey: privateKeyBytes,
		PublicKey:  pubKeyBytes,
	}, nil
}

// saveKeyPair saves the encrypted private key, public key, and address to a file.
func saveKeyPair(keyPair *keyPair, password string) (string, string, error) {
	err := createKeystoreDir()
	if err != nil {
		return "", "", err
	}

	encryptedPrivateKey, err := encryptKey(keyPair.PrivateKey, password)
	if err != nil {
		return "", "", err
	}

	address := generateAddress(keyPair.PublicKey)

	data := map[string]interface{}{
		"address":    address,
		"publicKey":  keyPair.PublicKey,
		"privateKey": hex.EncodeToString(encryptedPrivateKey),
	}

	jsonData, err := json.MarshalIndent(data, "", "    ") // 4 spaces indentation
	if err != nil {
		return "", "", err
	}

	filename := generateFilename(address)
	filepath := fmt.Sprintf("keystore/%s", filename)
	err = os.WriteFile(filepath, jsonData, 0600)
	return address, filepath, nil
}

// SignTransaction signs the transaction with the given private key.
func SignTransaction(privateKey *ecdsa.PrivateKey, hash []byte) error {
	_, _, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return err
	}
	return nil
}

// // verifyTransactionSignature verifies the authenticity and integrity of a transaction using the public key.
// func verifyTransactionSignature(pubKeyBytes []byte, signatureHex string, r *big.Int, s *big.Int, pubKey ecdsa.PublicKey, transactionDataHex string) (bool, error) {
// 	signatureBytes, err := hex.DecodeString(signatureHex)
// 	if err != nil {
// 		return false, err
// 	}

// 	transactionData, err := hex.DecodeString(transactionDataHex)
// 	if err != nil {
// 		return false, err
// 	}
// 	hash := sha256.Sum256(transactionData)

// 	xBytes := pubKeyBytes[:32] // Assuming X component is 32 bytes
// 	yBytes := pubKeyBytes[32:] // Assuming Y component is 32 bytes

// 	publicKey := ecdsa.PublicKey{
// 		Curve: elliptic.P256(), // Use the appropriate elliptic curve
// 		X:     new(big.Int).SetBytes(xBytes),
// 		Y:     new(big.Int).SetBytes(yBytes),
// 	}

// 	r_, s_ := new(big.Int).SetBytes(signatureBytes[:32]), new(big.Int).SetBytes(signatureBytes[32:])
// 	if err != nil {
// 		return false, err
// 	}

// 	return ecdsa.Verify(&publicKey, hash[:], r_, s_), nil
// }
