package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/comment", withCORS(commentHandler))

	fmt.Println("Server started at http://localhost:8086")
	log.Fatal(http.ListenAndServe("127.0.0.1:8086", nil))
}
