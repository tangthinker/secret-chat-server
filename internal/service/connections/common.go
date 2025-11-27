package connections

import (
	"encoding/json"
	"time"
)

type Connections interface {
	// SendSource 发送消息
	SendSource(uid string, message string) error
	// HandleSource 处理接受消息
	HandleSource(uid string, message string) error
}

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

func ToMessage(message string) (*Message, error) {
	var msg Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

type MessageType int

const (
	MessageTypeSingle    MessageType = 1
	MessageTypeGroup     MessageType = 2
	MessageTypeBroadcast MessageType = 3
)
