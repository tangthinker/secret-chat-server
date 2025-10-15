package chat

import "context"

type Service interface {
	Send(ctx context.Context, uid string, msg Message) error
	Handle(ctx context.Context, uid string, msg Message) error
}

type Message struct {
	MessageType MessageType
	Destination string
	Content     string
}

type MessageType int

const (
	MessageTypeSingle    MessageType = 1
	MessageTypeGroup     MessageType = 2
	MessageTypeBroadcast MessageType = 3
)
