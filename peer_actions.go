package main

type PeerActions struct {
	joinRoom           chan *Peer
	leaveRoom          chan *Peer
	broadcastRoomEvent chan *BaseRoomEvent
}

func createPeerActionsForRoom() *PeerActions {
	return &PeerActions{
		joinRoom:           make(chan *Peer),
		leaveRoom:          make(chan *Peer),
		broadcastRoomEvent: make(chan *BaseRoomEvent),
	}
}
