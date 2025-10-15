package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tangthinker/secret-chat-server/internal/model"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
)

type WebSocketConnections struct {
	connections map[string]*websocket.Conn
	mutex       sync.RWMutex

	messagesModel *model.MessagesModel
}

func NewWebSocketConnections() *WebSocketConnections {
	return &WebSocketConnections{
		connections: make(map[string]*websocket.Conn),
		mutex:       sync.RWMutex{},

		messagesModel: model.NewMessagesModel(),
	}
}

func (ws *WebSocketConnections) AddConnection(uid string, conn *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.connections[uid] = conn

	ctx, cal := context.WithTimeout(context.Background(), 3*time.Second)
	defer cal()
	msgs, err := ws.messagesModel.GetListByUid(ctx, uid)
	if err != nil {
		return
	}
	msgIds := make([]uint, 0)
	for _, msg := range msgs {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		if err != nil {
			log.Infof("write message to websocket fail, uid: %s, msg: %s", uid, msg.Content)
			continue
		}
		msgIds = append(msgIds, msg.ID)
	}
	if len(msgIds) > 0 {
		if err := ws.messagesModel.Delete(context.Background(), msgIds); err != nil {
			log.Infof("delete synced messages error: %s", err)
		}
	}
}

func (ws *WebSocketConnections) RemoveConnection(uid string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	delete(ws.connections, uid)
}

func (ws *WebSocketConnections) SendSource(uid string, message SourceMessage) error {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	conn, ok := ws.connections[uid]
	if !ok {
		return errors.New("connection not found")
	}
	return conn.WriteMessage(int(message.Type), message.Content)
}

func (ws *WebSocketConnections) HandleSource(uid string, message SourceMessage) error {
	if message.Type == SourceTypeString {
		var msg Message
		if err := json.Unmarshal(message.Content, &msg); err != nil {
			return fmt.Errorf("unmarshal msg err:%w", err)
		}
		msg.From = uid
		msg.Timestamp = time.Now()
		return ws.HandleMessage(msg)
	}
	return nil
}

func (ws *WebSocketConnections) HandleMessage(message Message) error {
	if message.MessageType == MessageTypeSingle {
		err := ws.SendSource(message.Destination, SourceMessage{
			Type:    SourceTypeString,
			Content: []byte(message.String()),
		})
		if err != nil {
			log.Infof("send message: %s; err:%s", message.String(), err)
			err := ws.messagesModel.Create(context.Background(), &schema.Messages{
				Uid:     message.Destination,
				Content: message.String(),
			})
			if err != nil {
				log.Errorf("create messages error: %s", err)
				return err
			}
		}
	}
	return nil
}
