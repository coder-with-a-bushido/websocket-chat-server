package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Peer struct {
	room *Room
	conn *websocket.Conn
	// msg from room to be sent to the peer
	sendMsg chan []byte
}

func (peer *Peer) handleFromPeer() {
	defer func() {
		peer.room.unregister <- peer
		peer.conn.Close()
	}()

	for {
		_, msgVal, err := peer.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		newMsg := &Message{
			peer:  peer,
			value: msgVal,
		}
		peer.room.broadcast <- newMsg
	}
}

func (peer *Peer) handleToPeer() {
	defer func() {
		peer.conn.Close()
	}()

	for msg := range peer.sendMsg {
		err := peer.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
