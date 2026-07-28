package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

func NotificationListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.NotificationListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetNotificationSrv().UserList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func NotificationUnreadCountHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetNotificationSrv().UserUnreadCount(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func NotificationMarkReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.NotificationMarkReadReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetNotificationSrv().UserMarkRead(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func NotificationMarkAllReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetNotificationSrv().UserMarkAllRead(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func NotificationStreamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		service.GetNotificationSrv().StreamUnreadCount(c, model.NotificationRecipientUser)
	}
}

func AdminNotificationListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.NotificationListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetNotificationSrv().AdminList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNotificationUnreadCountHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetNotificationSrv().AdminUnreadCount(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNotificationMarkReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.NotificationMarkReadReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetNotificationSrv().AdminMarkRead(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNotificationMarkAllReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetNotificationSrv().AdminMarkAllRead(c.Request.Context())
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AdminNotificationStreamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		service.GetNotificationSrv().StreamUnreadCount(c, model.NotificationRecipientAdmin)
	}
}
