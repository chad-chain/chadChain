package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	t "github.com/malay44/chadChain/core/types"
	r "github.com/malay44/chadChain/core/utils"
	"github.com/malay44/chadChain/core/validator"
	"github.com/multiformats/go-multiaddr"
)

var (
	hostVar       host.Host
	CtxVar        context.Context
	PeerAddrs     []string
	VoteThreshold = 2
)

type vote struct {
	blockNumber uint64
	yesVotes    int
	noVotes     int
}

var blockVotes map[uint64]*vote // Map to track votes for each block

func setupHost() (host.Host, error) {

	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0000"))
	if err != nil {
		return nil, err
	}

	return host, nil
}

func GetHostAddr() []string {
	addrs := []string{hostVar.Addrs()[0].String() + "/p2p/" + hostVar.ID().String(),
		hostVar.Addrs()[1].String() + "/p2p/" + hostVar.ID().String()}
	// Check if the first address is private, if not swap them
	if strings.Contains(addrs[0], "127.0.0.1") {
		return addrs
	} else {
		return []string{addrs[1], addrs[0]} // Swap the addresses
	}
}

func ConnectToPeer(addr string) error {

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
	SendPingToPeer(addr)
	return nil
}

func checkForSelf(addr string) bool {
	if addr == GetHostAddr()[0] || addr == GetHostAddr()[1] {
		return true
	}
	return false
}

func CreateStreamToPeer(peerAddr string) (network.Stream, error) {
	// Parse peer address string into a multiaddress
	peerMA, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return nil, fmt.Errorf("error creating multiaddress: %v", err)
	}

	// Extract peer ID from multiaddress
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		return nil, fmt.Errorf("error creating peer.AddrInfo: %v", err)
	}

	// Create a new stream to the peer
	s, err := hostVar.NewStream(CtxVar, peerAddrInfo.ID, "/")
	if err != nil {
		return nil, fmt.Errorf("error creating stream: %v", err)
	}

	return s, nil
}

func sendToAllPeers(msg message) {

	for _, p := range PeerAddrs {
		if checkForSelf(p) {
			println("Skipping self address")
			continue
		}

		s, err := CreateStreamToPeer(p)
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
	var decodedData []byte
	var msg message
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	senderID := s.Conn().RemotePeer()

	switch msg.ID {
	case 0:
		err := r.DecodeData(msg.Data, &decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		fmt.Println("Received PING: ", string(decodedData))
		if string(decodedData) == "PING" {
			sendPongToPeer(senderID)
		}

	case 1:
		err := r.DecodeData(msg.Data, &decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		fmt.Println("Received PONG:", string(decodedData))

	case 2:
		err := r.DecodeData(msg.Data, &decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		Address(string(decodedData))

	case 3:
		err := r.DecodeData(msg.Data, &decodedData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		ReceiveAddress(string(decodedData))

	case 4:
		fmt.Println("Request (to which this response should be made): List of block numbers (upto 10 max) Response: Encoded version of a list of asked blocks")

	case 5:
		block := t.Block{}
		err := r.DecodeData(decodedData, &block)
		if err != nil {
			fmt.Println("Error decoding data:", err)
			return
		}
		receiveBlock(block)

	case 10:
		handleVoteMessage(msg.Data)

	default:
		fmt.Println("ERR", msg)
	}
}

func handleVoteMessage(data []byte) {
	var voteData []uint64 // Change to slice

	err := r.DecodeData(data, &voteData)
	if err != nil {
		fmt.Println("Error decoding vote data:", err)
		return
	}

	// Ensure the vote data has at least two elements
	if len(voteData) < 2 {
		fmt.Println("Invalid vote data format")
		return
	}

	blockNumber := voteData[0]
	voteValue := uint8(voteData[1])

	fmt.Printf("Received vote for block %d: %d\n", blockNumber, voteValue)

	// Check if a vote has already been received for the block
	if _, ok := blockVotes[blockNumber]; !ok {
		// Create a new vote object if no vote has been received
		blockVotes[blockNumber] = &vote{blockNumber: blockNumber}
	}

	// Increment the vote count based on the vote value
	if voteValue == 1 {
		blockVotes[blockNumber].yesVotes++
	} else {
		blockVotes[blockNumber].noVotes++
	}

	println("block 1 votes", blockVotes[1].yesVotes, blockVotes[1].noVotes)

	// Check if the block has received enough votes
	if blockVotes[blockNumber].yesVotes >= VoteThreshold {
		commitBlock(blockNumber)
	} else if blockVotes[blockNumber].noVotes >= VoteThreshold {
		discardBlock(blockNumber)
	}
}

// Function to commit the block to the database
func commitBlock(blockNumber uint64) {
	// Implement logic to commit the block to the database
	fmt.Printf("Block %d committed to the database\n", blockNumber)
}

// Function to discard the block
func discardBlock(blockNumber uint64) {
	// Implement logic to discard the block
	fmt.Printf("Block %d discarded\n", blockNumber)
}

// Function to send a vote message for a block
func SendVote(blockNumber uint64, vote uint8) {

	data, err := r.EncodeData([2]uint64{blockNumber, uint64(vote)}, false)
	if err != nil {
		fmt.Println("Error encoding vote data:", err)
		return
	}
	sendToAllPeers(message{ID: 10, Code: 0, Want: 0, Data: data})
	println("Sent Vote for block number = ", blockNumber, " Vote = ", vote)
}

func SendBlock(block t.Block) {
	data, err := r.EncodeData(block, true)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	sendToAllPeers(message{ID: 5, Code: 0, Want: 0, Data: data})
}

func receiveBlock(block t.Block) {
	fmt.Println("Received Block number = ", block.Header.Number)
	if validator.ValidateBlock(&block) {
		SendVote(block.Header.Number, 1)
	} else {
		SendVote(block.Header.Number, 0)
	}
}

func Address(receivedAddress string) {
	// PeerAddrs = append(PeerAddrs, receivedAddress)
	fmt.Println("Address Received: ", receivedAddress)
	data, err := r.EncodeData(GetHostAddr()[1], false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	sendToAllPeers(message{ID: 3, Code: 0, Want: 0, Data: data})
}

func SendAddress(addr string) {
	data, err := r.EncodeData(addr, false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	sendToAllPeers(message{ID: 2, Code: 0, Want: 3, Data: data})
}

func ReceiveAddress(addr string) {
	fmt.Println("address Received: ", addr)
}

func SendPing() {
	data, err := r.EncodeData("PING", false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	sendToAllPeers(message{ID: 0, Code: 0, Want: 0, Data: data})
}

func SendPingToPeer(peerAddr string) {
	// Create a stream to the peer
	s, err := CreateStreamToPeer(peerAddr)
	if err != nil {
		fmt.Println("Error creating stream:", err)
		return
	}
	defer s.Close() // Close the stream when done

	data, err := r.EncodeData("PING", false)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	// Send PING message
	encoder := json.NewEncoder(s)
	if err := encoder.Encode(message{ID: 0, Code: 0, Want: 0, Data: data}); err != nil {
		fmt.Println("Error encoding message:", err)
	}
}

func sendPongToPeer(peerID peer.ID) {

	data, err := r.EncodeData("PONG", false)
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
	blockVotes = make(map[uint64]*vote)
	hostVar, err = setupHost()
	hostVar.SetStreamHandler("/", streamHandler)
	if err != nil {
		panic(err)
	}
	defer hostVar.Close()

	PeerAddrs = append(PeerAddrs, GetHostAddr()[1])

	go func() {
		Rpc()
	}()
	fmt.Println("Addresses:", hostVar.Addrs())
	// fmt.Println("ID:", hostVar.ID())
	fmt.Println("Concnated Addr:", hostVar.Addrs()[0].String()+"/p2p/"+hostVar.ID().String())
	// fmt.Println("Peer_ADDR:", os.Getenv("PEER_ADDR"))

	GetAllAddrsFromRoot()

	for _, addr := range PeerAddrs {
		if checkForSelf(addr) {
			println("Skipping self address")
			continue
		}
		if err := ConnectToPeer(addr); err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
		fmt.Println("Connected to peer:", addr)
	}
	// SendVote(1, 1)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	<-sigCh
}
