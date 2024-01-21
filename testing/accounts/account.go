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
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// KeyPair holds a public/private key pair.
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// EncryptKey encrypts the private key using AES.
func EncryptKey(privateKey []byte, password string) ([]byte, error) {
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

// GenerateEthereumAddress generates an Ethereum-style address from a public key.
func GenerateEthereumAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[len(hash)-20:]) // Use the last 20 bytes
}

// CreateKeystoreDir creates the keystore directory if it does not exist.
func CreateKeystoreDir() error {
	path := "keystore"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0700)
	}
	return nil
}

// GenerateEthereumStyleFilename generates a filename in Ethereum keystore format.
func GenerateEthereumStyleFilename(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	address := hex.EncodeToString(hash[:])
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z")
	return fmt.Sprintf("UTC--%s--%s.json", timestamp, address)
}

// GenerateKeyPair generates a new ECDSA key pair.
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := privateKey.D.Bytes()
	// publicKeyBytes := elliptic.MarshalCompressed(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
	pubKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &KeyPair{
		PrivateKey: privateKeyBytes,
		PublicKey:  pubKeyBytes,
	}, nil
}

// SaveKeyPair saves the encrypted private key, public key, and address to a file.
func SaveKeyPair(keyPair *KeyPair, password string) (string, []byte, error) {
	err := CreateKeystoreDir()
	if err != nil {
		return "", nil, err
	}

	encryptedPrivateKey, err := EncryptKey(keyPair.PrivateKey, password)
	if err != nil {
		return "", nil, err
	}

	address := GenerateEthereumAddress(keyPair.PublicKey)

	// data := map[string]string{
	// 	"address":    address,
	// 	"publicKey":  hex.EncodeToString(keyPair.PublicKey),
	// 	"privateKey": hex.EncodeToString(encryptedPrivateKey),
	// }

	// data := map[string][]byte{
	// 	"address":    []byte(address),
	// 	"publicKey":  (keyPair.PublicKey),
	// 	"privateKey": (encryptedPrivateKey),
	// }

	data := map[string]interface{}{
		"address":    address,
		"publicKey":  keyPair.PublicKey,
		"privateKey": hex.EncodeToString(encryptedPrivateKey),
	}

	// Use MarshalIndent for pretty JSON formatting
	jsonData, err := json.MarshalIndent(data, "", "    ") // 4 spaces indentation
	if err != nil {
		return "", nil, err
	}

	filename := fmt.Sprintf("UTC--%s--%s.json", time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z"), address)
	filepath := fmt.Sprintf("keystore/%s", filename)
	err = ioutil.WriteFile(filepath, jsonData, 0600)
	return address, keyPair.PublicKey, err
}

func creatNewAddress() (string, []byte) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	password := "your-secure-password"
	address, pubKey, err := SaveKeyPair(keyPair, password)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Key pair generated and saved successfully. Address: %s\n", address)

	return address, pubKey
}
