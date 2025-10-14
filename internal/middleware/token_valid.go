package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/user-center/pkg"
)

func sendForbidden(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusForbidden)
	return ctx.SendString("Forbidden: Invalid Token")
}

func TokenValid(ctx *fiber.Ctx) error {
	token := ""

	// 优先检查 Authorization 头
	headers := ctx.GetReqHeaders()
	if len(headers) > 0 {
		authorization := headers["Authorization"]
		if len(authorization) > 0 {
			token = authorization[0]
		}
	}

	// 如果 Authorization 头中没有 token，检查查询参数
	if token == "" {
		queryToken := ctx.Query(TokenKey)
		if queryToken != "" {
			token = queryToken
		}
	}

	// 如果 token 为空，返回错误
	if token == "" {
		log.Printf("Token validation failed: no token provided")
		return sendForbidden(ctx)
	}

	// 验证 token
	uid, err := pkg.TokenValid(token)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.SendString("Unauthorized: " + err.Error())
	}

	// 存储到上下文
	ctx.Locals(UIDKey, uid)
	ctx.Locals(TokenKey, token)

	return ctx.Next()
}
