package main

import (
	"encoding/json"
	"net/http"
)

type MessageRequest struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Reply string `json:"reply"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// 加上 CORS 头，允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 如果是预检请求，直接返回
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 只允许 POST
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
