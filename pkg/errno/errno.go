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
	Success             = ErrNo{Code: 0, Msg: "success"}
	ParamError          = ErrNo{Code: 10001, Msg: "参数错误"}
	UserNotExist        = ErrNo{Code: 10002, Msg: "用户不存在"}
	PasswordError       = ErrNo{Code: 10003, Msg: "密码错误"}
	UsernameExists      = ErrNo{Code: 10004, Msg: "用户名已存在"}
	Unauthorized        = ErrNo{Code: 10005, Msg: "未授权"}
	TokenInvalid        = ErrNo{Code: 10006, Msg: "token无效"}
	TokenExpired        = ErrNo{Code: 10007, Msg: "token过期"}
	Forbidden           = ErrNo{Code: 10008, Msg: "禁止访问"}
	DBError             = ErrNo{Code: 20001, Msg: "数据库错误"}
	FileTooLarge        = ErrNo{Code: 20002, Msg: "文件过大"}
	FileTypeError       = ErrNo{Code: 20003, Msg: "文件类型错误"}
	RedisError          = ErrNo{Code: 20004, Msg: "缓存错误"}
	InternalServerError = ErrNo{Code: 20005, Msg: "内部服务器错误"}
)
