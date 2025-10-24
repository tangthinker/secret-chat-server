package user_info

import (
	"context"

	"github.com/tangthinker/secret-chat-server/internal/model"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
)

type Service struct {
	userInfoMode *model.UserInfoModel
}

func NewService() *Service {
	return &Service{
		userInfoMode: model.NewUserInfoModel(),
	}
}

func (s *Service) GetUserInfo(ctx context.Context, userID string) (*schema.UserInfo, error) {
	return s.userInfoMode.GetByUid(ctx, userID)
}

func (s *Service) UpdateUserInfo(ctx context.Context, req *schema.UserInfo) error {
	return s.userInfoMode.Save(ctx, req)
}
