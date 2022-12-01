package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	room := newRoom()
	go room.handleRoom()
	mux.HandleFunc("/ws", room.serveWs)

	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
