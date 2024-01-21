package block

import (
	"GoGoChain/core"
	"GoGoChain/enum"
	"GoGoChain/leveldb"
	"crypto/sha256"
	"fmt"
)

type block struct {
	BlockNumber  int      `json:"blocknumber"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
	PrevHash     string   `json:"prevHash"`
	Hash         string   `json:"hash,omitempty"`
	// CONSENSUS RELATED
	// ValidatorAddresses []string `json:"validatorAddresses"`
}

// BlockHash calculates the hash of the block including the timestamp.
func blockHash(b *block) core.Hash {
	txInByte := []byte(fmt.Sprintf("%d-%s-%s-%s", b.BlockNumber, b.Timestamp, b.Transactions, b.PrevHash))
	hash := sha256.Sum256(txInByte)
	return hash // Ensure this conversion is correct based on 'Hash' type definition
}

// GetBalance retrieves the balance for the given address.
func GetBalance(address string) (string, error) {
	balanceInBytes := leveldb.FetchFromDatabase(string(enum.BlockstateDb), []byte(address))
	if balanceInBytes == nil {
		return "", fmt.Errorf("address not found: %s", address)
	}

	balance := (string(balanceInBytes))

	return balance, nil
}

// GetCurrentBlockNumber retrieves the current block number.
func GetCurrentBlockNumber() (string, error) {
	currentBlockNumber := leveldb.FetchFromDatabase(string(enum.BlockchianDb), []byte(string(enum.CurrentBlockNumber)))
	if currentBlockNumber == nil {
		return "", fmt.Errorf("current block number not found")
	}

	return string(currentBlockNumber), nil
}

// GetBlockHashByNumber retrieves the block hash for the given block number.
func GetBlockHashByNumber(blockNumber string) (string, error) {
	blockHash := leveldb.FetchFromDatabase(string(enum.BlockchianDb), []byte(blockNumber))
	if blockHash == nil {
		return "", fmt.Errorf("block not found for number: %s", blockNumber)
	}

	return string(blockHash), nil
}

// GetBlockByHash retrieves the block details for the given block hash.
func GetBlockByHash(blockHash string) (string, error) {
	blockDetails := leveldb.FetchFromDatabase(string(enum.BlockchianDb), []byte(blockHash))
	if blockDetails == nil {
		return "", fmt.Errorf("block not found for hash: %s", blockHash)
	}

	return string(blockDetails), nil
}

// GetTransactionByHash retrieves the transaction details for the given transaction hash.
func GetTransactionByHash(hashq string) (string, error) {
	transaction := leveldb.FetchFromDatabase(string(enum.BlockchianDb), []byte(hashq))
	if transaction == nil {
		return "", fmt.Errorf("transaction not found for hash: %s", hashq)
	}

	return string(transaction), nil
}
