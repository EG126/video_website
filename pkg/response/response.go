package response

import (
	"video_website/biz/model/base"
	"video_website/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

func SendResponse(c *app.RequestContext, err error, data interface{}) {
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"base": &base.BaseResp{
				Code: errno.GetCode(err),
				Msg:  err.Error(),
			},
			"data": data,
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"base": &base.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		"data": data,
	})
}
