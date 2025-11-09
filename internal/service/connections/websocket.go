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
	"github.com/tangthinker/secret-chat-server/pkg"
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
		// 历史消息内容已经是 JSON 字符串，需要加密后发送
		encryptedContent, err := ws.encryptBinaryContent([]byte(msg.Content))
		if err != nil {
			log.Errorf("encrypt history message failed, uid: %s, msg: %s, err: %v", uid, msg.Content, err)
			continue
		}
		err = conn.WriteMessage(websocket.BinaryMessage, encryptedContent)
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

	// 加密消息内容（发送加密后的二进制数据）
	encryptedContent, err := ws.encryptBinaryContent(message.Content)
	if err != nil {
		log.Errorf("encrypt message content failed: %v", err)
		return fmt.Errorf("encrypt message content failed: %w", err)
	}

	// 发送加密后的二进制数据
	return conn.WriteMessage(websocket.BinaryMessage, encryptedContent)
}

func (ws *WebSocketConnections) HandleSource(uid string, message SourceMessage) error {
	if message.Type == SourceTypeString && string(message.Content) == "PING" {
		return ws.sendPONG(uid)
	}
	// 客户端发送的是加密后的二进制数据，需要先解密
	decryptedContent, err := ws.decryptBinaryContent(message.Content)
	if err != nil {
		log.Errorf("decrypt message content failed: %v", err)
		return fmt.Errorf("decrypt message content failed: %w", err)
	}

	// 解密后解析 JSON
	var msg Message
	if err := json.Unmarshal(decryptedContent, &msg); err != nil {
		return fmt.Errorf("unmarshal msg err:%w", err)
	}
	msg.From = uid
	msg.Timestamp = time.Now()
	return ws.HandleMessage(msg)
}

func (ws *WebSocketConnections) HandleMessage(message Message) error {
	// 消息内容已经是解密后的 JSON 字符串，直接使用
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

// encryptBinaryContent 加密二进制消息内容
// 发送消息前需要先加密，然后发送加密后的二进制数据
func (ws *WebSocketConnections) encryptBinaryContent(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	// 加密二进制数据
	encryptedData, err := pkg.Encrypt(data)
	if err != nil {
		return nil, fmt.Errorf("encrypt failed: %w", err)
	}

	return encryptedData, nil
}

// decryptBinaryContent 解密二进制消息内容
// 客户端发送的是加密后的二进制数据，必须解密才能解析 JSON
// 如果解密失败，直接返回错误
func (ws *WebSocketConnections) decryptBinaryContent(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data is empty")
	}

	// 解密二进制数据
	decryptedData, err := pkg.Decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	return decryptedData, nil
}

func (ws *WebSocketConnections) sendPONG(uid string) error {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	conn, ok := ws.connections[uid]
	if !ok {
		return errors.New("connection not found")
	}
	return conn.WriteMessage(websocket.PongMessage, []byte("PONG"))
}
