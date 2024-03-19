package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const PORT = "8080"

type ActionRequest struct {
	ID   uint64      `json:"id"`
	Code int         `json:"code"`
	Want int         `json:"want"`
	Data interface{} `json:"data"`
}

type ActionResponse interface{}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var request ActionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response ActionResponse

	switch request.ID {
	case 0:
		response = hello(request)
	case 1:
		response = newTransaction(request)
	case 2:
		response = newBlock(request)
	case 3:
		response = getBlock(request)
	case 4:
		response = getblockdetails(request)
	default:
		response = "Invalid action ID."
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func hello(request ActionRequest) ActionResponse {
	return request
}

func newTransaction(request ActionRequest) ActionResponse {
	return "New transaction added."
}

func newBlock(request ActionRequest) ActionResponse {
	return "New block added."
}

func getBlock(request ActionRequest) ActionResponse {
	return "Block requested."
}

func getblockdetails(request ActionRequest) ActionResponse {
	return "Block requested."
}

// Http function sets up the HTTP server to handle requests
func Http() {
	http.HandleFunc("/", handleRequest)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
