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

	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8001"
	}
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
