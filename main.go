package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/router"
	userPkg "github.com/tangthinker/user-center/pkg"
)

var configPath = flag.String("config", "./config/local/config.toml", "Path to config file")

func main() {
	flag.Parse()

	core.Init(*configPath)

	app := fiber.New()

	router.RegisterRouters(app)

	userPkg.RegisterUserCenter(app, core.GetDBPath())

	core.StartServer(app)
}
