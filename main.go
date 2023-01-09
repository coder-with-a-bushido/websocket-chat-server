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

	client := http.FileServer(http.Dir("./client"))
	mux.Handle("/",client)

	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8001"
	}
	log.Println("Starting server...")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
