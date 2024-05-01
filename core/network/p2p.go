package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func setupHost() (host.Host, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		return nil, err
	}

	return host, nil
}

func connectToPeer(ctx context.Context, host host.Host, addr string) error {
	peerMA, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		return err
	}

	if err := host.Connect(ctx, *peerAddrInfo); err != nil {
		return err
	}

	fmt.Println("Connected to", peerAddrInfo.String())
	return nil
}

func handleIncomingStreams(ctx context.Context, host host.Host) {
	host.SetStreamHandler("/Hello", func(s network.Stream) {
		// Your existing stream handling code goes here
	})
}

func send(s network.Stream, msg string) {
	encoder := json.NewEncoder(s)
	if err := encoder.Encode(msg); err != nil {
		fmt.Println("Error encoding message:", err)
	}
}

func streamHandler(s network.Stream) {
	decoder := json.NewDecoder(s)
	var msg string
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	fmt.Println("Received message:", msg)
}

func sendInitialHelloMessage(ctx context.Context, host host.Host, peerAddrInfo peer.AddrInfo, peerMA multiaddr.Multiaddr) error {
	s, err := host.NewStream(ctx, peerAddrInfo.ID, "/Hello")
	if err != nil {
		return err
	}
	defer s.Close()

	send(s, "Hello")
	fmt.Println("Sent Hello message to", peerMA.String())
	return nil
}

func Run(ctx context.Context, peerAddrs []string) {
	host, err := setupHost()
	if err != nil {
		panic(err)
	}
	defer host.Close()

	fmt.Println("Addresses:", host.Addrs())
	fmt.Println("ID:", host.ID())
	// fmt.Println("Peer_ADDR:", os.Getenv("PEER_ADDR"))

	for _, addr := range peerAddrs {
		if err := connectToPeer(ctx, host, addr); err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
	}

	handleIncomingStreams(ctx, host)

	// Send initial Hello message to each peer
	for _, addr := range peerAddrs {
		peerMA, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			fmt.Println("Error parsing peer address:", err)
			continue
		}
		fmt.Println("peerMA:", peerMA)

		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			fmt.Println("Error creating peer address info:", err)
			continue
		}

		fmt.Println("peerAddrInfo:", peerAddrInfo)

		if err := sendInitialHelloMessage(ctx, host, *peerAddrInfo, peerMA); err != nil {
			fmt.Println("Error sending initial Hello message to peer:", err)
			continue
		}
	}

	host.SetStreamHandler("/Hello", streamHandler)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	<-sigCh
}
