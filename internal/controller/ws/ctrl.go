package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/tangthinker/secret-chat-server/internal/service/connections"
)

type Ctrl struct {
	connService *connections.WebSocketConnections
}

func New() *Ctrl {
	return &Ctrl{
		connService: connections.NewWebSocketConnections(),
	}
}

func (ctrl *Ctrl) HandleConn(conn *websocket.Conn) {
	uid := conn.Locals("uid").(string)
	ctrl.connService.AddConnection(uid, conn)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			ctrl.connService.RemoveConnection(uid)
			break
		}
		err = ctrl.connService.HandleSource(uid, connections.SourceMessage{
			Type:    connections.SourceType(messageType),
			Content: message,
		})
		if err != nil {
			ctrl.connService.RemoveConnection(uid)
			break
		}
	}
}
