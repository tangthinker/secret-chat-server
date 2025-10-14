package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func UserHook(ctx *fiber.Ctx) error {

	//uid := ctx.Locals(UIDKey).(string)
	//_, err := model.NewUserModel().CreateIfNotExists(ctx.Context(), uid)
	//if err != nil {
	//	ctx.Status(fiber.StatusInternalServerError)
	//	return ctx.SendString("User Hook: Internal Server Error")
	//}
	//
	//userInfo, err := model.NewUserModel().GetByUid(ctx.Context(), uid)
	//if err != nil {
	//	ctx.Status(fiber.StatusInternalServerError)
	//	return ctx.SendString("User Hook: Internal Server Error")
	//}
	//
	//ctx.Locals(UserInfoKey, userInfo)

	return ctx.Next()

}
