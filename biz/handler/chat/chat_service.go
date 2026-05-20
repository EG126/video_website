package chat

import (
	"context"
	"net/http"

	"video_website/biz/service/chat"
	userService "video_website/biz/service/user"
	"video_website/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func HTTPWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	var userID, username string
	if token != "" {
		claims, err := jwt.ParseToken(token)
		if err != nil {
			hlog.Errorf("Token 解析失败: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID = claims.UserID

		users, err := userService.Info(context.Background(), userID)
		if err != nil {
			hlog.Errorf("查询用户失败: %v", err)
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		if len(users) == 0 {
			hlog.Errorf("用户不存在: %s", userID)
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		username = users[0].Username
	} else {
		userID = r.URL.Query().Get("user_id")
		username = r.URL.Query().Get("username")
		if userID == "" || username == "" {
			http.Error(w, "Missing user info", http.StatusBadRequest)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hlog.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &chat.Client{
		Hub:      chat.DefaultHub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		Rooms:    make(map[string]bool),
	}

	chat.DefaultHub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
