package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	allPeers    map[*Peer]bool
	peerActions *PeerActions
}

func newRoom() *Room {
	return &Room{
		allPeers: make(map[*Peer]bool),
		peerActions: &PeerActions{
			joinRoom:     make(chan *Peer),
			leaveRoom:    make(chan *Peer),
			broadcastMsg: make(chan *Message),
		},
	}
}

func (room *Room) handleRoom() {
	for {
		select {
		case peer := <-room.peerActions.joinRoom:
			room.allPeers[peer] = true
		case peer := <-room.peerActions.leaveRoom:
			if _, ok := room.allPeers[peer]; ok {
				close(peer.msgToPeer)
				delete(room.allPeers, peer)
			}
		case msg := <-room.peerActions.broadcastMsg:
			for peer := range room.allPeers {
				if peer == msg.peer {
					continue
				}
				select {
				case peer.msgToPeer <- msg.value:
				// this case is when there's no receiver for chan `peer.sendMsg`
				default:
					close(peer.msgToPeer)
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
		conn:        conn,
		msgToPeer:   make(chan []byte, 256),
		peerActions: room.peerActions,
	}
	room.peerActions.joinRoom <- peer

	go peer.handleFromPeer()
	go peer.handleToPeer()
}
