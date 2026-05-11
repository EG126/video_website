package errno

import "fmt"

type ErrNo struct {
	Code int32
	Msg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

var (
	Success          = ErrNo{Code: 0, Msg: "success"}
	ParamError       = ErrNo{Code: 10001, Msg: "参数错误"}
	UserNotExist     = ErrNo{Code: 10002, Msg: "用户不存在"}
	PasswordError    = ErrNo{Code: 10003, Msg: "密码错误"}
	UsernameExists   = ErrNo{Code: 10004, Msg: "用户名已存在"}
	Unauthorized     = ErrNo{Code: 10005, Msg: "未授权"}
	TokenInvalid     = ErrNo{Code: 10006, Msg: "token无效"}
	TokenExpired     = ErrNo{Code: 10007, Msg: "token过期"}
	TokenError       = ErrNo{Code: 10008, Msg: "token生成失败"}
	Forbidden        = ErrNo{Code: 10009, Msg: "禁止访问"}
	VideoNotExist    = ErrNo{Code: 10010, Msg: "视频不存在"}
	CommentNotExist  = ErrNo{Code: 10011, Msg: "评论不存在"}
	RelationNotExist = ErrNo{Code: 10012, Msg: "关系不存在"}

	DBError             = ErrNo{Code: 20001, Msg: "数据库错误"}
	FileTooLarge        = ErrNo{Code: 20002, Msg: "文件过大"}
	FileTypeError       = ErrNo{Code: 20003, Msg: "文件类型错误"}
	FileError           = ErrNo{Code: 20004, Msg: "文件操作失败"}
	RedisError          = ErrNo{Code: 20005, Msg: "缓存错误"}
	EncryptError        = ErrNo{Code: 20006, Msg: "加密操作失败"}
	InternalServerError = ErrNo{Code: 20007, Msg: "内部服务器错误"}
)
