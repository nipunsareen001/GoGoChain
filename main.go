package main

import (
	accounts "GoGoChain/accounts"
	"GoGoChain/core/block"
	"GoGoChain/core/transaction"
	"GoGoChain/node"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func main() {
	newAccount := flag.Int("newAccount", 0, "Creates a new account if set to 1")
	initGenesis := flag.Int("initGenesis", 0, "Initializes the genesis block if set to 1")
	balance := flag.String("balance", "", "Gets the balance of the provided address")
	blockNumber := flag.Int("blockNumber", 0, "Prints the current block number if set to 1")
	getBlockHashByNumber := flag.String("getBlockHashByNumber", "", "Gets the block hash by block number")
	getBlockByHash := flag.String("getBlockByHash", "", "Gets block details by block hash")
	getTransactionByHash := flag.String("getTransactionByHash", "", "Gets transaction details by transaction hash")
	sendTrx := flag.Int("sendTrx", 0, "sends transaction")
	from := flag.String("from", "", "from address for transaction")
	to := flag.String("to", "", "to address for transaction")
	value := flag.String("value", "0", "value for transaction")
	auth := flag.String("auth", "", "authorization string")
	startNode := flag.Int("startNode", 0, "starts node")
	sourcePort := flag.Int("sourcePort", 3001, "source port number")
	validator := flag.Int("validator", 0, "is validator")
	validator_address := flag.String("address", "", "validator address")

	// Create a set of all defined flags for checking
	definedFlags := make(map[string]struct{})
	flag.VisitAll(func(f *flag.Flag) {
		definedFlags[f.Name] = struct{}{}
	})

	// Check for undefined flags
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flagName := strings.TrimPrefix(arg, "-")
			if _, defined := definedFlags[flagName]; !defined {
				fmt.Printf("Oops! It looks like '%s' is not a recognized option.\n\n", flagName)
				fmt.Println("Here are some options you can use:")
				flag.PrintDefaults()
				fmt.Println("\nPlease check the options and try again.")
				os.Exit(1)
			}
		}
	}

	flag.Parse()

	if *newAccount == 1 {
		accounts.NewAccount()
	}

	if *initGenesis == 1 {
		block.InitGenesis()
		fmt.Println("Genesis has been loaded successfully")
	}

	if *balance != "" {
		balanceOfAddress, _ := block.GetBalance(*balance)
		fmt.Println("Balance of address is:", balanceOfAddress)
	}

	if *blockNumber == 1 {
		blockNumber, _ := block.GetCurrentBlockNumber()
		fmt.Println("Current block number is:", blockNumber)
	}

	if *getBlockHashByNumber != "" {
		blockHash, _ := block.GetBlockHashByNumber(*getBlockHashByNumber)
		fmt.Println("Block hash is:", blockHash)
	}

	if *getBlockByHash != "" {
		b, _ := block.GetBlockByHash(*getBlockByHash)
		fmt.Println("Block is:", b)
	}

	if *getTransactionByHash != "" {
		transaction, _ := block.GetTransactionByHash(*getTransactionByHash)
		fmt.Println("Transaction is:", transaction)
	}

	if *sendTrx == 1 {
		// Verifying 'from', 'to', 'value' and 'auth'
		if !isValidHexAddress(*from) {
			fmt.Println("Error: 'from' address is not valid")
			os.Exit(1)
		}
		if !isValidHexAddress(*to) {
			fmt.Println("Error: 'to' address is not valid")
			os.Exit(1)
		}
		if len(*auth) == 0 {
			fmt.Println("Error: Authorization token is required")
			os.Exit(1)
		}

		// Convert string value to *big.Int
		valueBigInt := new(big.Int)
		_, ok := valueBigInt.SetString(*value, 10) // base 10
		if !ok {
			fmt.Println("Invalid value for transaction")
			os.Exit(1)
		}
		// Compare the valueBigInt with 0. Cmp returns -1 if valueBigInt < 0, 0 if valueBigInt == 0, and +1 if valueBigInt > 0
		if valueBigInt.Cmp(big.NewInt(0)) != 1 {
			fmt.Println("Error: Transaction value must be greater than 0")
			os.Exit(1)
		}

		// If all parameters are valid, send the transaction
		fmt.Println("Sending transaction with the following details:")
		fmt.Printf("From: %s, To: %s, Value: %s\n", *from, *to, *value)
		err := transaction.SendTransaction(*from, *to, valueBigInt, *auth)
		if err != nil {
			fmt.Printf("Transaction cant be sent: %v\n", err)
			os.Exit(1)
		}
	}

	if *startNode == 1 {
		canStartNode := true

		// Perform the address validation only if the node is a validator
		if *validator == 1 {
			if *validator_address == "" || !isValidHexAddress(*validator_address) {
				fmt.Println("Enter a valid, non-empty hexadecimal validator address")
				// If the address is invalid, set the flag to false
				canStartNode = false
			}
		}

		// If all conditions are met, start the node
		if canStartNode {
			node.StartNode(*sourcePort, *validator, *validator_address)
		}
	}
}

func isValidHexAddress(address string) bool {
	// Define a regular expression for a 40-character hexadecimal string.
	re := regexp.MustCompile(`^(0x)?[0-9a-fA-F]{40}$`)
	return re.MatchString(address)
}
