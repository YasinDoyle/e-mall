package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

func SellerApplyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.SellerApplyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetSellerSrv().Apply(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func SellerProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetSellerSrv().Profile(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
