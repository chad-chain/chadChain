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

func send(s network.Stream, msg message) {
	encoder := json.NewEncoder(s)
	if err := encoder.Encode(msg); err != nil {
		fmt.Println("Error encoding message:", err)
	}
}

func streamHandler(s network.Stream) {
	decoder := json.NewDecoder(s)
	var msg message
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	fmt.Println("Received message:\n",
		"\nID:", msg.ID,
		"\nCode:", msg.Code,
		"\nWant:", msg.Want,
		"\nData:", msg.Data,
	)

	switch msg.ID {
	case 0:
		send(s, msg)
		fmt.Println("Sent Hello message to", s.Conn().RemoteMultiaddr().String())

	case 1:
		fmt.Println("Received transaction. Response: List of encoded transactions (can be 1 or more).")

	case 2:
		fmt.Println("Received block. Response: Response: Encoded version of a single block (which was just mined)")

	case 3:
		fmt.Println("Request: List of block numbers (upto 10 max) Response (expected): Encoded version of a list of asked blocks")

	case 4:
		fmt.Println("Request (to which this response should be made): List of block numbers (upto 10 max) Response: Encoded version of a list of asked blocks")

	default:
		fmt.Println("ERR", msg)
	}
}

func sendInitialHelloMessage(ctx context.Context, host host.Host, peerAddrInfo peer.AddrInfo, peerMA multiaddr.Multiaddr) error {
	s, err := host.NewStream(ctx, peerAddrInfo.ID, "/")
	if err != nil {
		return err
	}
	defer s.Close()

	send(s, message{ID: 0, Code: 0, Want: 0, Data: "Hello"})
	fmt.Println("Sent Hello message to", peerMA.String())
	return nil
}

type message struct {
	ID   uint64      `json:"id"`
	Code int         `json:"code"`
	Want int         `json:"want"`
	Data interface{} `json:"data"`
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

	host.SetStreamHandler("/", streamHandler)

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

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	<-sigCh
}
