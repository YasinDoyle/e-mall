package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

func UploadReviewImageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, fileHeader, _ := c.Request.FormFile("file")
		if fileHeader == nil {
			err := errors.New("请选择评价图片")
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetReviewSrv().ReviewImageUpload(c.Request.Context(), file, fileHeader.Size)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// CreateReviewHandler 用户创建评价
func CreateReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ReviewCreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetReviewSrv().ReviewCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ListReviewsHandler 获取商品评价列表（公开）
func ListReviewsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ReviewListReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetReviewSrv().ReviewList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// AdminDeleteReviewHandler 管理员删除评价
func AdminDeleteReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ReviewAdminDeleteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetReviewSrv().ReviewAdminDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
