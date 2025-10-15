package schema

import "gorm.io/gorm"

type Messages struct {
	gorm.Model
	Uid     string `gorm:"type:varchar(255);not null"`
	Content string `gorm:"type:text"`
}

func (m *Messages) TableName() string {
	return "messages"
}
