package network

import (
	"context"
	"encoding/hex"
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
	rlp "github.com/malay44/chadChain/core/utils"
	"github.com/multiformats/go-multiaddr"
)

var (
	hostVar   host.Host
	CtxVar    context.Context
	PeerAddrs []string
)

func setupHost() (host.Host, error) {
	Hex := os.Getenv("PRIV_HEX")

	// Decode hex string to bytes
	privBytes, err := hex.DecodeString(Hex)
	if err != nil {
		panic(err)
	}

	// Parse bytes into a private key
	privKey, err := crypto.UnmarshalEd25519PrivateKey(privBytes)
	if err != nil {
		panic(err)
	}

	// priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)

	host, err := libp2p.New(libp2p.Identity(privKey), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/3000"))
	if err != nil {
		return nil, err
	}

	return host, nil
}

func connectToPeer(addr string) error {

	peerMA, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		return err
	}

	if err := hostVar.Connect(CtxVar, *peerAddrInfo); err != nil {
		return err
	}

	fmt.Println("Connected to", peerAddrInfo.String())
	return nil
}

func sendToAllPeers(msg message) {

	for _, p := range PeerAddrs {
		if p == hostVar.Addrs()[0].String()+"/p2p/"+hostVar.ID().String() {
			println("Skipping self address")
			continue
		}

		peerMA, err := multiaddr.NewMultiaddr(p)
		if err != nil {
			fmt.Println("Error creating multiaddr:", err)
			continue
		}

		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			fmt.Println("Error creating peer.AddrInfo:", err)
			continue
		}

		s, err := hostVar.NewStream(CtxVar, peerAddrInfo.ID, "/")
		if err != nil {
			fmt.Println("Error creating stream:", err)
			continue
		}

		send(s, msg)
	}
}

func send(s network.Stream, msg message) {
	encoder := json.NewEncoder(s)
	if err := encoder.Encode(msg); err != nil {
		fmt.Println("Error encoding message:", err)
	}
}

func streamHandler(s network.Stream) {
	decoder := json.NewDecoder(s)
	var decodedData interface{}
	var msg message
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	senderID := s.Conn().RemotePeer()

	switch msg.ID {
	case 0:

		err := rlp.DecodeData(msg.Data, decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		fmt.Println("Received PING:", decodedData)
		sendPongToPeer(senderID)

	case 1:
		err := rlp.DecodeData(msg.Data, decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		fmt.Println("Received PONG:", decodedData)

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

func sendPing() {
	data, err := rlp.EncodeData("PING", false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	sendToAllPeers(message{ID: 0, Code: 0, Want: 0, Data: data})
}

func sendPongToPeer(peerID peer.ID) {

	data, err := rlp.EncodeData("PONG", false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	// Create a new stream to the peer
	s, err := hostVar.NewStream(CtxVar, peerID, "/")
	if err != nil {
		fmt.Println("Error creating stream:", err)
		return
	}
	defer s.Close() // Close the stream when done

	// Send PONG message
	encoder := json.NewEncoder(s)
	if err := encoder.Encode(message{ID: 1, Code: 0, Want: 0, Data: data}); err != nil {
		fmt.Println("Error encoding message:", err)
	}
}

type message struct {
	ID   uint64 `json:"id"`
	Code int    `json:"code"`
	Want int    `json:"want"`
	Data []byte `json:"data"`
}

func Run() {
	var err error
	hostVar, err = setupHost()
	if err != nil {
		panic(err)
	}
	defer hostVar.Close()

	fmt.Println("Addresses:", hostVar.Addrs())
	// fmt.Println("ID:", hostVar.ID())
	fmt.Println("Concnated Addr:", hostVar.Addrs()[0].String()+"/p2p/"+hostVar.ID().String())
	// fmt.Println("Peer_ADDR:", os.Getenv("PEER_ADDR"))

	for _, addr := range PeerAddrs {
		if addr == hostVar.Addrs()[0].String()+"/p2p/"+hostVar.ID().String() {
			println("Skipping self address")
			continue
		}
		if err := connectToPeer(addr); err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
		fmt.Println("Connected to peer:", addr)
	}

	sendPing()

	hostVar.SetStreamHandler("/", streamHandler)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	<-sigCh
}
