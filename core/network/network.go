package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const PORT = "8080"

// ActionRequest represents the structure of the incoming requests
type ActionRequest struct {
	ID string `json:"id"`
}

// ActionResponse represents the structure of the response for different actions
type ActionResponse struct {
	Response string `json:"response"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var request ActionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response ActionResponse

	switch request.ID {
	case "0":
		response.Response = "Should receive a same response back."
	case "1":
		response.Response = "List of encoded transactions (can be 1 or more)."
	case "2":
		response.Response = "Encoded version of a single block (which was just mined)."
	case "3":
		response.Response = "Encoded version of a list of asked blocks."
	case "4":
		response.Response = "Encoded version of a list of asked blocks."
	default:
		response.Response = "Invalid action ID."
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
