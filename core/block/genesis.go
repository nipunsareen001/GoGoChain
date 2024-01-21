package block

import (
	"GoGoChain/core/transaction"
	"GoGoChain/enum"
	"GoGoChain/leveldb"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
)

type allocJSON struct {
	Allocs []alloc `json:"alloc"`
}

type alloc struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}

// InitGenesis initializes the genesis block with predefined allocations with geneis.json file in root directory.
func InitGenesis() {
	allocJSON, err := readGenesis(string(enum.GenesisFilepath))
	if err != nil {
		fmt.Println("Failed to read genesis file:", err)
		return
	}

	var transactionHashes []string
	for _, alloc := range allocJSON.Allocs {
		to := alloc.Address
		value := alloc.Balance
		hashInString, trxBody := transaction.SignGenesisTrx(to, value)

		// Add transaction to the BlockchianDb
		leveldb.AddDataToDatabase((string(enum.BlockchianDb)), []byte(hashInString), []byte(trxBody))

		// Update BlockstateDb
		leveldb.AddDataToDatabase((string(enum.BlockstateDb)), []byte(to), []byte(value.Text(10)))

		transactionHashes = append(transactionHashes, hashInString)
	}

	blockHashInString, blockBody := genesisBlockHash(transactionHashes)

	// Add block 0 data to BlockchianDb
	leveldb.AddDataToDatabase((string(enum.BlockchianDb)), []byte(blockHashInString), []byte(blockBody))

	// Add blocknumber 0 with block 0 hash in BlockchianDb
	blockNumber := (string(enum.GenesisBlockNumber))
	leveldb.AddDataToDatabase((string(enum.BlockchianDb)), []byte(blockNumber), []byte(blockHashInString))

	// Add lasthash represented by 'lh' with block 0 hash in BlockchianDb
	lastHash := (string(enum.LastHash))
	leveldb.AddDataToDatabase((string(enum.BlockchianDb)), []byte(lastHash), []byte(blockHashInString))

	// Update blocknumber as block 0  in BlockchianDb
	currentBlockNumber := (string(enum.CurrentBlockNumber))
	leveldb.AddDataToDatabase((string(enum.BlockchianDb)), []byte(currentBlockNumber), []byte(blockNumber))

	// SHARDING RELATED, look after words
	// // Update current block extended number
	// currentBlockExtendedNumber := string(enum.CurrentBlockExtendedNumber)
	// db.AddDataToDatabase(string(enum.Blockextended), []byte(currentBlockExtendedNumber), []byte(blockNumber))

	// executedBlockExtendedNumber := string(enum.ExecutedBlockExtendedNumber)
	// db.AddDataToDatabase(string(enum.Blockextended), []byte(executedBlockExtendedNumber), []byte(blockNumber))
}

// readGenesis reads and parses the genesis file.
func readGenesis(filename string) (allocJSON, error) {
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return allocJSON{}, fmt.Errorf("error reading genesis file: %w", err)
	}

	var allocJSON_ allocJSON
	if err := json.Unmarshal(byteValue, &allocJSON_); err != nil {
		return allocJSON{}, fmt.Errorf("error parsing genesis file: %w", err)
	}

	return allocJSON_, nil
}

// GenesisBlockHash generates a hash for the genesis block and returns it along with the block's JSON representation.
func genesisBlockHash(transactions []string) (string, string) {
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"
	timestamp := "123654789"
	block_ := &block{
		BlockNumber:  0,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     prevHash,
	}

	blockHash := blockHash(block_)
	blockAfterHash := &block{
		BlockNumber:  0,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         blockHash.String(),
	}

	blockBody, err := json.MarshalIndent(blockAfterHash, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal block: %v", err)
	}

	fmt.Println(string(blockBody))

	return blockHash.String(), string(blockBody)
}
