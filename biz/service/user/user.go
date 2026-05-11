package user

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/pkg/bcrypt"
	"video_website/pkg/errno"
	"video_website/pkg/jwt"
	"video_website/pkg/utils"

	user "video_website/biz/model/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	redis2 "github.com/redis/go-redis/v9"
)

func Register(ctx context.Context, username, password string) error {
	hlog.CtxInfof(ctx, "开始处理注册请求, 用户名: %s", username)

	exist, err := mysql.IsUserExist(ctx, username)
	if err != nil {
		hlog.CtxErrorf(ctx, "检查用户名失败: %v", err)
		return errno.DBError
	}
	if exist {
		return errno.UsernameExists
	}

	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		hlog.CtxErrorf(ctx, "密码加密失败: %v", err)
		return errno.EncryptError
	}

	now := time.Now()
	newUser := &entity.User{
		ID:        utils.GenerateID(),
		Username:  username,
		Password:  hashedPwd,
		AvatarURL: "http://127.0.0.1:8888/static/avatars/default_avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := mysql.CreateUser(ctx, newUser); err != nil {
		hlog.CtxErrorf(ctx, "创建用户失败: %v", err)
		return errno.DBError
	}

	hlog.CtxInfof(ctx, "用户注册成功, 用户ID: %s", newUser.ID)
	return nil
}

type LoginResult struct {
	Users        []*user.UserDataResp
	AccessToken  string
	RefreshToken string
}

func Login(ctx context.Context, username, password string) (*LoginResult, error) {
	hlog.CtxInfof(ctx, "开始处理登录请求, 用户名: %s", username)

	u, err := mysql.GetUserByUsername(ctx, username)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询用户失败: %v", err)
		return nil, errno.DBError
	}
	if u == nil {
		return nil, errno.UserNotExist
	}

	if !bcrypt.CheckPasswordHash(password, u.Password) {
		return nil, errno.PasswordError
	}

	accessToken, err := jwt.GenerateAccessToken(u.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "生成access token失败: %v", err)
		return nil, errno.TokenError
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		hlog.CtxErrorf(ctx, "生成refresh token失败: %v", err)
		return nil, errno.TokenError
	}

	redisKey := "refresh:" + refreshToken
	if err := redis.RDB.Set(ctx, redisKey, u.ID, 7*24*time.Hour).Err(); err != nil {
		hlog.CtxErrorf(ctx, "存储refresh token失败: %v", err)
		return nil, errno.RedisError
	}
	hlog.CtxInfof(ctx, "用户登录成功, 用户ID: %s", u.ID)

	return &LoginResult{
		Users: []*user.UserDataResp{{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: "",
		}},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func Info(ctx context.Context, userID string) ([]*user.UserDataResp, error) {
	hlog.CtxInfof(ctx, "开始获取用户信息, 请求用户ID: %s", userID)

	u, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询用户失败: %v", err)
		return nil, errno.DBError
	}
	if u == nil {
		return nil, errno.UserNotExist
	}

	hlog.CtxInfof(ctx, "获取用户信息成功, 用户ID: %s", userID)
	return []*user.UserDataResp{{
		ID:        u.ID,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: "",
	}}, nil
}

func UploadAvatar(ctx context.Context, c *app.RequestContext, userID string) ([]*user.UserDataResp, error) {
	hlog.CtxInfof(ctx, "开始处理上传头像请求, 用户ID: %s", userID)

	file, err := c.FormFile("data")
	if err != nil {
		hlog.CtxErrorf(ctx, "获取文件表单失败: %v", err)
		return nil, errno.ParamError
	}

	src, err := file.Open()
	if err != nil {
		hlog.CtxErrorf(ctx, "打开文件失败: %v", err)
		return nil, errno.ParamError
	}
	defer src.Close()

	dataBytes, err := io.ReadAll(src)
	if err != nil {
		hlog.CtxErrorf(ctx, "读取文件失败: %v", err)
		return nil, errno.ParamError
	}
	hlog.CtxInfof(ctx, "获取头像文件成功, 文件名: %s, 大小: %d bytes", file.Filename, file.Size)

	fileName := userID + "_" + time.Now().Format("20060102150405") + ".jpg"
	filePath := filepath.Join("./static/avatars", fileName)
	if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
		hlog.CtxErrorf(ctx, "保存头像失败: %v", err)
		return nil, errno.FileError
	}
	avatarURL := "http://127.0.0.1:8888/static/avatars/" + fileName
	hlog.CtxInfof(ctx, "头像文件保存成功, 路径: %s", filePath)

	if err := mysql.UpdateUserAvatar(ctx, userID, avatarURL); err != nil {
		hlog.CtxErrorf(ctx, "更新头像失败: %v", err)
		return nil, errno.DBError
	}
	hlog.CtxInfof(ctx, "更新用户头像成功, 用户ID: %s", userID)

	u, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		hlog.CtxErrorf(ctx, "查询用户失败: %v", err)
		return nil, errno.DBError
	}

	hlog.CtxInfof(ctx, "上传头像处理完成, 用户ID: %s", userID)
	return []*user.UserDataResp{{
		ID:        u.ID,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: "",
	}}, nil
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

func Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	hlog.CtxInfof(ctx, "开始处理刷新token请求")

	redisKey := "refresh:" + refreshToken
	userID, err := redis.RDB.Get(ctx, redisKey).Result()
	if err == redis2.Nil {
		return nil, errno.TokenInvalid
	} else if err != nil {
		hlog.CtxErrorf(ctx, "Redis错误: %v", err)
		return nil, errno.RedisError
	}
	hlog.CtxInfof(ctx, "刷新token: 从Redis获取用户ID成功, 用户ID: %s", userID)

	newAccessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		hlog.CtxErrorf(ctx, "生成access token失败: %v", err)
		return nil, errno.TokenError
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		hlog.CtxErrorf(ctx, "生成refresh token失败: %v", err)
		return nil, errno.TokenError
	}

	redis.RDB.Del(ctx, redisKey)
	hlog.CtxInfof(ctx, "刷新token: 删除旧refresh token成功")

	newKey := "refresh:" + newRefreshToken
	if err := redis.RDB.Set(ctx, newKey, userID, 7*24*time.Hour).Err(); err != nil {
		hlog.CtxErrorf(ctx, "存储refresh token失败: %v", err)
		return nil, errno.RedisError
	}
	hlog.CtxInfof(ctx, "刷新token: 存储新refresh token成功")

	hlog.CtxInfof(ctx, "刷新token处理完成, 用户ID: %s", userID)
	return &RefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
