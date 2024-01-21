# GoGoChain - Blockchain Implementation in Go
GoGoChain: A robust, Go-based blockchain implementation designed from scratch. This project features a custom-built blockchain architecture, focusing on security, efficient data storage, and network synchronization. Ideal for understanding core blockchain principles and exploring innovative blockchain solutions.

## Introduction
GoGoChain is born out of a passion for understanding and innovating within the blockchain space. It's more than just a code repository; it's a learning journey into the intricacies of blockchain technology. Developed from the ground up, GoGoChain is designed to illustrate the foundational elements of a blockchain, such as security protocols, transaction processes, and network synchronization, all through the lens of Go, a language known for its simplicity and efficiency.

## Simplified Overview
For those new to blockchain technology or Go:

Blockchain, Simplified: At its heart, GoGoChain is a chain of data blocks, each securely linked to the next, creating a tamper-resistant timeline of data.
Go, The Language of Choice: Go, or Golang, is known for its simplicity and efficiency, making it an ideal choice for building a streamlined and performant blockchain.
What You Can Do: With GoGoChain, you can create and manage user accounts, initiate and track transactions, and understand how blockchain maintains data integrity and security. All through simple command-line interactions.

## Features
- Monolithic Architecture: The blockchain is designed with a monolithic architecture, providing a solid foundation for seamless integration of future features and improvements.
- Command-Line Interface (CLI): The blockchain comes equipped with a powerful CLI that allows users to interact with and manage the blockchain effortlessly.
- Secure Account Creation: Using ECDSA, the system generates private keys, public keys, and addresses, ensuring the security of user accounts.
- Keystore Storage: Account information is stored securely in a keystore file, safeguarding sensitive data.
- Genesis Initialization: The blockchain starts with Block 0, initialized using a genesis.json file that defines initial allocations to addresses, laying the groundwork for the entire network.
- Transaction and Block Structures: Well-defined transaction and block structures are in place, with each block accommodating multiple transactions, maintaining the integrity of the ledger.
- Robust Hashing: The system employs robust hashing mechanisms for transactions and blocks, enhancing security and data integrity.
- Efficient Data Storage: Transactions and blocks are efficiently stored in a directory in key-pair format, utilizing LevelDB to ensure data consistency.
- Block State Database: A dedicated block state database meticulously tracks the balances of each address across all nodes, providing real-time insights into the network's financial health.
- Insightful CLI Commands: A variety of CLI commands are available for retrieving essential information, including current block details, address balances, transaction data, and block data.
- SendTransaction Capability: Users can seamlessly send their native tokens to other addresses, a process that involves transaction signing, transaction body creation, and memory pool insertion.
- Memory Pool Management: Unverified transactions find their place in a memory pool, awaiting validation before being included in the blockchain. Validators maintain their memory pools as well.
- Transaction Verification: The blockchain incorporates a robust verification process, ensuring that transactions adhere to network rules. It checks if the sent value is within the sender's balance and validates transaction signatures to prevent fraudulent activities.
- Networking with LibP2P: Leveraging LibP2P, the blockchain seamlessly connects nodes, facilitates transaction transmission to validators, and enhances overall network robustness.
- Peer Discovery: Automatic peer discovery and maintenance mechanisms establish a resilient network, allowing nodes to communicate and share vital information effectively.
- Bootnode Logic: Initial node discovery is made easy through bootnode logic, streamlining the onboarding process for new nodes.
- Flexible Configuration: Configuration options are available via the config.toml file, empowering users to customize network settings, including bootnode addresses.
- Node Synchronization: Nodes can synchronize with each other, ensuring that all participants operate on the same blockchain, maintaining network consensus.

## Getting Started

### Prerequisites
- Go installed on your system.

### Installation
1. Clone the GoGoChain repository to your local machine.
2. Build the project using `go build`.

## Usage
You can interact with GoGoChain using the following CLI commands:

### Initialize the Genesis Block
./gogochain -initGenesis 1

### Get Address Balance
./gogochain -balance "address"

### Get Current Block Number
./gogochain -blockNumber 1

### Send Transaction
./gogochain -sendTrx 1 -from "senderAddress" -to "recipientAddress" -value "amount" -auth "authorizationToken"

### Start a Node
./gogochain -startNode 1 -sourcePort 'portNumber' 

### Start a Validator Node
./gogochain -startNode 1 -sourcePort 'portNumber' -validator 1 -address "validatorAddress"
