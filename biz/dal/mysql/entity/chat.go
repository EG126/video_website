package entity

import (
	"time"

	"gorm.io/gorm"
)

type ChatMessage struct {
	ID         uint           `gorm:"primaryKey;column:id"`
	FromUserID string         `gorm:"index;column:from_user_id;type:varchar(64)"`
	ToUserID   string         `gorm:"index;column:to_user_id;type:varchar(64)"`
	RoomID     string         `gorm:"index;column:room_id;type:varchar(64)"`
	Content    string         `gorm:"column:content;type:text"`
	MsgType    int            `gorm:"column:msg_type;comment:1:私聊 2:群聊"`
	IsRead     bool           `gorm:"column:is_read;default:false"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

type ChatRoom struct {
	ID        uint           `gorm:"primaryKey;column:id"`
	RoomID    string         `gorm:"uniqueIndex;column:room_id;type:varchar(64)"`
	RoomName  string         `gorm:"column:room_name;type:varchar(100)"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}

type RoomMember struct {
	ID        uint           `gorm:"primaryKey;column:id"`
	RoomID    string         `gorm:"index;column:room_id;type:varchar(64)"`
	UserID    string         `gorm:"index;column:user_id;type:varchar(64)"`
	JoinedAt  time.Time      `gorm:"column:joined_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (RoomMember) TableName() string {
	return "room_members"
}
