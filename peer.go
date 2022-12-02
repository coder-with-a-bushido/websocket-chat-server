package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Peer struct {
	conn *websocket.Conn
	// msg from room to be sent to the peer
	msgToPeer chan []byte
	// peer actions
	peerActions *PeerActions
}

func (peer *Peer) handleFromPeer() {
	defer func() {
		peer.peerActions.leaveRoom <- peer
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
		peer.peerActions.broadcastMsg <- newMsg
	}
}

func (peer *Peer) handleToPeer() {
	defer func() {
		peer.conn.Close()
	}()

	for msg := range peer.msgToPeer {
		err := peer.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
