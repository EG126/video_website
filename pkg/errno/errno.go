package errno

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

var (
	Success          = errors.New("success")
	ParamError       = errors.New("参数错误")
	UserNotExist     = errors.New("用户不存在")
	PasswordError    = errors.New("密码错误")
	UsernameExists   = errors.New("用户名已存在")
	Unauthorized     = errors.New("未授权")
	TokenInvalid     = errors.New("token无效")
	TokenExpired     = errors.New("token过期")
	TokenError       = errors.New("token生成失败")
	Forbidden        = errors.New("禁止访问")
	VideoNotExist    = errors.New("视频不存在")
	CommentNotExist  = errors.New("评论不存在")
	RelationNotExist = errors.New("关系不存在")
	MFARequired      = errors.New("需要MFA验证码")
	MFAError         = errors.New("MFA验证码错误")

	DBError             = errors.New("数据库错误")
	FileTooLarge        = errors.New("文件过大")
	FileTypeError       = errors.New("文件类型错误")
	FileError           = errors.New("文件操作失败")
	RedisError          = errors.New("缓存错误")
	EncryptError        = errors.New("加密操作失败")
	InternalServerError = errors.New("内部服务器错误")
)

var errorCodeMap = map[error]int32{
	Success:          0,
	ParamError:       10001,
	UserNotExist:     10002,
	PasswordError:    10003,
	UsernameExists:   10004,
	Unauthorized:     10005,
	TokenInvalid:     10006,
	TokenExpired:     10007,
	TokenError:       10008,
	Forbidden:        10009,
	VideoNotExist:    10010,
	CommentNotExist:  10011,
	RelationNotExist: 10012,
	MFARequired:      10013,
	MFAError:         10014,

	DBError:             20001,
	FileTooLarge:        20002,
	FileTypeError:       20003,
	FileError:           20004,
	RedisError:          20005,
	EncryptError:        20006,
	InternalServerError: 20007,
}

func GetCode(err error) int32 {
	originalErr := Cause(err)
	if code, ok := errorCodeMap[originalErr]; ok {
		return code
	}
	if code, ok := errorCodeMap[err]; ok {
		return code
	}
	return 20001 // 默认返回内部错误
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Log(ctx context.Context, msg string, err error) {
	originalErr := Cause(err)
	hlog.CtxErrorf(ctx, "%s: original error: %T %v", msg, originalErr, originalErr)
	hlog.CtxErrorf(ctx, "stack trace:\n%+v", err)
}
