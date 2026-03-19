package entity

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID           string         `gorm:"primary_key;column:id;type:varchar(64)"`
	UserID       string         `gorm:"column:user_id;type:varchar(64)"`
	Title        string         `gorm:"column:title;type:varchar(255)"`
	Description  string         `gorm:"column:description;type:text"`
	VideoURL     string         `gorm:"column:video_url;type:varchar(255)"`
	CoverURL     string         `gorm:"column:cover_url;type:varchar(255)"`
	VisitCount   uint64         `gorm:"column:visit_count;default:0"`
	LikeCount    uint64         `gorm:"column:like_count;default:0"`
	CommentCount uint64         `gorm:"column:comment_count;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Video) TableName() string {
	return "videos"
}
