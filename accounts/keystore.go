package accounts

import (
	"GoGoChain/enum"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"strings"
)

const (
	keystoreFolderPath = string(enum.KeystorePath)
)

// KeystoreData represents the structure of the keystore file.
type keystoreData struct {
	Address    string `json:"address"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// CreateKeystoreDir creates the keystore directory if it does not exist.
func createKeystoreDir() error {
	path := keystoreFolderPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0700)
	}
	return nil
}

// LoadPrivateKeyFromKeystore loads and decrypts the private key from the keystore file.
func LoadPrivateKeyFromKeystore(keystorePath, password string) (string, *ecdsa.PrivateKey, error) {
	fileData, err := os.ReadFile(keystorePath)
	if err != nil {
		return "", nil, err
	}

	var keystore keystoreData
	err = json.Unmarshal(fileData, &keystore)
	if err != nil {
		return "", nil, err
	}

	encryptedPrivateKey, err := hex.DecodeString(keystore.PrivateKey)
	if err != nil {
		return "", nil, err
	}
	privkey, err := decryptKey(encryptedPrivateKey, password)
	if err != nil {
		return "", nil, fmt.Errorf("Password doest not match!")
	}

	return keystore.PublicKey, privkey, err
}

// findKeystoreFileByAddress finds the keystore file with the matching address.
func FindKeystoreFileByAddress(address string) (string, error) {
	files, err := os.ReadDir(keystoreFolderPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") && strings.Contains(file.Name(), address) {
			return keystoreFolderPath + "/" + file.Name(), nil
		}
	}

	return "", fmt.Errorf("no matching .json file found for address: %s", address)
}

// NewAccount generates a new account, encrypts it, and stores it in the keystore.
func NewAccount() {
	fmt.Println("Enter a password for the new account:") // password for the address/account
	var password string
	fmt.Scanln(&password)

	keyPair, err := generateKeyPair()
	if err != nil {
		panic(err)
	}

	address, filepath, err := saveKeyPair(keyPair, password)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Key pair generated and saved successfully.\n Address: %s\n Filepath: %s\n", address, filepath)

}
