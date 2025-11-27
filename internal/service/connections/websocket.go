package connections

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tangthinker/secret-chat-server/internal/model"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
)

type WebSocketConnections struct {
	connections map[string]*Conn
	mutex       sync.RWMutex

	messagesModel *model.MessagesModel
}

func NewWebSocketConnections() *WebSocketConnections {
	return &WebSocketConnections{
		connections: make(map[string]*Conn),
		mutex:       sync.RWMutex{},

		messagesModel: model.NewMessagesModel(),
	}
}

func (ws *WebSocketConnections) AddConnection(uid string, conn *Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if _, ok := ws.connections[uid]; ok {
		ws.connections[uid].Close()
		delete(ws.connections, uid)
	}
	ws.connections[uid] = conn

	ctx, cal := context.WithTimeout(context.Background(), 3*time.Second)
	defer cal()
	msgs, err := ws.messagesModel.GetListByUid(ctx, uid)
	if err != nil {
		return
	}
	msgIds := make([]uint, 0)
	for _, msg := range msgs {
		err = conn.SendMessage(msg.Content)
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
	conn, ok := ws.connections[uid]
	if !ok {
		return
	}
	conn.Close()
	delete(ws.connections, uid)
}

func (ws *WebSocketConnections) Send2User(uid string, message string) error {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	conn, ok := ws.connections[uid]
	if !ok {
		return errors.New("connection not found")
	}

	return conn.SendMessage(message)
}

func (ws *WebSocketConnections) Handle(uid string, message string) error {
	if message == "PING" {
		return ws.sendPONG(uid)
	}

	msg, err := ToMessage(message)
	if err != nil {
		return fmt.Errorf("unmarshal msg err:%w", err)
	}
	msg.From = uid
	msg.Timestamp = time.Now()

	if msg.MessageType == MessageTypeSingle {
		// 发送单聊消息
		err := ws.Send2User(msg.Destination, msg.String())
		// 发送失败，则保存消息到数据库
		if err != nil {
			log.Infof("send message to user failed, uid: %s, message: %s, err: %v", uid, message, err)
			err := ws.messagesModel.Create(context.Background(), &schema.Messages{
				Uid:     msg.Destination,
				Content: msg.String(),
			})
			if err != nil {
				log.Errorf("create messages error: %s", err)
				return err
			}
		}
	}
	return nil
}

func (ws *WebSocketConnections) sendPONG(uid string) error {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	conn, ok := ws.connections[uid]
	if !ok {
		return errors.New("connection not found")
	}
	return conn.SendMessage("PONG")
}
