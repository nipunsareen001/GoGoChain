package p2p

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

// MakeHost creates a new P2P host with a new key pair.
func MakeHost(port int, randomness io.Reader) (host.Host, error) {
	// Generates a new Secp256k1 key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Secp256k1, 2048, randomness)
	if err != nil {
		return nil, err
	}

	// Assemble the multiaddress for the host to listen on.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	// Create a new libp2p Host that listens on the given multiaddress with the generated private key.
	return libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
}

// GetHostAddress builds the full multiaddress for a given host.
func GetHostAddress(ha host.Host) string {
	// Build host multiaddress.
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().String()))

	// Now we can build a full multiaddress to reach this host by encapsulating both addresses:
	addr := ha.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr).String()

	return fullAddr
}

// StartListener sets up the host to listen for incoming streams and handle them using provided handlers.
func StartListener(ctx context.Context, ha host.Host, listenPort int, streamHandler network.StreamHandler) {
	// Set stream handlers for each of the protocols we want to listen for.
	ha.SetStreamHandler("/peerlist/1.0.0", streamHandler)

	// Attempt to find the actual local port if it's not specified.
	var port string
	for _, la := range ha.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		log.Println("Was not able to find the actual local port")
		return
	}

	// Log that the host is listening for connections.
	log.Printf("Host is listening on port: %s", port)

	// Additional setup or routines can go here, such as setting up more stream handlers or other initialization tasks.

	// Continue to listen and serve connections until context is done.
	<-ctx.Done()
	// Perform any necessary cleanup here.
	log.Println("Listener is shutting down...")
}
