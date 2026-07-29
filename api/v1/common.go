package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func ErrorResponse(ctx *gin.Context, err error) *ctl.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", fieldError.Field()))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", fieldError.Tag()))
			return ctl.RespError(ctx, err, fmt.Sprintf("%s%s", field, tag))
		}
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return ctl.RespError(ctx, err, "JSON类型不匹配")
	}

	if code, ok := e.CodeFromError(err); ok {
		msg := e.GetMsgByLocale(code, requestLocale(ctx))
		resp := ctl.RespError(ctx, err, msg, code)
		resp.Msg = msg
		resp.MsgKey = e.GetMsgKey(code)
		resp.Error = msg
		return resp
	}

	return ctl.RespError(ctx, err, err.Error(), e.ERROR)
}

func requestLocale(ctx *gin.Context) string {
	if ctx == nil || ctx.Request == nil {
		return e.DefaultLocale
	}
	if locale := strings.TrimSpace(ctx.GetHeader("X-Locale")); locale != "" {
		return locale
	}
	acceptLanguage := strings.TrimSpace(ctx.GetHeader("Accept-Language"))
	if acceptLanguage == "" {
		return e.DefaultLocale
	}
	first := strings.Split(acceptLanguage, ",")[0]
	return strings.TrimSpace(strings.Split(first, ";")[0])
}
