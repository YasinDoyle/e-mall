package v1

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func ErrorResponse(ctx *gin.Context, err error) *ctl.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", fieldError.Field))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", fieldError.Tag))
			return ctl.RespError(ctx, err, fmt.Sprintf("%s%s", field, tag))
		}
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return ctl.RespError(ctx, err, "JSON类型不匹配")
	}

	return ctl.RespError(ctx, err, err.Error(), e.ERROR)
}
