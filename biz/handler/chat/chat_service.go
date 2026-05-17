
package chat

import (
	"net/http"

	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/service/chat"
	"video_website/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
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
		
		var user entity.User
		if err := mysql.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				hlog.Errorf("用户不存在: %v", err)
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}
			hlog.Errorf("查询用户失败: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		username = user.Username
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
		Hub:      chat.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		Rooms:    make(map[string]bool),
	}

	chat.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

