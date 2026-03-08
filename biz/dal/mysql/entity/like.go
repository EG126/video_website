package entity

import "time"

type Like struct {
	ID        string    `gorm:"primary_key;column:id;type:varchar(64)"`
	UserID    string    `gorm:"uniqueIndex:idx_user_video;column:user_id;type:varchar(64)"`
	VideoID   string    `gorm:"uniqueIndex:idx_user_video;column:video_id;type:varchar(64)"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
