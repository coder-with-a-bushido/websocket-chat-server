package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	allPeers    map[*Peer]bool
	peerActions *PeerActions
}

func createNewRoom() *Room {
	return &Room{
		allPeers:    make(map[*Peer]bool),
		peerActions: createPeerActionsForRoom(),
	}
}

func (room *Room) handleRoom() {
	for {
		select {
		case peer := <-room.peerActions.joinRoom:
			log.Println("Peer joined ", peer)
			room.allPeers[peer] = true
		case peer := <-room.peerActions.leaveRoom:
			log.Println("Peer left ", peer)
			if _, ok := room.allPeers[peer]; ok {
				close(peer.msgToPeer)
				delete(room.allPeers, peer)
			}
		case roomEvent := <-room.peerActions.broadcastRoomEvent:
			msg, err := json.Marshal(roomEvent.Data)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("broadcasting msg from peer", roomEvent.peer, ":", msg)

			for peer := range room.allPeers {
				if peer == roomEvent.peer {
					continue
				}
				select {
				case peer.msgToPeer <- msg:
				// this case is when there's no receiver for chan `peer.msgToPeer`
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

	log.Println("New peer connecting...")
	createNewPeer(conn, room)
}
