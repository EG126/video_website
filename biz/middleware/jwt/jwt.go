package jwt

import (
	"context"
	"video_website/pkg/errno"
	"video_website/pkg/jwt"
	"video_website/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserInfo struct {
	UserID string
}

var UserInfoKey = "user_info"

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader("Access-Token"))
		if token == "" {
			response.SendResponse(c, errno.Unauthorized, nil)
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.SendResponse(c, err, nil)
			c.Abort()
			return
		}

		userInfo := &UserInfo{
			UserID: claims.UserID,
		}

		c.Set(UserInfoKey, userInfo)
		c.Next(ctx)
	}
}

func GetUserID(_ context.Context, c *app.RequestContext) string {
	value, exists := c.Get(UserInfoKey)
	if !exists {
		return ""
	}
	if userInfo, ok := value.(*UserInfo); ok {
		return userInfo.UserID
	}
	return ""
}
