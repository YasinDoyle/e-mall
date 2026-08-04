package event

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YasinDoyle/e-mall/consts"
	notificationhub "github.com/YasinDoyle/e-mall/domain/notification"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	esrepo "github.com/YasinDoyle/e-mall/repository/es"
)

func handleNotificationEvent(ctx context.Context, event any) error {
	switch e := event.(type) {
	case SellerApplied:
		return notifyAdmins(ctx, model.NotificationSceneSellerAudit, "新的商家入驻申请", fmt.Sprintf("用户 %d 提交了店铺「%s」的入驻申请。", e.SellerID, e.ShopName), map[string]any{
			"seller_id": e.SellerID,
			"shop_name": e.ShopName,
		})
	case SellerAuditChanged:
		return notifySellerAuditChanged(ctx, e)
	case ProductSubmitted:
		if e.Product == nil {
			return nil
		}
		return notifyAdmins(ctx, model.NotificationSceneProductAudit, "新的商品审核待办", fmt.Sprintf("商家 %s 提交了商品「%s」。", e.Product.BossName, e.Product.Name), map[string]any{
			"product_id": e.Product.ID,
			"seller_id":  e.Product.BossID,
		})
	case ProductAuditChanged:
		return notifyProductAuditChanged(ctx, e)
	case OrderPaid:
		if err := notifyUser(ctx, e.BuyerID, model.NotificationSceneOrderPaid, "订单支付成功", fmt.Sprintf("订单 %d 已支付成功。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
		}); err != nil {
			return err
		}
		return notifyUser(ctx, e.SellerID, model.NotificationSceneOrderPaid, "你有新的待发货订单", fmt.Sprintf("订单 %d 已支付，请及时发货。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
		})
	case OrderShipped:
		return notifyUser(ctx, e.BuyerID, model.NotificationSceneOrderShipped, "订单已发货", fmt.Sprintf("订单 %d 已发货，物流单号：%s。", e.OrderNum, e.TrackingNo), map[string]any{
			"order_id":    e.OrderID,
			"order_num":   e.OrderNum,
			"tracking_no": e.TrackingNo,
		})
	case RefundRequested:
		if err := notifyUser(ctx, e.SellerID, model.NotificationSceneOrderRefunded, "买家申请退款", fmt.Sprintf("订单 %d 有新的退款申请。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
		}); err != nil {
			return err
		}
		return notifyAdmins(ctx, model.NotificationSceneOrderRefunded, "新的退款申请", fmt.Sprintf("订单 %d 等待平台处理退款。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
		})
	case OrderRefunded:
		if err := notifyUser(ctx, e.BuyerID, model.NotificationSceneOrderRefunded, "退款已完成", fmt.Sprintf("订单 %d 已退款，退款金额 %.2f。", e.OrderNum, e.RefundAmount), map[string]any{
			"order_id":      e.OrderID,
			"order_num":     e.OrderNum,
			"refund_amount": e.RefundAmount,
		}); err != nil {
			return err
		}
		return notifyUser(ctx, e.SellerID, model.NotificationSceneOrderRefunded, "订单已退款", fmt.Sprintf("订单 %d 已完成退款。", e.OrderNum), map[string]any{
			"order_id":      e.OrderID,
			"order_num":     e.OrderNum,
			"refund_amount": e.RefundAmount,
		})
	case AfterSaleClosed:
		if err := notifyUser(ctx, e.BuyerID, model.NotificationSceneOrderRefunded, "售后已关闭", fmt.Sprintf("订单 %d 的售后已关闭：%s。", e.OrderNum, e.Note), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
			"note":      e.Note,
		}); err != nil {
			return err
		}
		return notifyUser(ctx, e.SellerID, model.NotificationSceneOrderRefunded, "售后已关闭", fmt.Sprintf("订单 %d 的售后已由平台关闭。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
			"note":      e.Note,
		})
	case OrderReceived:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneSettlement, "订单已确认收货", fmt.Sprintf("订单 %d 已确认收货，结算金额已进入可结算。", e.OrderNum), map[string]any{
			"order_id":  e.OrderID,
			"order_num": e.OrderNum,
		})
	case SettlementGenerated:
		if e.SettlementID > 0 {
			return notifyUser(ctx, e.SellerID, model.NotificationSceneSettlement, "结算单已生成", fmt.Sprintf("结算单 #%d 已生成，可结算金额 %.2f。", e.SettlementID, e.Amount), map[string]any{
				"settlement_id": e.SettlementID,
				"order_id":      e.OrderID,
			})
		}
		if e.Count <= 0 {
			return nil
		}
		return notifyUser(ctx, e.SellerID, model.NotificationSceneSettlement, "结算单已生成", fmt.Sprintf("平台已生成 %d 笔可结算订单。", e.Count), map[string]any{
			"seller_id": e.SellerID,
			"count":     e.Count,
		})
	case SettlementPaid:
		if e.Settlement == nil {
			return nil
		}
		return notifyUser(ctx, e.Settlement.SellerID, model.NotificationSceneSettlement, "结算已打款", fmt.Sprintf("结算单 #%d 已打款，金额 %.2f。", e.Settlement.ID, e.Settlement.SettlementAmount), map[string]any{
			"settlement_id": e.Settlement.ID,
			"order_id":      e.Settlement.OrderID,
		})
	case WithdrawApplied:
		if e.Withdraw == nil {
			return nil
		}
		return notifyAdmins(ctx, model.NotificationSceneWithdraw, "新的提现申请", fmt.Sprintf("商家 %s 申请提现 %.2f。", e.ShopName, e.Withdraw.Amount), map[string]any{
			"withdraw_id": e.Withdraw.ID,
			"seller_id":   e.Withdraw.SellerID,
			"amount":      e.Withdraw.Amount,
		})
	case WithdrawAuditChanged:
		return notifyWithdrawAuditChanged(ctx, e)
	case WithdrawPaidStatusChanged:
		return notifyWithdrawPaidStatusChanged(ctx, e)
	default:
		return nil
	}
}

func handleProductIndexEvent(ctx context.Context, event any) error {
	switch e := event.(type) {
	case ProductSubmitted:
		return esrepo.NewProductIndexRepo().IndexProduct(ctx, e.Product)
	case ProductChanged:
		return esrepo.NewProductIndexRepo().IndexProduct(ctx, e.Product)
	case ProductDeleted:
		return esrepo.NewProductIndexRepo().DeleteProduct(ctx, e.ProductID)
	default:
		return nil
	}
}

func notifySellerAuditChanged(ctx context.Context, e SellerAuditChanged) error {
	switch e.Status {
	case consts.SellerStatusApproved:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneSellerAudit, "商家入驻审核通过", "你的店铺已通过审核，可以发布商品并使用卖家中心。", map[string]any{
			"status": e.Status,
		})
	case consts.SellerStatusRejected:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneSellerAudit, "商家入驻审核未通过", e.RejectReason, map[string]any{
			"status":        e.Status,
			"reject_reason": e.RejectReason,
		})
	case consts.SellerStatusBanned:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneSellerAudit, "商家账号已封禁", "你的卖家能力已被后台封禁。", map[string]any{
			"status": e.Status,
		})
	default:
		return nil
	}
}

func notifyProductAuditChanged(ctx context.Context, e ProductAuditChanged) error {
	switch e.AuditStatus {
	case consts.ProductAuditApproved:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneProductAudit, "商品审核通过", fmt.Sprintf("商品「%s」已通过审核并上架。", e.ProductName), map[string]any{
			"product_id": e.ProductID,
			"status":     e.AuditStatus,
		})
	case consts.ProductAuditRejected:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneProductAudit, "商品审核未通过", fmt.Sprintf("商品「%s」未通过审核，请修改后重新提交。", e.ProductName), map[string]any{
			"product_id": e.ProductID,
			"status":     e.AuditStatus,
		})
	default:
		return nil
	}
}

func notifyWithdrawAuditChanged(ctx context.Context, e WithdrawAuditChanged) error {
	switch e.Status {
	case model.SellerWithdrawStatusApproved:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneWithdraw, "提现审核通过", fmt.Sprintf("提现申请 #%d 已通过审核，等待打款。", e.WithdrawID), map[string]any{
			"withdraw_id": e.WithdrawID,
			"amount":      e.Amount,
		})
	case model.SellerWithdrawStatusRejected:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneWithdraw, "提现审核未通过", e.Reason, map[string]any{
			"withdraw_id": e.WithdrawID,
			"amount":      e.Amount,
			"reason":      e.Reason,
		})
	default:
		return nil
	}
}

func notifyWithdrawPaidStatusChanged(ctx context.Context, e WithdrawPaidStatusChanged) error {
	switch e.Status {
	case model.SellerWithdrawStatusPaid:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneWithdraw, "提现已打款", fmt.Sprintf("提现申请 #%d 已标记打款，金额 %.2f。", e.WithdrawID, e.Amount), map[string]any{
			"withdraw_id": e.WithdrawID,
			"amount":      e.Amount,
		})
	case model.SellerWithdrawStatusFailed:
		return notifyUser(ctx, e.SellerID, model.NotificationSceneWithdraw, "提现打款失败", e.Reason, map[string]any{
			"withdraw_id": e.WithdrawID,
			"amount":      e.Amount,
			"reason":      e.Reason,
		})
	default:
		return nil
	}
}

func notifyUser(ctx context.Context, recipientID uint, scene, title, content string, payload any) error {
	if recipientID == 0 {
		return nil
	}
	if err := dao.NewNotificationDao(ctx).Create(&model.Notification{
		RecipientType: model.NotificationRecipientUser,
		RecipientID:   recipientID,
		Scene:         strings.TrimSpace(scene),
		Title:         strings.TrimSpace(title),
		Content:       strings.TrimSpace(content),
		Payload:       marshalPayload(payload),
	}); err != nil {
		return err
	}
	notificationhub.Publish(model.NotificationRecipientUser, recipientID)
	return nil
}

func notifyAdmins(ctx context.Context, scene, title, content string, payload any) error {
	admins, err := dao.NewUserDao(ctx).ListAdminUsers()
	if err != nil {
		return err
	}
	notifications := make([]*model.Notification, 0, len(admins))
	for _, admin := range admins {
		notifications = append(notifications, &model.Notification{
			RecipientType: model.NotificationRecipientAdmin,
			RecipientID:   admin.ID,
			Scene:         strings.TrimSpace(scene),
			Title:         strings.TrimSpace(title),
			Content:       strings.TrimSpace(content),
			Payload:       marshalPayload(payload),
		})
	}
	if err := dao.NewNotificationDao(ctx).BatchCreate(notifications); err != nil {
		return err
	}
	for _, admin := range admins {
		notificationhub.Publish(model.NotificationRecipientAdmin, admin.ID)
	}
	return nil
}

func marshalPayload(payload any) string {
	if payload == nil {
		return ""
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
