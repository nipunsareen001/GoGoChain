package node

import (
	"GoGoChain/config"
	"GoGoChain/p2p"
	"context"
	"log"

	"github.com/libp2p/go-libp2p/core/host"
)

// GetBootnode returns the bootnodes from the application configuration.
func GetBootnode() string {
	return config.GetConfig().Node.P2P.Bootnodes[0]
}

// GetAllBootnodes returns the bootnodes from the application configuration.
func GetAllBootnodes() []string {
	config := config.GetConfig()
	return config.Node.P2P.Bootnodes
}

// connectToBootNode connects the current node to the bootnode if it's not a bootnode itself.
func connectToBootNode(ctx context.Context, h host.Host, bootnodes []string) {

	for _, target := range bootnodes {
		rw, err := p2p.ConnectionForPeerUpdates(ctx, h, target)
		if err != nil {
			log.Printf("Failed to connect to this bootnode: %s, error: %v\n Now Trying next one.", target, err)
			continue // Try the next bootnode
		}

		// Connection successful, create goroutines to read and write data.
		p2p.WritePeerUpdatesSender(rw, NodeAddr)

		log.Printf("Successfully connected to bootnode: %s\n", target)
		return // Exit the function once a successful connection is made
	}

	log.Println("Failed to connect to any bootnode.")
}
