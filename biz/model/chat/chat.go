
package chat

type PrivateMessageReq struct {
	ToUserID string `json:"to_user_id"`
	Content  string `json:"content"`
}

type PrivateHistoryReq struct {
	ToUserID   string `json:"to_user_id"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
	MarkAsRead bool   `json:"mark_as_read"`
}

type UnreadReq struct {
	FromUserID string `json:"from_user_id"`
}

type GroupMessageReq struct {
	RoomID   string `json:"room_id"`
	Content  string `json:"content"`
}

type GroupHistoryReq struct {
	RoomID string `json:"room_id"`
	Page   int64  `json:"page"`
	Size   int64  `json:"size"`
}

type CheckOnlineReq struct {
	UserID string `json:"user_id"`
}

type GetUnreadCountReq struct {
	FromUserID string `json:"from_user_id"`
}

type MessageInfo struct {
	ID         uint   `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id,omitempty"`
	RoomID     string `json:"room_id,omitempty"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
	IsRead     bool   `json:"is_read,omitempty"`
}

type PrivateHistoryResp struct {
	Messages []MessageInfo `json:"messages"`
	Total    int64        `json:"total"`
	Page     int64        `json:"page"`
	Size     int64        `json:"size"`
}

type UnreadResp struct {
	Messages []MessageInfo `json:"messages"`
	Count    int           `json:"count"`
}

type GroupHistoryResp struct {
	Messages []MessageInfo `json:"messages"`
	Total    int64        `json:"total"`
	Page     int64        `json:"page"`
	Size     int64        `json:"size"`
}

type CheckOnlineResp struct {
	UserID string `json:"user_id"`
	Online bool   `json:"online"`
}

type UnreadCountResp struct {
	FromUserID string `json:"from_user_id"`
	Count      int64  `json:"count"`
}

type AckResp struct {
	Status string `json:"status"`
	MsgID  string `json:"msg_id"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
