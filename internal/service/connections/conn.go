package connections

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	encrypt "github.com/tangthinker/encrypt-conn-tools/pkg"
	skep "github.com/tangthinker/skep-server-go/pkg"
)

var _ skep.Conn = (*Conn)(nil)

type Conn struct {
	conn       *websocket.Conn
	encryptKey string
	connId     string
}

func NewConn(conn *websocket.Conn) *Conn {
	connId := uuid.New().String()
	return &Conn{
		conn:   conn,
		connId: connId,
	}
}

func (c *Conn) GetConnId() string {
	return c.connId
}

func (c *Conn) SetEncryptKey(encryptKey string) {
	c.encryptKey = encryptKey
}

func (c *Conn) ReadFunc() (string, error) {
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	if messageType != websocket.TextMessage {
		return "", fmt.Errorf("invalid message type: %d", messageType)
	}
	return string(message), nil
}

func (c *Conn) WriteFunc(message string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) ReadMessage() (string, error) {
	messageType, message, err := c.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	if messageType != websocket.TextMessage {
		return "", fmt.Errorf("invalid message type: %d", messageType)
	}
	if string(message) == "PING" {
		return string(message), nil
	}
	if c.encryptKey == "" {
		return "", fmt.Errorf("encrypt key is not set")
	}
	decryptedMessage := encrypt.Decrypt(string(message), c.encryptKey)

	return decryptedMessage, nil
}

func (c *Conn) SendMessage(data string) error {
	if data == "PONG" {
		return c.conn.WriteMessage(websocket.TextMessage, []byte(data))
	}

	if c.encryptKey == "" {
		return fmt.Errorf("encrypt key is not set")
	}

	encryptedMessage := encrypt.Encrypt(data, c.encryptKey)

	return c.conn.WriteMessage(websocket.TextMessage, []byte(encryptedMessage))

}
