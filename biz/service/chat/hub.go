package chat

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/biz/model/chat"
	"video_website/pkg/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Username string
	Rooms    map[string]bool
	mu       sync.RWMutex
}

type BroadcastMsg struct {
	RoomID  string
	Message []byte
	Exclude *Client
}

type PrivateMsg struct {
	Type    string                 `json:"type"`
	From    string                 `json:"from"`
	To      string                 `json:"to"`
	Payload map[string]interface{} `json:"payload"`
}

type WSResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMsg
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

var DefaultHub *Hub

func InitHub() {
	ctx, cancel := context.WithCancel(context.Background())
	DefaultHub = &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		Broadcast:  make(chan *BroadcastMsg, 512),
		ctx:        ctx,
		cancel:     cancel,
	}
	go DefaultHub.Run()
	go DefaultHub.SubscribePrivateMessages()
	go DefaultHub.SubscribeRoomMessages()
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			if err := redis.SetUserOnline(client.UserID); err != nil {
				hlog.Errorf("设置用户 %s 上线状态失败: %v", client.UserID, err)
			}
			hlog.Infof("用户 %s 上线", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, exists := h.Clients[client.UserID]; exists {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			if err := redis.SetUserOffline(client.UserID); err != nil {
				hlog.Errorf("设置用户 %s 离线状态失败: %v", client.UserID, err)
			}
			hlog.Infof("用户 %s 离线", client.UserID)

		case msg := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				if msg.Exclude != client {
					client.mu.RLock()
					if client.Rooms[msg.RoomID] {
						select {
						case client.Send <- msg.Message:
						default:
							close(client.Send)
							delete(h.Clients, client.UserID)
						}
					}
					client.mu.RUnlock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SubscribePrivateMessages() {
	pubsub := redis.SubscribePrivateMessages(h.ctx)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-h.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			var pm PrivateMsg
			if err := json.Unmarshal([]byte(msg.Payload), &pm); err != nil {
				hlog.Errorf("解析私聊消息失败: %v", err)
				continue
			}
			h.deliverPrivateMessage(pm.To, pm.Payload)
		}
	}
}

func (h *Hub) SubscribeRoomMessages() {
	pubsub := redis.SubscribeRoomMessages(h.ctx)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-h.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			var roomMsg struct {
				RoomID  string                 `json:"room_id"`
				Payload map[string]interface{} `json:"payload"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &roomMsg); err != nil {
				hlog.Errorf("解析群聊消息失败: %v", err)
				continue
			}
			h.Broadcast <- &BroadcastMsg{
				RoomID:  roomMsg.RoomID,
				Message: mustMarshal(roomMsg.Payload),
			}
		}
	}
}

func (h *Hub) deliverPrivateMessage(toUserID string, payload map[string]interface{}) {
	h.mu.RLock()
	client, exists := h.Clients[toUserID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	resp := WSResponse{
		Type:    "type1_push",
		Payload: payload,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
		hlog.Warnf("用户 %s 发送队列已满，消息可能丢失", toUserID)
	}
}

func (h *Hub) IsUserOnline(userID string) bool {
	return redis.IsUserOnline(userID)
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				hlog.Errorf("WebSocket error: %v", err)
			}
			break
		}

		var msgData map[string]interface{}
		if err := json.Unmarshal(message, &msgData); err != nil {
			c.SendError("invalid_json", "无效的JSON格式")
			continue
		}

		msgType, ok := msgData["type"].(string)
		if !ok {
			c.SendError("missing_type", "缺少消息类型")
			continue
		}

		switch msgType {
		case "type1":
			c.handleType1(msgData)
		case "type2":
			c.handleType2(msgData)
		case "type3":
			c.handleType3(msgData)
		case "type4":
			c.handleType4(msgData)
		case "type5":
			c.handleType5(msgData)
		case "check_online":
			c.handleCheckOnline(msgData)
		case "get_unread_count":
			c.handleGetUnreadCount(msgData)
		default:
			c.SendError("unknown_type", "未知消息类型")
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) SendError(code, message string) {
	resp := WSResponse{
		Type: "error",
		Payload: chat.ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
	data, _ := json.Marshal(resp)
	c.Send <- data
}

func (c *Client) SendResponse(msgType string, payload interface{}) {
	resp := WSResponse{
		Type:    msgType,
		Payload: payload,
	}
	data, _ := json.Marshal(resp)
	c.Send <- data
}

func (c *Client) handleType1(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	toUserID, ok := payload["to_user_id"].(string)
	content, contentOk := payload["content"].(string)

	if !ok || !contentOk || toUserID == "" || content == "" {
		c.SendError("missing_fields", "缺少必要字段")
		return
	}

	chatMsg := &entity.ChatMessage{
		FromUserID: c.UserID,
		ToUserID:   toUserID,
		Content:    content,
		MsgType:    1,
		CreatedAt:  time.Now(),
	}

	go func() {
		if err := mysql.SavePrivateMessage(chatMsg); err != nil {
			hlog.Errorf("保存私聊消息失败: %v", err)
		}
	}()
	go func() {
		redis.IncrementUnreadCount(toUserID, c.UserID)
	}()
	go c.publishPrivateMessageToUser(toUserID, chatMsg)

	c.SendResponse("type1_ack", chat.AckResp{
		Status: "success",
		MsgID:  utils.GenerateID(),
	})
}

func (c *Client) handleType2(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	toUserID, ok := payload["to_user_id"].(string)
	page := int64(getFloatValue(payload, "page", 1))
	size := int64(getFloatValue(payload, "size", 20))
	markAsRead := payload["mark_as_read"] == true

	if !ok || toUserID == "" {
		c.SendError("missing_fields", "缺少必要字段")
		return
	}

	messages, total := c.getPrivateHistory(c.UserID, toUserID, page, size)

	if markAsRead {
		go func() {
			mysql.MarkMessagesAsRead(c.UserID, toUserID)
		}()
		go func() {
			redis.ClearUnreadCount(c.UserID, toUserID)
		}()
	}

	c.SendResponse("type2_resp", chat.PrivateHistoryResp{
		Messages: messages,
		Total:    total,
		Page:     page,
		Size:     size,
	})
}

func (c *Client) handleType3(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	fromUserID, ok := payload["from_user_id"].(string)
	if !ok || fromUserID == "" {
		c.SendError("missing_fields", "缺少必要字段(需要from_user_id)")
		return
	}

	messages := c.getUnreadMessages(c.UserID, fromUserID)

	go func() {
		mysql.MarkMessagesAsRead(c.UserID, fromUserID)
	}()
	go func() {
		redis.ClearUnreadCount(c.UserID, fromUserID)
	}()

	c.SendResponse("type3_resp", chat.UnreadResp{
		Messages: messages,
		Count:    len(messages),
	})
}

func (c *Client) handleType4(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	roomID, ok := payload["room_id"].(string)
	content, contentOk := payload["content"].(string)

	if !ok || !contentOk || roomID == "" || content == "" {
		c.SendError("missing_fields", "缺少必要字段")
		return
	}

	c.mu.Lock()
	c.Rooms[roomID] = true
	c.mu.Unlock()

	go func() {
		if err := mysql.CreateChatRoomIfNotExists(context.Background(), roomID, "群聊房间"+roomID); err != nil {
			hlog.Errorf("创建房间失败: %v", err)
		}
		if err := mysql.AddRoomMemberIfNotExists(context.Background(), roomID, c.UserID); err != nil {
			hlog.Errorf("添加房间成员失败: %v", err)
		}
	}()

	chatMsg := &entity.ChatMessage{
		FromUserID: c.UserID,
		RoomID:     roomID,
		Content:    content,
		MsgType:    2,
		CreatedAt:  time.Now(),
	}

	go func() {
		if err := mysql.SaveGroupMessage(chatMsg); err != nil {
			hlog.Errorf("保存群聊消息失败: %v", err)
		}
	}()
	go c.publishGroupMessage(roomID, chatMsg)

	c.SendResponse("type4_ack", chat.AckResp{
		Status: "success",
		MsgID:  utils.GenerateID(),
	})
}

func (c *Client) handleType5(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	roomID, ok := payload["room_id"].(string)
	page := int64(getFloatValue(payload, "page", 1))
	size := int64(getFloatValue(payload, "size", 20))

	if !ok || roomID == "" {
		c.SendError("missing_fields", "缺少必要字段")
		return
	}

	if !mysql.IsRoomMember(context.Background(), roomID, c.UserID) {
		c.SendError("not_member", "你不是该群聊的成员")
		return
	}

	messages, total := c.getGroupHistory(roomID, page, size)

	c.SendResponse("type5_resp", chat.GroupHistoryResp{
		Messages: messages,
		Total:    total,
		Page:     page,
		Size:     size,
	})
}

func (c *Client) handleCheckOnline(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	userID, ok := payload["user_id"].(string)
	if !ok || userID == "" {
		c.SendError("missing_fields", "缺少必要字段(user_id)")
		return
	}

	isOnline := c.Hub.IsUserOnline(userID)

	c.SendResponse("check_online_resp", chat.CheckOnlineResp{
		UserID: userID,
		Online: isOnline,
	})
}

func (c *Client) handleGetUnreadCount(msgData map[string]interface{}) {
	payload, ok := msgData["payload"].(map[string]interface{})
	if !ok {
		c.SendError("invalid_payload", "无效的payload")
		return
	}

	fromUserID, ok := payload["from_user_id"].(string)
	if !ok || fromUserID == "" {
		c.SendError("missing_fields", "缺少必要字段(from_user_id)")
		return
	}

	count := redis.GetUnreadCount(c.UserID, fromUserID)

	c.SendResponse("unread_count_resp", chat.UnreadCountResp{
		FromUserID: fromUserID,
		Count:      count,
	})
}

func (c *Client) getPrivateHistory(fromUserID, toUserID string, page, size int64) ([]chat.MessageInfo, int64) {
	messages, total := mysql.GetPrivateHistory(fromUserID, toUserID, page, size)

	result := make([]chat.MessageInfo, 0, len(messages))
	for _, m := range messages {
		result = append(result, chat.MessageInfo{
			ID:         m.ID,
			FromUserID: m.FromUserID,
			ToUserID:   m.ToUserID,
			Content:    m.Content,
			CreatedAt:  m.CreatedAt.Unix(),
			IsRead:     m.IsRead,
		})
	}

	return result, total
}

func (c *Client) getUnreadMessages(userID, fromUserID string) []chat.MessageInfo {
	messages := mysql.GetUnreadMessages(userID, fromUserID)

	result := make([]chat.MessageInfo, 0, len(messages))
	for _, m := range messages {
		result = append(result, chat.MessageInfo{
			ID:         m.ID,
			FromUserID: m.FromUserID,
			Content:    m.Content,
			CreatedAt:  m.CreatedAt.Unix(),
		})
	}

	return result
}

func (c *Client) getGroupHistory(roomID string, page, size int64) ([]chat.MessageInfo, int64) {
	messages, total := mysql.GetGroupHistory(roomID, page, size)

	result := make([]chat.MessageInfo, 0, len(messages))
	for _, m := range messages {
		result = append(result, chat.MessageInfo{
			ID:         m.ID,
			FromUserID: m.FromUserID,
			RoomID:     m.RoomID,
			Content:    m.Content,
			CreatedAt:  m.CreatedAt.Unix(),
		})
	}

	return result, total
}

func (c *Client) publishPrivateMessageToUser(toUserID string, msg *entity.ChatMessage) {
	payload := map[string]interface{}{
		"id":           msg.ID,
		"from_user_id": msg.FromUserID,
		"content":      msg.Content,
		"created_at":   msg.CreatedAt.Unix(),
	}

	c.Hub.mu.RLock()
	_, exists := c.Hub.Clients[toUserID]
	c.Hub.mu.RUnlock()

	if exists {
		c.Hub.deliverPrivateMessage(toUserID, payload)
	} else {
		pm := PrivateMsg{
			Type:    "type1_push",
			From:    c.UserID,
			To:      toUserID,
			Payload: payload,
		}
		pmData, _ := json.Marshal(pm)
		go func() {
			if err := redis.PublishPrivateMessage(toUserID, pmData); err != nil {
				hlog.Errorf("发布私聊消息失败: %v", err)
			}
		}()
	}
}

func (c *Client) publishGroupMessage(roomID string, msg *entity.ChatMessage) {
	payload := map[string]interface{}{
		"id":           msg.ID,
		"from_user_id": msg.FromUserID,
		"room_id":      msg.RoomID,
		"content":      msg.Content,
		"created_at":   msg.CreatedAt.Unix(),
	}

	c.Hub.mu.RLock()
	hasLocalMembers := false
	for _, client := range c.Hub.Clients {
		client.mu.RLock()
		if client.Rooms[roomID] {
			hasLocalMembers = true
			client.mu.RUnlock()
			break
		}
		client.mu.RUnlock()
	}
	c.Hub.mu.RUnlock()

	if hasLocalMembers {
		c.Hub.Broadcast <- &BroadcastMsg{
			RoomID:  roomID,
			Message: mustMarshal(payload),
			Exclude: c,
		}
	}

	go func() {
		roomPayload := map[string]interface{}{
			"room_id": roomID,
			"payload": payload,
		}
		roomData, _ := json.Marshal(roomPayload)
		if err := redis.PublishRoomMessage(roomID, roomData); err != nil {
			hlog.Errorf("发布房间消息失败: %v", err)
		}
	}()
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func getFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if v, ok := data[key].(float64); ok {
		return v
	}
	return defaultValue
}
