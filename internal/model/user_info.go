package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
	"gorm.io/gorm"
)

type UserInfoModel struct {
	db *gorm.DB
}

func NewUserInfoModel() *UserInfoModel {
	d := core.GlobalHelper.DB.GetDB()
	if err := d.AutoMigrate(&schema.UserInfo{}); err != nil {
		panic(fmt.Sprintf("auto migrate error, %s", err))
	}
	return &UserInfoModel{
		db: d,
	}
}

func (ui *UserInfoModel) GetByUid(ctx context.Context, uid string) (*schema.UserInfo, error) {
	var user schema.UserInfo
	if err := ui.db.WithContext(ctx).Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ui *UserInfoModel) Save(ctx context.Context, req *schema.UserInfo) error {
	tbName := req.TableName()

	return ui.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user schema.UserInfo
		if err := tx.Table(tbName).Where("uid = ?", req.UID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(req).Error
			}
			return err
		}
		return tx.Table(tbName).Where("id = ?", user.ID).Updates(&schema.UserInfo{
			Nickname:   req.Nickname,
			Avatar:     req.Avatar,
			Phone:      req.Phone,
			Email:      req.Email,
			Background: req.Background,
		}).Error
	})
}

func (ui *UserInfoModel) CreateIfNotExists(ctx context.Context, req *schema.UserInfo) error {
	tbName := req.TableName()
	return ui.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user schema.UserInfo
		if err := tx.Table(tbName).Where("uid = ?", req.UID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(req).Error
			}
			return err
		}
		return nil
	})
}
