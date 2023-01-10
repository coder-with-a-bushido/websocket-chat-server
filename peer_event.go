package main

type PeerEventKind int

const (
	PEERJOIN PeerEventKind = iota
	PEERMESSAGE
	PEERLEAVE
)

type BasePeerEvent struct {
	Kind PeerEventKind
	Data any
}

type JoinPeerEventData struct {
	Username string
}

type MessagePeerEventData struct {
	Content []byte
}

type LeavePeerEventData struct{}

var PeerEventDataFromKind = map[PeerEventKind]func() any{
	PEERJOIN:    func() any { return &JoinPeerEventData{} },
	PEERMESSAGE: func() any { return &MessagePeerEventData{} },
	PEERLEAVE:   func() any { return &LeavePeerEventData{} },
}
