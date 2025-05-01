package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

const dataFilePath string = "./comment.json"

type CommentSentRequest struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type CommentEntry struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	IP      string `json:"ip"`
	Time    uint64 `json:"time"`
}

type CommentListResponse struct {
	List []CommentEntry `json:"list"`
}

var mu sync.Mutex

// curl http://localhost:8086/comment
func commentListHandler(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	// defer 表示延迟到函数结束时执行（不管是正常结束还是错误结束）
	defer mu.Unlock()

	var list []CommentEntry
	if _, err := os.Stat(dataFilePath); err == nil {
		data, err := os.ReadFile(dataFilePath)
		if err == nil {
			json.Unmarshal(data, &list)
		}
	}

	start := 0
	if len(list) > 5 {
		start = len(list) - 5
	}
	latestList := list[start:]

	resp := CommentListResponse{List: latestList}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(resp)
}

/*
	curl -X POST http://localhost:8086/comment \
	  -H "Content-Type: application/json" \
	  -d '{"name":"g","message":"What?"}'
*/
func commentSentHandler(w http.ResponseWriter, r *http.Request) {
	var req CommentSentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" || req.Message == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newComment := CommentEntry{
		Name:    req.Name,
		Message: req.Message,
		IP:      getRealIP(r),
		Time:    uint64(time.Now().Unix()),
	}

	mu.Lock()
	defer mu.Unlock()

	var list []CommentEntry
	if _, err := os.Stat(dataFilePath); err == nil {
		data, err := os.ReadFile(dataFilePath)
		if err == nil {
			json.Unmarshal(data, &list)
		}
	}

	list = append(list, newComment)

	data, _ := json.MarshalIndent(list, "", "  ")
	os.WriteFile(dataFilePath, data, 0644)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"status":"ok"}`))
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		commentListHandler(w, r)
	} else if r.Method == http.MethodPost {
		commentSentHandler(w, r)
	} else {
		http.Error(w, "Only POST and GET method are allowed", http.StatusMethodNotAllowed)
		return
	}
}
