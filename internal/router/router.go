package router

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/internal/controller/user_info"
	"github.com/tangthinker/secret-chat-server/internal/controller/ws"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
)

func RegisterRouters(router fiber.Router) {
	rootGroup := router.Group("/api/v1/", middleware.TokenValid, middleware.UserHook)

	rootGroup.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	userInfoCtrl := user_info.New()
	rootGroup.Post("/user/info/get", userInfoCtrl.GetUserInfo)
	rootGroup.Post("/user/info/update", userInfoCtrl.UpdateUserInfo)
	rootGroup.Post("/user/info/exists", userInfoCtrl.Exists)

	websocketCtrl := ws.New()
	rootGroup.Get("/websocket/conn", websocket.New(websocketCtrl.HandleConn))
}
