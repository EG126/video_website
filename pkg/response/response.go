package response

import (
	"video_website/biz/model/base"
	"video_website/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

func SendResponse(c *app.RequestContext, err error, data interface{}) {
	if err != nil {
		if e, ok := err.(errno.ErrNo); ok {
			c.JSON(200, map[string]interface{}{
				"base": &base.BaseResp{
					Code: e.Code,
					Msg:  e.Msg,
				},
				"data": data,
			})
			return
		}
		//未知错误
		c.JSON(200, map[string]interface{}{
			"base": &base.BaseResp{
				Code: errno.InternalServerError.Code,
				Msg:  err.Error(),
			},
			"data": data,
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"base": &base.BaseResp{
			Code: errno.Success.Code,
			Msg:  errno.Success.Msg,
		},
		"data": data,
	})
	return
}
