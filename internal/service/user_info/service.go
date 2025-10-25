package user_info

import (
	"context"

	"github.com/tangthinker/secret-chat-server/internal/model"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
	"github.com/tangthinker/secret-chat-server/internal/proto"
)

type Service struct {
	userInfoMode *model.UserInfoModel
}

func NewService() *Service {
	return &Service{
		userInfoMode: model.NewUserInfoModel(),
	}
}

func (s *Service) GetUserInfo(ctx context.Context, req *proto.UserInfoGetReq) (*proto.UserInfoGetResp, error) {
	userInfo, err := s.userInfoMode.GetByUid(ctx, req.UID)
	if err != nil {
		return nil, err
	}
	return &proto.UserInfoGetResp{
		UserInfo: userInfo,
	}, nil
}

func (s *Service) UpdateUserInfo(ctx context.Context, req *schema.UserInfo) error {
	return s.userInfoMode.Save(ctx, req)
}

func (s *Service) Exists(ctx context.Context, req *proto.UserExistsReq) (*proto.UserExistsResp, error) {
	exists, err := s.userInfoMode.Exists(ctx, req.UID)
	if err != nil {
		return nil, err
	}
	return &proto.UserExistsResp{
		Exists: exists,
	}, nil
}
