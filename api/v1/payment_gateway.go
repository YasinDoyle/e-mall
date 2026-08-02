package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/application"
	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

func OrderBalancePayHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PaymentDownReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := application.NewPaymentUsecase().PayOrderByBalance(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func OrderWechatPayHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderGatewayPayReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().WechatOrderPay(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func OrderAlipayPayHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderGatewayPayReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().AlipayOrderPay(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func OrderPaymentStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentNo := c.Query("payment_no")
		if paymentNo == "" {
			c.JSON(http.StatusOK, ErrorResponse(c, errors.New("payment_no不能为空")))
			return
		}
		resp, err := service.GetPayGatewaySrv().OrderPaymentStatus(c.Request.Context(), paymentNo)
		if err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// WechatRechargeHandler 微信充值发起
func WechatRechargeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RechargeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().WechatRecharge(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// AlipayRechargeHandler 支付宝充值发起
func AlipayRechargeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RechargeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().AlipayRecharge(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// RechargeStatusHandler 查询充值订单状态（前端轮询）
func RechargeStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNum := c.Query("order_num")
		if orderNum == "" {
			c.JSON(http.StatusOK, ErrorResponse(c, errors.New("order_num不能为空")))
			return
		}
		resp, err := service.GetPayGatewaySrv().RechargeStatus(c.Request.Context(), orderNum)
		if err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// GetPendingCreditHandler 查询待入账金额
func GetPendingCreditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		amount, err := service.GetPayGatewaySrv().GetPendingCredit(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, gin.H{"pending": amount}))
	}
}

// ApplyPendingCreditHandler 用户确认入账（需要支付密码）
func ApplyPendingCreditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Key string `json:"key" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().ApplyPendingCredit(c.Request.Context(), req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// WechatNotifyHandler 微信支付回调（无需 JWT）
func WechatNotifyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.GetPayGatewaySrv().HandleWechatNotifyRequest(c.Request.Context(), c.Request); err != nil {
			log.LogrusObj.Error(err)
			c.String(http.StatusInternalServerError, "FAIL")
			return
		}
		c.String(http.StatusOK, "SUCCESS")
	}
}

// AlipayNotifyHandler 支付宝回调（无需 JWT）
func AlipayNotifyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.GetPayGatewaySrv().HandleAlipayNotifyRequest(c.Request.Context(), c.Request); err != nil {
			log.LogrusObj.Error(err)
			c.String(http.StatusInternalServerError, "fail")
			return
		}
		c.String(http.StatusOK, "success")
	}
}

func WechatRefundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RechargeRefundReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().WechatRefund(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func AlipayRefundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.RechargeRefundReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		resp, err := service.GetPayGatewaySrv().AlipayRefund(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusOK, ErrorResponse(c, err))
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
