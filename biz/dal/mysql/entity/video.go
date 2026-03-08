package entity

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID           string         `gorm:"primary_key;column:id"`
	UserID       string         `gorm:"column:user_id"`
	Title        string         `gorm:"column:title"`
	Description  string         `gorm:"column:description"`
	VideoURL     string         `gorm:"column:video_url"`
	CoverURL     string         `gorm:"column:cover_url"`
	VisitCount   uint64         `gorm:"column:visit_count;default:0"`
	LikeCount    uint64         `gorm:"column:like_count;default:0"`
	CommentCount uint64         `gorm:"column:comment_count;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
