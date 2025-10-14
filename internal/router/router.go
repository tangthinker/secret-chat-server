package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
)

func RegisterRouters(router fiber.Router) {
	rootGroup := router.Group("/api/v1/", middleware.TokenValid, middleware.UserHook)

	rootGroup.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})
}
