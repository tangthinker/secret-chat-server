package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/tangthinker/secret-chat-server/internal/service/chat"
)

type WebSocketConnections struct {
	connections map[string]*websocket.Conn
	mutex       sync.RWMutex

	chatService chat.Service
}

func NewWebSocketConnections() *WebSocketConnections {
	return &WebSocketConnections{
		connections: make(map[string]*websocket.Conn),
		mutex:       sync.RWMutex{},
	}
}

func (ws *WebSocketConnections) AddConnection(uid string, conn *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.connections[uid] = conn
}

func (ws *WebSocketConnections) RemoveConnection(uid string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	delete(ws.connections, uid)
}

func (ws *WebSocketConnections) SendMessage(ctx context.Context, uid string, message SourceMessage) error {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	conn, ok := ws.connections[uid]
	if !ok {
		return errors.New("connection not found")
	}
	return conn.WriteMessage(int(message.Type), message.Content)
}

func (ws *WebSocketConnections) HandleMessage(ctx context.Context, uid string, message SourceMessage) error {
	if message.Type == SourceTypeString {
		var msg chat.Message
		if err := json.Unmarshal(message.Content, &msg); err != nil {
			return fmt.Errorf("unmarshal msg err:%w", err)
		}
		return ws.chatService.Handle(ctx, uid, msg)
	}
	return nil
}
