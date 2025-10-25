package proto

import "github.com/tangthinker/secret-chat-server/internal/model/schema"

type UserExistsReq struct {
	UID string `json:"uid"`
}

type UserExistsResp struct {
	Exists bool `json:"exists"`
}

type UserInfoGetReq struct {
	UID string `json:"uid"`
}

type UserInfoGetResp struct {
	UserInfo *schema.UserInfo `json:"user_info"`
}
