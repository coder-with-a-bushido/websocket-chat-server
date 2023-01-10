package main

type BaseRoomEvent struct {
	peer *Peer
	Data any
}

type MessageRoomEventData struct {
	PeerName string
	Value    []byte
}
