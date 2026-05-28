package ctl

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/track"
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
	if spanCtxInterface == nil {
		return "", errors.New("获取 track id 错误")
	}

	return track.TraceIDFromSpanContext(spanCtxInterface)
}
