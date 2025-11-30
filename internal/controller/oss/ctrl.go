package oss

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/secret-chat-server/helper/response"
	"github.com/tangthinker/secret-chat-server/internal/service/oss"
)

type Ctrl struct {
	ossService *oss.Service
}

func New() *Ctrl {
	return &Ctrl{
		ossService: oss.NewService(),
	}
}

func (ctrl *Ctrl) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Upload: Bad Request")
	}
	openFile, err := file.Open()
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, "Upload: Internal Server Error")
	}
	defer openFile.Close()
	url, err := ctrl.ossService.Upload(ctx.Context(), openFile, file.Filename)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, "Upload: Internal Server Error")
	}
	return response.Success(ctx, url)
}
