package ctl

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

type Response struct {
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Msg     string `json:"msg"`
	Error   string `json:"error"`
	TrackId string `json:"track_id"`
}

func RespSuccess(c *gin.Context, data any, code ...int) *Response {
	trackId, _ := getTrackIdFromCtx(c)
	status := e.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = "操作成功"
	}

	r := &Response{
		Status:  status,
		Data:    data,
		Msg:     e.GetMsg(status),
		TrackId: trackId,
	}

	return r
}

func RespError(c *gin.Context, err error, data string, code ...int) *Response {
	trackId, _ := getTrackIdFromCtx(c)
	status := e.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status:  status,
		Msg:     e.GetMsg(status),
		Data:    data,
		Error:   err.Error(),
		TrackId: trackId,
	}

	return r
}

func getTrackIdFromCtx(c *gin.Context) (trackId string, err error) {
	spanCtxInterface, _ := c.Get(consts.SpanCTX)
	str := fmt.Sprintf("%v", spanCtxInterface)
	re := regexp.MustCompile(`([0-9a-fA-F]{16})`)

	match := re.FindStringSubmatch(str)
	if len(match) > 0 {
		return match[1], nil
	}
	return "", errors.New("获取 track id 错误")
}
