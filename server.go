package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Request body structure
type MessageRequest struct {
	Message string `json:"message"`
}

// Response body structure
type MessageResponse struct {
	Reply string `json:"reply"`
}

// Handle /chat
func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reply := req.Message + "!"
	resp := MessageResponse{Reply: reply}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("Server started at http://localhost:8086")
	log.Fatal(http.ListenAndServe("127.0.0.1:8086", nil))
}
