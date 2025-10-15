package connections

import (
	"encoding/json"
	"time"
)

type Connections interface {
	// SendSource 发送消息
	SendSource(uid string, message SourceMessage) error
	// HandleSource 处理接受消息
	HandleSource(uid string, message SourceMessage) error
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

type Message struct {
	MessageType MessageType `json:"message_type"`
	From        string      `json:"from"`
	Destination string      `json:"destination"`
	Content     string      `json:"content"`
	Timestamp   time.Time   `json:"timestamp"`
}

func (m *Message) String() string {
	jsData, _ := json.Marshal(m)
	return string(jsData)
}

type MessageType int

const (
	MessageTypeSingle    MessageType = 1
	MessageTypeGroup     MessageType = 2
	MessageTypeBroadcast MessageType = 3
)
