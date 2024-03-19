package network

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
)

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
	http.HandleFunc("/sendTx", sendTx)
	http.HandleFunc("/blockNumber", blockNumber)
	http.HandleFunc("/block", block)
	http.HandleFunc("/tx", tx)
	http.HandleFunc("/getNonce", getNonce)
	http.HandleFunc("/getBalance", getBalance)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
