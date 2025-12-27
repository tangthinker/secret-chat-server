package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
)

func StartServer(app *fiber.App) {

	app.Use(middleware.LoggerInConsole())

	serverPort := core.GlobalHelper.Config.GetString("server.port")
	log.Fatal(app.Listen(":" + serverPort))
}
