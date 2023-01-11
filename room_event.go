package main

type BaseRoomEvent struct {
	peer *Peer
	Data any
}

type MessageRoomEventData struct {
	PeerName string `json:"name"`
	Value    string `json:"value"`
}
