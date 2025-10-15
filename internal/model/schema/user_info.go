package schema

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	UID        string `gorm:"type:varchar(128);not null;unique_index" json:"uid"`
	Nickname   string `gorm:"type:varchar(128);not null;default:''" json:"nickname"`
	Avatar     string `gorm:"type:varchar(128);not null;default:''" json:"avatar"`
	Phone      string `gorm:"type:varchar(128);not null;default:''" json:"phone"`
	Email      string `gorm:"type:varchar(128);not null;default:''" json:"email"`
	Background string `gorm:"type:varchar(128);not null;default:''" json:"background"`
}

func (ui *UserInfo) TableName() string {
	return "user_info"
}
