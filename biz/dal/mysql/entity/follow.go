package entity

import "time"

type Follow struct {
	ID          string    `gorm:"primary_key;column:id;type:varchar(64)"`
	FollowerID  string    `gorm:"column:follower_id;type:varchar(64);uniqueIndex:idx_follow"`
	FollowingID string    `gorm:"column:following_id;type:varchar(64);uniqueIndex:idx_follow"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
