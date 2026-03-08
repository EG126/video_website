package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primary_key;column:id;type:varchar(64)"`
	Username  string         `gorm:"uniqueIndex;column:username;type:varchar(50)"`
	Password  string         `gorm:"column:password;type:varchar(100)"`
	AvatarURL string         `gorm:"column:avatar_url;type:varchar(255)"`
	CreatedAt time.Time      `gorm:"column:created_at;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;"`
}
