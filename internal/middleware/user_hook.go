package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/internal/model"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
)

func UserHook(ctx *fiber.Ctx) error {

	uid := ctx.Locals(UIDKey).(string)
	err := model.NewUserInfoModel().CreateIfNotExists(ctx.Context(), &schema.UserInfo{
		UID: uid,
	})
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString("User Hook: Internal Server Error")
	}

	userInfo, err := model.NewUserInfoModel().GetByUid(ctx.Context(), uid)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString("User Hook: Internal Server Error")
	}

	ctx.Locals(UserInfoKey, userInfo)

	return ctx.Next()

}
