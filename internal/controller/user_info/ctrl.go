package user_info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tangthinker/secret-chat-server/helper/response"
	"github.com/tangthinker/secret-chat-server/internal/middleware"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
	"github.com/tangthinker/secret-chat-server/internal/proto"
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
	req := &proto.UserInfoGetReq{}
	if err := ctx.BodyParser(req); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Get User Info: Bad Request")
	}
	userInfoResp, err := ctrl.userInfoService.GetUserInfo(ctx.Context(), req)
	if err != nil {
		log.Errorf("get user info error: %s", err)
		return response.Error(ctx, fiber.StatusInternalServerError, "Get User Info: Internal Server Error")
	}
	return response.Success(ctx, userInfoResp)
}

func (ctrl *Ctrl) UpdateUserInfo(ctx *fiber.Ctx) error {
	userInfo := &schema.UserInfo{}
	if err := ctx.BodyParser(userInfo); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Update User Info: Bad Request")
	}
	uid := ctx.Locals(middleware.UIDKey).(string)
	userInfo.UID = uid
	err := ctrl.userInfoService.UpdateUserInfo(ctx.Context(), userInfo)
	if err != nil {
		log.Errorf("update user info error: %s", err)
		return response.Error(ctx, fiber.StatusInternalServerError, "Update User Info: Internal Server Error")
	}
	return response.Success(ctx, "Update User Info Success")
}

func (ctrl *Ctrl) Exists(ctx *fiber.Ctx) error {
	req := &proto.UserExistsReq{}
	if err := ctx.BodyParser(req); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Exists User Info: Bad Request")
	}
	exists, err := ctrl.userInfoService.Exists(ctx.Context(), req)
	if err != nil {
		log.Errorf("exists user info error: %s", err)
		return response.Error(ctx, fiber.StatusInternalServerError, "Exists User Info: Internal Server Error")
	}
	return response.Success(ctx, exists)
}
