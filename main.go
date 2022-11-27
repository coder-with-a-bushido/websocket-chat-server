package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	room := newRoom()
	go room.handleRoom()
	mux.HandleFunc("/ws", room.serveWs)

	err := http.ListenAndServe(":"+"8001", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
