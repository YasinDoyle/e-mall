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

// ===== 分类 =====

func AdminCategoryCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCategoryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CategoryCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCategoryUpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCategoryUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CategoryUpdate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCategoryDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminIDReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CategoryDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 轮播图 =====

func AdminCarouselListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetAdminSrv().CarouselList(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCarouselCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCarouselReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		log.LogrusObj.Infof("imagepath:%s productId:%d", req.ImgPath, req.ProductID)
		resp, err := service.GetAdminSrv().CarouselCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCarouselUploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, fileHeader, _ := c.Request.FormFile("file")
		if fileHeader == nil {
			err := errors.New("请选择轮播图图片")
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CarouselUpload(c.Request.Context(), file, fileHeader.Size)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCarouselUpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminCarouselUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CarouselUpdate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminCarouselDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminIDReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().CarouselDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 公告 =====

func AdminNoticeListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetAdminSrv().NoticeList(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNoticeCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminNoticeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().NoticeCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNoticeUpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminNoticeUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().NoticeUpdate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNoticeDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminIDReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().NoticeDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 用户管理 =====

func AdminUserListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminListReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().UserList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminUserBanHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminUserBanReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().UserBan(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 商家管理 =====

func AdminSellerListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminSellerListReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().SellerList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminSellerAuditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminSellerAuditReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().SellerAudit(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 商品审核 =====

func AdminProductListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminProductListReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().ProductList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminProductAuditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminProductAuditReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().ProductAudit(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminProductDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminIDReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().ProductDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ===== 统计 =====

func AdminStatsOverviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetAdminSrv().StatsOverview(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminStatsOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.AdminStatsOrdersReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetAdminSrv().StatsOrders(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
