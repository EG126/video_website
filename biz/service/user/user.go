package user

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/mysql/entity"
	"video_website/biz/dal/redis"
	"video_website/config"
	"video_website/pkg/bcrypt"
	"video_website/pkg/constants"
	"video_website/pkg/errno"
	"video_website/pkg/jwt"
	"video_website/pkg/utils"

	user "video_website/biz/model/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgErrors "github.com/pkg/errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	redis2 "github.com/redis/go-redis/v9"
)

func Register(ctx context.Context, username, password string) error {
	exist, err := mysql.IsUserExist(ctx, username)
	if err != nil {
		return pkgErrors.Wrap(err, "service.Register: IsUserExist failed")
	}
	if exist {
		return errno.UsernameExists
	}

	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		return pkgErrors.Wrap(err, "service.Register: HashPassword failed")
	}

	now := time.Now()
	newUser := &entity.User{
		ID:        utils.GenerateID(),
		Username:  username,
		Password:  hashedPwd,
		AvatarURL: config.Conf.Static.BaseURL + config.Conf.Static.DefaultAvatar,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := mysql.CreateUser(ctx, newUser); err != nil {
		return pkgErrors.Wrap(err, "service.Register: CreateUser failed")
	}

	hlog.CtxInfof(ctx, "用户注册成功, userID=%s, username=%s", newUser.ID, username)
	return nil
}

type LoginResult struct {
	Users        []*user.UserDataResp
	AccessToken  string
	RefreshToken string
}

func Login(ctx context.Context, username, password, code string) (*LoginResult, error) {
	u, err := mysql.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Login: GetUserByUsername failed")
	}
	if u == nil {
		return nil, errno.UserNotExist
	}

	if !bcrypt.CheckPasswordHash(password, u.Password) {
		return nil, errno.PasswordError
	}

	if u.MFASecret != "" {
		if code == "" {
			return nil, errno.MFARequired
		}
		if !totp.Validate(code, u.MFASecret) {
			return nil, errno.MFAError
		}
	}

	accessToken, err := jwt.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Login: GenerateAccessToken failed")
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Login: GenerateRefreshToken failed")
	}

	redisKey := constants.RefreshTokenPrefix + refreshToken
	refreshExpire := time.Duration(config.Conf.Expire.RefreshToken) * time.Second
	if err := redis.RDB.Set(ctx, redisKey, u.ID, refreshExpire).Err(); err != nil {
		return nil, pkgErrors.Wrap(err, "service.Login: Redis Set failed")
	}

	hlog.CtxInfof(ctx, "用户登录成功, userID=%s, username=%s", u.ID, username)

	return &LoginResult{
		Users: []*user.UserDataResp{{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			CreatedAt: u.CreatedAt.Format(constants.DateTimeFormat),
			UpdatedAt: u.UpdatedAt.Format(constants.DateTimeFormat),
			DeletedAt: "",
		}},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func Info(ctx context.Context, userID string) ([]*user.UserDataResp, error) {
	u, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Info: GetUserByID failed")
	}
	if u == nil {
		return nil, errno.UserNotExist
	}

	return []*user.UserDataResp{{
		ID:        u.ID,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format(constants.DateTimeFormat),
		UpdatedAt: u.UpdatedAt.Format(constants.DateTimeFormat),
		DeletedAt: "",
	}}, nil
}

func UploadAvatar(ctx context.Context, c *app.RequestContext, userID string) ([]*user.UserDataResp, error) {
	file, err := c.FormFile("data")
	if err != nil {
		return nil, errno.ParamError
	}

	src, err := file.Open()
	if err != nil {
		return nil, errno.ParamError
	}
	defer src.Close()

	dataBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.UploadAvatar: io.ReadAll failed")
	}

	fileName := userID + "_" + time.Now().Format(constants.FileNameDateFormat) + constants.AvatarExt
	filePath := filepath.Join("./static/avatars", fileName)
	if err := os.WriteFile(filePath, dataBytes, 0644); err != nil {
		return nil, pkgErrors.Wrap(err, "service.UploadAvatar: os.WriteFile failed")
	}
	avatarURL := config.Conf.Static.BaseURL + config.Conf.Static.AvatarPath + fileName

	if err := mysql.UpdateUserAvatar(ctx, userID, avatarURL); err != nil {
		return nil, pkgErrors.Wrap(err, "service.UploadAvatar: UpdateUserAvatar failed")
	}

	u, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.UploadAvatar: GetUserByID failed")
	}

	hlog.CtxInfof(ctx, "上传头像成功, userID=%s", userID)
	return []*user.UserDataResp{{
		ID:        u.ID,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format(constants.DateTimeFormat),
		UpdatedAt: u.UpdatedAt.Format(constants.DateTimeFormat),
		DeletedAt: "",
	}}, nil
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

func Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	redisKey := constants.RefreshTokenPrefix + refreshToken
	userID, err := redis.RDB.Get(ctx, redisKey).Result()
	if err == redis2.Nil {
		return nil, errno.TokenInvalid
	} else if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Refresh: Redis Get failed")
	}

	newAccessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Refresh: GenerateAccessToken failed")
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.Refresh: GenerateRefreshToken failed")
	}

	redis.RDB.Del(ctx, redisKey)

	newKey := constants.RefreshTokenPrefix + newRefreshToken
	refreshExpire := time.Duration(config.Conf.Expire.RefreshToken) * time.Second
	if err := redis.RDB.Set(ctx, newKey, userID, refreshExpire).Err(); err != nil {
		return nil, pkgErrors.Wrap(err, "service.Refresh: Redis Set new token failed")
	}

	return &RefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

type GetMFAResult struct {
	Secret string
	Qrcode string
}

func GetMFA(ctx context.Context, userID string) (*GetMFAResult, error) {
	u, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.GetMFA: GetUserByID failed")
	}
	if u == nil {
		return nil, errno.UserNotExist
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "VideoWebsite",
		AccountName: u.Username,
		Period:      uint(config.Conf.JWT.MFAPeriod),
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.GetMFA: totp.Generate failed")
	}

	mfaKey := constants.MFAKeyPrefix + userID
	mfaExpire := time.Duration(config.Conf.Expire.MFA) * time.Second
	if err := redis.RDB.Set(ctx, mfaKey, key.Secret(), mfaExpire).Err(); err != nil {
		return nil, pkgErrors.Wrap(err, "service.GetMFA: Redis Set MFA secret failed")
	}

	qrCodeImage, err := key.Image(200, 200)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "service.GetMFA: key.Image failed")
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, qrCodeImage); err != nil {
		return nil, pkgErrors.Wrap(err, "service.GetMFA: png.Encode failed")
	}

	qrcodeBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	qrcodeURL := "data:image/png;base64," + qrcodeBase64

	return &GetMFAResult{
		Secret: key.Secret(),
		Qrcode: qrcodeURL,
	}, nil
}

func BindMFA(ctx context.Context, userID, code, secret string) error {
	if !totp.Validate(code, secret) {
		return errno.MFAError
	}

	if err := mysql.UpdateUserMFA(ctx, userID, secret); err != nil {
		return pkgErrors.Wrap(err, "service.BindMFA: UpdateUserMFA failed")
	}

	redisKey := constants.MFAKeyPrefix + userID
	redis.RDB.Del(ctx, redisKey)

	hlog.CtxInfof(ctx, "绑定MFA成功, userID=%s", userID)
	return nil
}
