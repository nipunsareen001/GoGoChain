package node

import (
	"GoGoChain/config"
	"GoGoChain/enum"
	"GoGoChain/filemanagement"
	"GoGoChain/p2p"
	"context"
	"crypto/rand"
	"log"

	"github.com/libp2p/go-libp2p/core/host"
)

var NodeAddr string

// LATER
// var Host host.Host

var isBootNode bool = false

// StartNode initializes and starts a P2P node.
func StartNode(sourcePort int, isValidator int, validatorAddress string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Loads configs as per the config.toml file, including bootnodes
	config.LoadConfig()

	_, host, err := setupHost(ctx, sourcePort)
	if err != nil {
		log.Fatalf("Failed to set up host: %v", err)
	}

	bootnodes := GetAllBootnodes()
	if len(bootnodes) < 1 {
		isBootNode = true
	}

	// Handle bootnode connection if not a bootnode itself
	if !isBootNode {
		connectToBootNode(ctx, host, bootnodes)
	}

	// connectToBootNode(ctx, host)

	// CONSENSUS RELATED
	// else {
	// 	consensus.IsInValidatorPool = true
	// }

	// Start periodic tasks for node management
	// startPeriodicTasks(ctx, host, NodeAddr, isValidator, validatorAddress)

	<-ctx.Done()
}

// setupHost initializes the host for the P2P network.
func setupHost(ctx context.Context, sourcePort int) (string, host.Host, error) {
	r := rand.Reader

	host, err := p2p.MakeHost(sourcePort, r)
	if err != nil {
		return "", nil, err
	}

	nodeAddr := p2p.GetHostAddress(host)

	NodeAddr = nodeAddr
	// LATER
	// p2p.Host = host

	file, err := filemanagement.OpenFileForAppending(string(enum.Peerlist))
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	filemanagement.AppendDataInFile(file, []byte(nodeAddr))

	p2p.StartListener(ctx, host, sourcePort, p2p.HandlePeerUpdateStream)

	return nodeAddr, host, nil
}

// startPeriodicTasks starts various periodic tasks necessary for node management.
func startPeriodicTasks(ctx context.Context, host host.Host, nodeAddr string, isValidator int, validatorAddress string) {
	// Start routine for updating peerlist
	go updatePeerList(ctx, host, nodeAddr)

	// // If the node is a validator, handle additional validator tasks
	// if isValidator == 1 {
	// 	handleValidatorNode(ctx, host, validatorAddress, nodeAddr)
	// }

}

// updatePeerList contains the logic for periodically updating the peer list.
func updatePeerList(ctx context.Context, host host.Host, nodeAddr string) {
	// Implementation of peer list updating
}
