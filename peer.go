package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Peer struct {
	conn *websocket.Conn
	// msg from room to be sent to the peer
	msgToPeer chan []byte
	// peer actions
	peerActions *PeerActions
	info        *PeerInfo
}

func createNewPeer(conn *websocket.Conn, room *Room) {
	peer := &Peer{
		conn:        conn,
		msgToPeer:   make(chan []byte, 256),
		peerActions: room.peerActions,
		info:        nil,
	}

	go peer.handleFromPeer()
	go peer.handleToPeer()
}

// Handle Incoming Messages from Peer
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

		log.Println("Message from Peer", peer, ":", string(msgVal))

		var raw json.RawMessage
		baseEvent := BasePeerEvent{
			Data: &raw,
		}
		err = json.Unmarshal(msgVal, &baseEvent)
		if err != nil {
			log.Fatal(err)
		}

		actualEventData := PeerEventDataFromKind[baseEvent.Kind]()
		err = json.Unmarshal(raw, actualEventData)
		if err != nil {
			log.Fatal(err)
		}

		switch baseEvent.Kind {
		case PEERJOIN:
			peer.onJoin(actualEventData.(*JoinPeerEventData))
		case PEERMESSAGE:
			peer.onMessage(actualEventData.(*MessagePeerEventData))
		case PEERLEAVE:
			peer.onLeave(actualEventData.(*LeavePeerEventData))
		default:
			log.Println(baseEvent)
		}
	}
}

// Handle Sending Messages to Peer
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

// Event Handlers
//  |
//  V

func (peer *Peer) onJoin(event *JoinPeerEventData) {
	info := &PeerInfo{
		Name: event.Username,
	}
	peer.info = info

	peer.peerActions.joinRoom <- peer
}

func (peer *Peer) onMessage(event *MessagePeerEventData) {
	roomEvent := &BaseRoomEvent{
		peer: peer,
		Data: &MessageRoomEventData{
			PeerName: peer.info.Name,
			Value:    event.Content,
		},
	}

	peer.peerActions.broadcastRoomEvent <- roomEvent
}

func (peer *Peer) onLeave(event *LeavePeerEventData) {
	peer.peerActions.leaveRoom <- peer
	peer.conn.Close()
}
