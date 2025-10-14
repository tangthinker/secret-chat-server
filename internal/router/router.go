package router

import (
	"git.tangthinker.com/tangthinker/secret-chat-server/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRouters(router fiber.Router) {
	rootGroup := router.Group("/api/v1/server", middleware.TokenValid, middleware.UserHook)

	rootGroup.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})
}
