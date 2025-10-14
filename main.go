package main

import (
	"flag"

	"git.tangthinker.com/tangthinker/secret-chat-server/core"
	"git.tangthinker.com/tangthinker/secret-chat-server/internal/router"
	"github.com/gofiber/fiber/v2"
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
