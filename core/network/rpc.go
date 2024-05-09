package network

import (
	// "encoding/json"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	PORT = "3001" // Port number for the RPC server
)

type PeerResponse struct {
	PeerAddrs []string `json:"peerAddrs"`
}

func GetAllAddrsFromRoot() {
	// Encode self string into JSON
	requestBody, err := json.Marshal(hostVar.Addrs()[0].String() + "/p2p/" + hostVar.ID().String())
	if err != nil {
		fmt.Println("Error encoding self string:", err)
		return
	}
	resp, err := http.Post("http://localhost:3000/getP2pAdr", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error getting response:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.StatusCode)
		return
	}

	var peerResp PeerResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	if err := json.Unmarshal(body, &peerResp); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println("Addresses:------------------------")
	for _, addr := range peerResp.PeerAddrs {
		fmt.Println(addr)
		ConnectToPeer(addr)
	}
	println("------------------")
}

func getP2pAdr(w http.ResponseWriter, r *http.Request) {
	var addr string
	if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	PeerAddrs = append(PeerAddrs, addr)
	PeerAddrs = removeDuplicates(PeerAddrs)

	fmt.Println("Addresses:------------------------")
	for _, addr := range PeerAddrs {
		fmt.Println(addr)
	}
	println("------------------")

	// Encode PeerAddrs to JSON
	peerResp := PeerResponse{PeerAddrs: PeerAddrs}
	peerAddrsJSON, err := json.Marshal(peerResp)
	if err != nil {
		http.Error(w, "Error encoding PeerAddrs to JSON", http.StatusInternalServerError)
		return
	}
	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(peerAddrsJSON)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func removeDuplicates(strs []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, str := range strs {
		if str == "" {
			continue
		}
		if !encountered[str] {
			encountered[str] = true
			result = append(result, str)
		}
	}

	return result
}

func sendTx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signed := r.URL.Query().Get("signed")
	fmt.Println("Value of signed:", signed)

	w.WriteHeader(http.StatusOK)
}

// /blockNumber: Returns the recent most block number
// /block?number={number}: Given the block number, returns the contents of a block else nil (Response type: json struct of block type described above)
// /block?hash={hash}: Given the block hash, returns the contents of a block else nil (Response type: json struct of block type described above)
// /tx?hash={hash}: Given the transaction hash, returns the contents of a transaction (without the v, r, s values) else nil (Response type: json struct of block type described above)
// /getNonce?address={address}: Given the address, returns the current nonce of that account
// /getBalance?address={address}: Given the address, returns the current balance/amount of that account

func blockNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func block(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("number") == "" && r.URL.Query().Get("hash") == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("number") != "" {
		number := r.URL.Query().Get("number")
		fmt.Println("Value of number:", number)
	}
	if r.URL.Query().Get("hash") != "" {
		hash := r.URL.Query().Get("hash")
		fmt.Println("Value of hash:", hash)
	}
	w.WriteHeader(http.StatusOK)
}

func tx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("hash") == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	hash := r.URL.Query().Get("hash")
	fmt.Println("Value of hash:", hash)
	w.WriteHeader(http.StatusOK)
}

func getNonce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("address") == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	address := r.URL.Query().Get("address")
	fmt.Println("Value of address:", address)
	w.WriteHeader(http.StatusOK)
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("address") == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	address := r.URL.Query().Get("address")
	fmt.Println("Value of address:", address)
	w.WriteHeader(http.StatusOK)
}

func Rpc() {

	http.HandleFunc("/getP2pAdr", getP2pAdr)
	http.HandleFunc("/sendTx", sendTx)
	http.HandleFunc("/blockNumber", blockNumber)
	http.HandleFunc("/block", block)
	http.HandleFunc("/tx", tx)
	http.HandleFunc("/getNonce", getNonce)
	http.HandleFunc("/getBalance", getBalance)

	fmt.Println("RPC server listening on port", PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
	}
}
