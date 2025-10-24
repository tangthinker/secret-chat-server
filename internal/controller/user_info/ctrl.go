package user_info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tangthinker/secret-chat-server/helper/response"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
	"github.com/tangthinker/secret-chat-server/internal/service/user_info"
)

type Ctrl struct {
	userInfoService *user_info.Service
}

func New() *Ctrl {
	return &Ctrl{
		userInfoService: user_info.NewService(),
	}
}

func (ctrl *Ctrl) GetUserInfo(ctx *fiber.Ctx) error {
	userID := ctx.Locals(middleware.UIDKey).(string)
	userInfo, err := ctrl.userInfoService.GetUserInfo(ctx.Context(), userID)
	if err != nil {
		log.Errorf("get user info error: %s", err)
		return response.Error(ctx, fiber.StatusInternalServerError, "Get User Info: Internal Server Error")
	}
	return response.Success(ctx, userInfo)
}

func (ctrl *Ctrl) UpdateUserInfo(ctx *fiber.Ctx) error {
	userInfo := &schema.UserInfo{}
	if err := ctx.BodyParser(userInfo); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Update User Info: Bad Request")
	}

	err := ctrl.userInfoService.UpdateUserInfo(ctx.Context(), userInfo)
	if err != nil {
		log.Errorf("update user info error: %s", err)
		return response.Error(ctx, fiber.StatusInternalServerError, "Update User Info: Internal Server Error")
	}
	return response.Success(ctx, "Update User Info Success")
}
