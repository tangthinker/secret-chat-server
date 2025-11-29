package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
	"github.com/tangthinker/secret-chat-server/internal/service/connections"
	skep "github.com/tangthinker/skep-server-go/pkg"
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
	uid := conn.Locals(middleware.UIDKey).(string)
	token := conn.Locals(middleware.TokenKey).(string)

	mConn := connections.NewConn(conn)

	// 握手
	ecdsaPrivKey := core.GlobalHelper.Config.GetString("encrypt-conn.ecdsa-priv-key")
	handshakeTimeout := core.GlobalHelper.Config.GetDuration("encrypt-conn.handshake-timeout")
	skepProcessor := skep.NewSkep(mConn, handshakeTimeout, []string{uid, token}, ecdsaPrivKey)
	sharedKey, err := skepProcessor.Handshake()
	if err != nil {
		log.Errorf("handshake failed, uid: %s, token: %s, err: %v", uid, token, err)
		mConn.Close()
		return
	}

	mConn.SetEncryptKey(sharedKey)

	ctrl.connService.AddConnection(uid, mConn)
	for {
		message, err := mConn.ReadMessage()
		if err != nil {
			log.Errorf("read message failed, uid: %s, err: %v", uid, err)
			ctrl.connService.RemoveConnection(uid, mConn.GetConnId())
			break
		}
		err = ctrl.connService.Handle(uid, mConn.GetConnId(), message)
		if err != nil {
			log.Errorf("handle message failed, uid: %s, message: %s, err: %v", uid, message, err)
			ctrl.connService.RemoveConnection(uid, mConn.GetConnId())
			break
		}
	}
}
