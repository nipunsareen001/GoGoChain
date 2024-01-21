package p2p

import (
	"GoGoChain/enum"
	"GoGoChain/filemanagement"
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

// HandlePeerUpdateStream sets up read and write goroutines for an incoming stream.
func HandlePeerUpdateStream(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadPeerUpdatesReciver(rw)
}

// ReadPeerUpdatesReciver continuously reads from the stream and updates the peer list.
func ReadPeerUpdatesReciver(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read from stream: %v", err)
			return // Exit the loop and end the goroutine
		}

		if str == "" {
			return
		}
		if str != "\n" {

			// go node.SendPeersToSender()
			//saves Sender's nodeAddrs in Reciver's peerlist, if not already
			str = strings.TrimSpace(str)
			isAlreadyInFile, err := filemanagement.IsAlreadyInFile(string(enum.Peerlist), str)
			if err != nil {
				log.Print(err)
				return
			}
			if !isAlreadyInFile {
				file, err := filemanagement.OpenFileForAppending(string(enum.Peerlist))
				if err != nil {
					log.Print(err)
					return
				}
				defer file.Close()
				filemanagement.AppendDataInFile(file, []byte(str))
			}

		}

	}
}

// WritePeerUpdatesSender continuously writes updates to the stream.
func WritePeerUpdatesSender(rw *bufio.ReadWriter, peerAddr string) {
	if peerAddr == "" {
		return
	}

	if _, err := rw.WriteString(fmt.Sprintf("%s\n", peerAddr)); err != nil {
		log.Printf("Failed to write to stream: %v", err)
		return // Handle the error appropriately (maybe try again or end the goroutine)
	}

	if err := rw.Flush(); err != nil {
		log.Printf("Failed to flush writer: %v", err)
	}
}

// ConnectionForPeerUpdates establishes a connection with the given peer and returns a buffered read-writer.
func ConnectionForPeerUpdates(ctx context.Context, host host.Host, destination string) (*bufio.ReadWriter, error) {
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to parse multiaddr: %v", err) // Return the error to the caller
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get addr info: %v", err) // Return the error to the caller
	}

	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	stream, err := host.NewStream(context.Background(), info.ID, "/peerlist/1.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create new stream: %v", err) // Return the error to the caller
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	return rw, nil
}
