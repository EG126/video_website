package mysql

import (
	"context"
	"time"

	"video_website/biz/dal/mysql/entity"
)

func SavePrivateMessage(msg *entity.ChatMessage) error {
	return DB.Create(msg).Error
}

func SaveGroupMessage(msg *entity.ChatMessage) error {
	return DB.Create(msg).Error
}

func GetPrivateHistory(fromUserID, toUserID string, page, size int64) ([]entity.ChatMessage, int64) {
	var messages []entity.ChatMessage
	var total int64

	offset := (page - 1) * size

	DB.Model(&entity.ChatMessage{}).Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		fromUserID, toUserID, toUserID, fromUserID,
	).Count(&total)

	DB.Where(
		"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		fromUserID, toUserID, toUserID, fromUserID,
	).Order("created_at DESC").Limit(int(size)).Offset(int(offset)).Find(&messages)

	return messages, total
}

func GetUnreadMessages(userID, fromUserID string) []entity.ChatMessage {
	var messages []entity.ChatMessage

	DB.Where(
		"to_user_id = ? AND from_user_id = ? AND is_read = ?",
		userID, fromUserID, false,
	).Order("created_at DESC").Find(&messages)

	return messages
}

func MarkMessagesAsRead(userID, fromUserID string) error {
	return DB.Model(&entity.ChatMessage{}).Where(
		"to_user_id = ? AND from_user_id = ? AND is_read = ?",
		userID, fromUserID, false,
	).Update("is_read", true).Error
}

func GetGroupHistory(roomID string, page, size int64) ([]entity.ChatMessage, int64) {
	var messages []entity.ChatMessage
	var total int64

	offset := (page - 1) * size

	DB.Model(&entity.ChatMessage{}).Where("room_id = ?", roomID).Count(&total)

	DB.Where("room_id = ?", roomID).Order("created_at DESC").Limit(int(size)).Offset(int(offset)).Find(&messages)

	return messages, total
}

func CreateChatRoomIfNotExists(ctx context.Context, roomID, roomName string) error {
	var count int64
	DB.WithContext(ctx).Model(&entity.ChatRoom{}).Where("room_id = ?", roomID).Count(&count)
	if count == 0 {
		room := &entity.ChatRoom{
			RoomID:   roomID,
			RoomName: roomName,
		}
		return DB.WithContext(ctx).Create(room).Error
	}
	return nil
}

func AddRoomMemberIfNotExists(ctx context.Context, roomID, userID string) error {
	var count int64
	DB.WithContext(ctx).Model(&entity.RoomMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	if count == 0 {
		member := &entity.RoomMember{
			RoomID:   roomID,
			UserID:   userID,
			JoinedAt: time.Now(),
		}
		return DB.WithContext(ctx).Create(member).Error
	}
	return nil
}

func IsRoomMember(ctx context.Context, roomID, userID string) bool {
	var count int64
	DB.WithContext(ctx).Model(&entity.RoomMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	return count > 0
}
