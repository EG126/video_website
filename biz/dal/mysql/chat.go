
package mysql

import (
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
