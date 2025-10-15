package model

import (
	"context"
	"fmt"

	"github.com/tangthinker/secret-chat-server/core"
	"github.com/tangthinker/secret-chat-server/internal/model/schema"
	"gorm.io/gorm"
)

type MessagesModel struct {
	db *gorm.DB
}

func NewMessagesModel() *MessagesModel {
	d := core.GlobalHelper.DB.GetDB()
	if err := d.AutoMigrate(&schema.Messages{}); err != nil {
		panic(fmt.Sprintf("auto migrate err:%v", err))
	}
	return &MessagesModel{db: d}
}

func (m *MessagesModel) GetListByUid(ctx context.Context, uid string) ([]*schema.Messages, error) {
	var result []*schema.Messages
	if err := m.db.WithContext(ctx).Where("uid = ?", uid).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MessagesModel) Delete(ctx context.Context, msgIds []uint) error {
	return m.db.WithContext(ctx).Delete(&schema.Messages{}, "id in (?)", msgIds).Error
}

func (m *MessagesModel) Create(ctx context.Context, req *schema.Messages) error {
	return m.db.WithContext(ctx).Create(req).Error
}
