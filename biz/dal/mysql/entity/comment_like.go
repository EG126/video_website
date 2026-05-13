package entity

import (
	"time"

	"gorm.io/gorm"
)

type CommentLike struct {
	ID        string         `gorm:"primary_key;column:id;type:varchar(64)"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);index"`
	CommentID string         `gorm:"column:comment_id;type:varchar(64);index"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}