package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	allPeers   map[*Peer]bool
	register   chan *Peer
	unregister chan *Peer
	// msg from peer to be sent to room
	broadcast chan *Message
}

func newRoom() *Room {
	return &Room{
		allPeers:   make(map[*Peer]bool),
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
		broadcast:  make(chan *Message),
	}
}

func (room *Room) handleRoom() {
	for {
		select {
		case peer := <-room.register:
			room.allPeers[peer] = true
		case peer := <-room.unregister:
			if _, ok := room.allPeers[peer]; ok {
				close(peer.sendMsg)
				delete(room.allPeers, peer)
			}
		case msg := <-room.broadcast:
			for peer := range room.allPeers {
				if peer == msg.peer {
					continue
				}
				select {
				case peer.sendMsg <- msg.value:
				// this case is when there's no receiver for chan `peer.sendMsg`
				default:
					close(peer.sendMsg)
					delete(room.allPeers, peer)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{}

func (room *Room) serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	peer := &Peer{
		room:    room,
		conn:    conn,
		sendMsg: make(chan []byte, 256),
	}
	room.register <- peer

	go peer.handleFromPeer()
	go peer.handleToPeer()
}
