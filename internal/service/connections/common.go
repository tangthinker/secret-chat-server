package connections

import "context"

type Connections interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, uid string, message SourceMessage) error
	// HandleMessage 处理接受消息
	HandleMessage(ctx context.Context, uid string, message SourceMessage) error
}

type SourceMessage struct {
	Type    SourceType
	Content []byte
}

type SourceType int

const (
	SourceTypeString SourceType = 1
	SourceTypeBinary SourceType = 2
)
