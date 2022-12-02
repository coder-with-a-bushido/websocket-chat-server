package main

type PeerActions struct {
	joinRoom     chan *Peer
	leaveRoom    chan *Peer
	broadcastMsg chan *Message
}
