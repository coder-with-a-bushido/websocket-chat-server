package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	room := createNewRoom()
	go room.handleRoom()
	mux.HandleFunc("/ws", room.serveWs)

	client := http.FileServer(http.Dir("./client"))
	mux.Handle("/", client)

	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8080"
	}
	log.Println("Starting server on port ", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
