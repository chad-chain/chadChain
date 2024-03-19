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

	w.WriteHeader(http.StatusCreated)
}

func Rpc() {
	http.HandleFunc("/sendTx", sendTx)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
