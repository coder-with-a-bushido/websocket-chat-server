package main

type PeerEventKind int

const (
	PEERJOIN PeerEventKind = iota
	PEERMESSAGE
	PEERLEAVE
)

type BasePeerEvent struct {
	Kind PeerEventKind `json:"kind"`
	Data any           `json:"data"`
}

type JoinPeerEventData struct {
	Username string `json:"username"`
}

type MessagePeerEventData struct {
	Content string `json:"content"`
}

type LeavePeerEventData struct{}

var PeerEventDataFromKind = map[PeerEventKind]func() any{
	PEERJOIN:    func() any { return &JoinPeerEventData{} },
	PEERMESSAGE: func() any { return &MessagePeerEventData{} },
	PEERLEAVE:   func() any { return &LeavePeerEventData{} },
}
