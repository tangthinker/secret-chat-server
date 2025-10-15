package schema

import "gorm.io/gorm"

type UserRelations struct {
	gorm.Model
	UID       string `gorm:"unique"`
	FriendIds string `gorm:"type:text"`
}
