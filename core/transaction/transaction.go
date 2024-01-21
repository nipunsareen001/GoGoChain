package transaction

import (
	"GoGoChain/accounts"
	"GoGoChain/core"
	"GoGoChain/enum"
	"GoGoChain/filemanagement"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

type transaction struct {
	To        string   `json:"to"`
	From      string   `json:"from"`
	Value     *big.Int `json:"value"`
	Timestamp int64    `json:"timestamp"`
	Hash      string   `json:"hash,omitempty"`
	PublicKey string   `json:"publickey,omitempty"`
}

// SignGenesisTrx creates and signs a genesis transaction, returning its hash and JSON representation.
func SignGenesisTrx(_to string, _value *big.Int) (string, string) {
	timestamp := int64(123654789)
	_from := "00000000000000000GENESIS0000000000000000"

	tx := &transaction{
		To:        _to,
		From:      _from,
		Value:     _value,
		Timestamp: timestamp,
	}

	txhash := transactionHash(tx)

	// transactionInBytes := []byte(fmt.Sprintf("%v", tx))
	// txHash := core.HashFromBytes(transactionInBytes)
	// txxHash := txHash.String()

	txAfterHash := &transaction{
		To:        _to,
		From:      _from,
		Value:     _value,
		Timestamp: timestamp,
		Hash:      txhash.String(),
	}

	transactionBody, err := json.MarshalIndent(txAfterHash, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}

	return txhash.String(), string(transactionBody)
}

// SendTransaction handles the entire process of sending a transaction
func SendTransaction(from string, to string, value *big.Int, password string) error {
	filename, err := accounts.FindKeystoreFileByAddress(from)
	if err != nil {
		return err
	}

	publicKey, privateKey, err := accounts.LoadPrivateKeyFromKeystore(filename, password)
	if err != nil {
		return err
	}

	txBody, hash := createAndHashTransaction(publicKey, from, to, value)

	err = accounts.SignTransaction(privateKey, hash[:])
	if err != nil {
		return err
	}

	return appendTransactionToPool(txBody)
}

// createAndHashTransaction creates a transaction and calculates its hash
func createAndHashTransaction(publicKey string, from string, to string, value *big.Int) (*transaction, core.Hash) {
	timestamp := time.Now().Unix()
	tx := &transaction{
		To:        to,
		From:      from,
		Value:     value,
		Timestamp: timestamp,
		PublicKey: publicKey,
	}
	hash := transactionHash(tx) // Calculating transaction hash
	tx.Hash = hash.String()
	return tx, hash
}

func transactionHash(tx *transaction) core.Hash {
	txInByte := []byte(fmt.Sprintf("%s-%s-%d-%d", tx.From, tx.To, tx.Value, tx.Timestamp))
	hash := sha256.Sum256(txInByte)
	return hash
}

// appendTransactionToPool appends the transaction to the transaction pool (CSV file in this case)
func appendTransactionToPool(tx *transaction) error {

	file, err := filemanagement.OpenFileForAppending(string(enum.TransactionMemoryPool))
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("Error marshaling transaction: %v", err)
	}

	err = filemanagement.AppendDataInFile(file, content)
	if err != nil {
		return err
	}

	return nil
}
