package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        string         `gorm:"primary_key;column:id;type:varchar(64)"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);index"`
	VideoID   string         `gorm:"column:video_id;type:varchar(64);index"`
	Content   string         `gorm:"column:content;type:text"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Comment) TableName() string {
	return "comments"
}
