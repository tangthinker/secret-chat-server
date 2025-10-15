package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
)

func StartServer(app *fiber.App) {

	logFilepath := core.GlobalHelper.Config.GetString("server.log-file-path")
	app.Use(middleware.LoggerInFile(logFilepath))

	serverPort := core.GlobalHelper.Config.GetString("server.port")
	log.Fatal(app.Listen(":" + serverPort))
}
