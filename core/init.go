package core

import (
	"fmt"
	"log"

	"git.tangthinker.com/tangthinker/secret-chat-server/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type globalHelper struct {
	Config *Config
	DB     DB
}

var GlobalHelper *globalHelper

func Init(configPath string) {
	config := NewConfig(configPath)
	dbPath := config.GetString("database.path")
	GlobalHelper = &globalHelper{
		Config: config,
		DB:     NewSqliteDB(dbPath),
	}

	fmt.Println("------init cnf success------")
}

func StartServer(app *fiber.App) {

	logFilePath := GlobalHelper.Config.GetString("server.log-file-path")
	app.Use(middleware.LoggerInFile(logFilePath))

	serverPort := GlobalHelper.Config.GetString("server.port")
	log.Fatal(app.Listen(":" + serverPort))
}

func GetDBPath() string {
	return GlobalHelper.Config.GetString("database.path")
}
