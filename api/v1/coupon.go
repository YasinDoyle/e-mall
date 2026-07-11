package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

// AdminCouponCreateHandler 管理员创建优惠券
func AdminCouponCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCouponCreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetCouponSrv().AdminCouponCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCouponListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetCouponSrv().AdminCouponList(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCouponOfflineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCouponOfflineReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetCouponSrv().AdminCouponOffline(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// CouponListHandler 查询可领取的优惠券列表（公开）
func CouponListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetCouponSrv().CouponList(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// CouponClaimHandler 用户领券
func CouponClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CouponClaimReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetCouponSrv().CouponClaim(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// UserCouponListHandler 用户查看自己的优惠券
func UserCouponListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetCouponSrv().UserCouponList(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
