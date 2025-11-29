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
	connections map[string][]*Conn
	mutex       sync.RWMutex

	messagesModel *model.MessagesModel
}

func NewWebSocketConnections() *WebSocketConnections {
	return &WebSocketConnections{
		connections: make(map[string][]*Conn),
		mutex:       sync.RWMutex{},

		messagesModel: model.NewMessagesModel(),
	}
}

func (ws *WebSocketConnections) AddConnection(uid string, conn *Conn) {
	ws.mutex.Lock()
	if _, ok := ws.connections[uid]; ok {
		ws.connections[uid] = append(ws.connections[uid], conn)
		ws.mutex.Unlock()
		return
	}
	ws.connections[uid] = []*Conn{conn}
	ws.mutex.Unlock()

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

func (ws *WebSocketConnections) RemoveConnection(uid string, connId string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	conns, ok := ws.connections[uid]
	if !ok {
		return
	}
	for i, conn := range conns {
		if conn.connId == connId {
			conn.Close()
			ws.connections[uid] = append(conns[:i], conns[i+1:]...)

			if len(ws.connections[uid]) == 0 {
				delete(ws.connections, uid)
			}
			return
		}
	}
}

func (ws *WebSocketConnections) Send2User(uid string, message string) error {
	ws.mutex.RLock()
	conns, ok := ws.connections[uid]
	if !ok {
		ws.mutex.RUnlock()
		return errors.New("connection not found")
	}

	targetConns := make([]*Conn, len(conns))
	copy(targetConns, conns)
	ws.mutex.RUnlock()

	successCount := 0
	for _, conn := range targetConns {
		err := conn.SendMessage(message)
		if err != nil {
			log.Infof("send message to user failed, uid: %s, message: %s, err: %v", uid, message, err)
			continue
		}
		successCount++
	}
	if successCount == 0 {
		return errors.New("send message to user failed")
	}
	return nil
}

func (ws *WebSocketConnections) Handle(uid string, connId string, message string) error {
	if message == "PING" {
		return ws.sendPONG(uid, connId)
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

func (ws *WebSocketConnections) sendPONG(uid string, connId string) error {
	ws.mutex.RLock()
	conns, ok := ws.connections[uid]
	if !ok {
		ws.mutex.RUnlock()
		return errors.New("connection not found")
	}

	targetConns := make([]*Conn, len(conns))
	copy(targetConns, conns)
	ws.mutex.RUnlock()

	for _, conn := range targetConns {
		if conn.connId == connId {
			return conn.SendMessage("PONG")
		}
	}
	return errors.New("connection not found")
}
