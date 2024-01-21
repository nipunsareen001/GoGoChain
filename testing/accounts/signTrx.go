package testing

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
	"io/ioutil"
	"math/big"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// KeystoreData represents the structure of the keystore file.
type KeystoreData struct {
	Address    string `json:"address"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// DecryptKey decrypts the private key using the user's password.
func DecryptKey(encryptedPrivateKey []byte, password string) (*ecdsa.PrivateKey, error) {
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

// LoadPrivateKeyFromKeystore loads and decrypts the private key from the keystore file.
func LoadPrivateKeyFromKeystore(keystorePath, password string) (*ecdsa.PrivateKey, error) {
	fileData, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		return nil, err
	}

	var keystore KeystoreData
	err = json.Unmarshal(fileData, &keystore)
	if err != nil {
		return nil, err
	}

	encryptedPrivateKey, err := hex.DecodeString(keystore.PrivateKey)
	if err != nil {
		return nil, err
	}

	return DecryptKey(encryptedPrivateKey, password)
}

// SignTransaction signs the transaction with the given private key.
func SignTransaction(privateKey *ecdsa.PrivateKey, transactionData2 []byte) ([]byte, *big.Int, *big.Int, error) {

	transactionData := []byte("data to sign")
	hash := sha256.Sum256(transactionData)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, nil, nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Println("r: ", r)
	fmt.Println("s: ", s)
	return signature, r, s, nil
}

func signTrx(address string) (*big.Int, *big.Int, *ecdsa.PrivateKey, string) {
	// Example usage
	// address := "a04467fb2ef4dcb8d067cc39ae002cfc9423b6a1"
	filename, err := FindKeystoreFileByAddress(address)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Keystore file found for address %s: %s\n", address, filename)
	}

	password := "your-secure-password"

	privateKey, err := LoadPrivateKeyFromKeystore(filename, password)
	if err != nil {
		panic(err)
	}

	transactionData := []byte("data to sign") // Replace with actual transaction data
	signature, r, s, err := SignTransaction(privateKey, transactionData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction signed successfully: %x\n", signature)
	return r, s, privateKey, string(signature)
}

func FindKeystoreFileByAddress(address string) (string, error) {
	keystoreDir := "C:/WorkSpace/GoGoChain/Development/GoGoChain/testing/accounts/keystore"
	files, err := ioutil.ReadDir(keystoreDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") && strings.Contains(file.Name(), address) {
			return keystoreDir + "/" + file.Name(), nil
		}
	}

	return "", fmt.Errorf("no matching .json file found for address: %s", address)
}
