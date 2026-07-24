package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	wechat "github.com/go-pay/gopay/wechat/v3"
	"gorm.io/gorm"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/repository/rabbitmq"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

const pendingCreditPrefix = "pending_credit:"
const pendingCreditExpire = 7 * 24 * time.Hour

var PayGatewaySrvIns *PayGatewaySrv
var PayGatewaySrvOnce sync.Once

type PayGatewaySrv struct{}

func GetPayGatewaySrv() *PayGatewaySrv {
	PayGatewaySrvOnce.Do(func() { PayGatewaySrvIns = &PayGatewaySrv{} })
	return PayGatewaySrvIns
}

func pendingCreditKey(userID uint) string { return fmt.Sprintf("%s%d", pendingCreditPrefix, userID) }

func genRechargeOrderNum() string { return fmt.Sprintf("R%d", time.Now().UnixNano()) }

func genRefundNo(orderNum string) string {
	return fmt.Sprintf("RF%s%d", orderNum, time.Now().UnixNano())
}

func amountToFen(amount float64) int {
	return int(math.Round(amount * 100))
}

func fenToYuan(amount int) float64 {
	return float64(amount) / 100
}

func newWechatClient() (*wechat.ClientV3, error) {
	wConf := conf.Config.WechatPay
	if wConf == nil || wConf.AppID == "" || wConf.MchID == "" || wConf.ApiV3Key == "" || wConf.SerialNo == "" || wConf.PrivateKey == "" {
		return nil, errors.New("微信支付未配置，请联系管理员")
	}

	client, err := wechat.NewClientV3(wConf.MchID, wConf.SerialNo, wConf.ApiV3Key, wConf.PrivateKey)
	if err != nil {
		return nil, err
	}
	client.DebugSwitch = gopay.DebugOff
	if wConf.PublicKey != "" && wConf.PublicKeyID != "" {
		if err = client.AutoVerifySignByPublicKey([]byte(wConf.PublicKey), wConf.PublicKeyID); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func newAlipayClient() (*alipay.Client, error) {
	aConf := conf.Config.Alipay
	if aConf == nil || aConf.AppID == "" || aConf.PrivateKey == "" {
		return nil, errors.New("支付宝未配置，请联系管理员")
	}

	client, err := alipay.NewClient(aConf.AppID, aConf.PrivateKey, !aConf.IsSandbox)
	if err != nil {
		return nil, err
	}
	client.SetCharset("utf-8").SetSignType(alipay.RSA2)
	if aConf.AlipayPublicKey != "" {
		client.AutoVerifySign([]byte(aConf.AlipayPublicKey))
	}

	return client, nil
}

func createRechargeOrder(ctx context.Context, userID uint, channel string, amount float64) (*model.RechargeOrder, error) {
	order := &model.RechargeOrder{
		OrderNum:     genRechargeOrderNum(),
		UserID:       userID,
		Channel:      channel,
		Amount:       amount,
		Status:       model.RechargeStatusPending,
		RefundStatus: model.RechargeRefundNone,
	}
	if err := dao.NewRechargeDao(ctx).Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

// WechatRecharge 发起微信 Native 充值。
func (s *PayGatewaySrv) WechatRecharge(ctx context.Context, req *types.RechargeReq) (resp *types.RechargeResp, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	if amountToFen(req.Amount) <= 0 {
		return nil, errors.New("充值金额不能小于 0.01 元")
	}

	client, err := newWechatClient()
	if err != nil {
		return nil, err
	}
	order, err := createRechargeOrder(ctx, u.Id, model.RechargeChannelWechat, req.Amount)
	if err != nil {
		return nil, err
	}

	wConf := conf.Config.WechatPay
	bm := make(gopay.BodyMap)
	bm.Set("appid", wConf.AppID).
		Set("mchid", wConf.MchID).
		Set("description", "E-Mall 钱包充值").
		Set("out_trade_no", order.OrderNum).
		Set("notify_url", wConf.NotifyURL).
		SetBodyMap("amount", func(b gopay.BodyMap) {
			b.Set("total", amountToFen(req.Amount)).Set("currency", "CNY")
		})

	wxRsp, err := client.V3TransactionNative(ctx, bm)
	if err != nil || wxRsp == nil || wxRsp.Response == nil || wxRsp.Code != wechat.Success {
		_ = dao.NewRechargeDao(ctx).MarkFailed(order.OrderNum)
		if err != nil {
			log.LogrusObj.Error(err)
			return nil, errors.New("微信支付下单失败，请检查配置")
		}
		if wxRsp == nil {
			return nil, errors.New("微信支付下单失败")
		}
		return nil, fmt.Errorf("微信支付下单失败: %s", wxRsp.Error)
	}

	resp = &types.RechargeResp{OrderNum: order.OrderNum, QRCodeURL: wxRsp.Response.CodeUrl}
	return
}

// AlipayRecharge 发起支付宝 PC/WAP 充值。
func (s *PayGatewaySrv) AlipayRecharge(ctx context.Context, req *types.RechargeReq) (resp *types.RechargeResp, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	if amountToFen(req.Amount) <= 0 {
		return nil, errors.New("充值金额不能小于 0.01 元")
	}

	client, err := newAlipayClient()
	if err != nil {
		return nil, err
	}
	order, err := createRechargeOrder(ctx, u.Id, model.RechargeChannelAlipay, req.Amount)
	if err != nil {
		return nil, err
	}

	aConf := conf.Config.Alipay
	bm := make(gopay.BodyMap)
	bm.Set("subject", "E-Mall 钱包充值").
		Set("out_trade_no", order.OrderNum).
		Set("total_amount", fmt.Sprintf("%.2f", req.Amount)).
		Set("return_url", aConf.ReturnURL).
		Set("notify_url", aConf.NotifyURL)

	payURL, err := client.TradePagePay(ctx, bm)
	if err != nil {
		_ = dao.NewRechargeDao(ctx).MarkFailed(order.OrderNum)
		log.LogrusObj.Error(err)
		return nil, errors.New("支付宝下单失败，请检查配置")
	}

	resp = &types.RechargeResp{OrderNum: order.OrderNum, PayURL: payURL}
	return
}

func (s *PayGatewaySrv) RechargeStatus(ctx context.Context, orderNum string) (*types.RechargeStatusResp, error) {
	order, err := dao.NewRechargeDao(ctx).GetByOrderNum(orderNum)
	if err != nil {
		return nil, err
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	if order.UserID != u.Id {
		return nil, errors.New("无权查看该充值订单")
	}
	pending := float64(0)
	if order.Status == model.RechargeStatusPaid {
		pending, err = dao.NewRechargeDao(ctx).SumPendingCreditByUser(order.UserID)
		if err != nil {
			return nil, err
		}
	}

	return &types.RechargeStatusResp{
		OrderNum:      order.OrderNum,
		Status:        order.Status,
		Amount:        order.Amount,
		Channel:       order.Channel,
		PendingCredit: pending,
		RefundStatus:  order.RefundStatus,
	}, nil
}

func (s *PayGatewaySrv) GetPendingCredit(ctx context.Context) (float64, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return 0, err
	}
	return dao.NewRechargeDao(ctx).SumPendingCreditByUser(u.Id)
}

func (s *PayGatewaySrv) ApplyPendingCredit(ctx context.Context, key string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	pendingKey := pendingCreditKey(u.Id)
	pending := float64(0)
	err = dao.NewRechargeDao(ctx).DB.Transaction(func(tx *gorm.DB) error {
		rechargeDao := dao.NewRechargeDaoByDB(tx)
		orders, txErr := rechargeDao.ListPendingCreditOrdersByUser(u.Id)
		if txErr != nil {
			return txErr
		}
		if len(orders) == 0 {
			return errors.New("暂无待入账余额")
		}

		ids := make([]uint, 0, len(orders))
		for _, order := range orders {
			pending += order.Amount
			ids = append(ids, order.ID)
		}
		if pending <= 0 {
			return errors.New("暂无待入账余额")
		}

		userDao := dao.NewUserDaoByDB(tx)
		user, txErr := userDao.GetUserById(u.Id)
		if txErr != nil {
			return txErr
		}
		if !user.HasPayKey() {
			return errors.New("请先设置支付密码")
		}

		current, txErr := user.DecryptMoney(key)
		if txErr != nil {
			return errors.New("支付密码错误")
		}

		user.Money = fmt.Sprintf("%f", current+pending)
		user.Money, txErr = user.EncryptMoney(key)
		if txErr != nil {
			return txErr
		}
		if txErr = userDao.UpdateUserById(u.Id, user); txErr != nil {
			return txErr
		}
		return rechargeDao.MarkOrdersCredited(u.Id, ids)
	})
	if err != nil {
		return nil, err
	}
	cache.RedisClient.Del(ctx, pendingKey)
	log.LogrusObj.Infof("用户 %d 入账 %.2f 成功", u.Id, pending)

	resp = fmt.Sprintf("入账成功，到账 ¥%.2f", pending)
	return
}

func (s *PayGatewaySrv) markRechargePaid(ctx context.Context, orderNum, providerTradeNo string, paidAt time.Time, paidAmount float64) error {
	order, err := dao.NewRechargeDao(ctx).GetByOrderNum(orderNum)
	if err != nil {
		return err
	}
	if order.Status == model.RechargeStatusPaid || order.Status == model.RechargeStatusCredited {
		return nil
	}
	if order.Status != model.RechargeStatusPending {
		return errors.New("充值订单状态不允许支付回调")
	}
	if math.Abs(order.Amount-paidAmount) > 0.001 {
		return errors.New("充值回调金额不匹配")
	}

	paidOrder, freshPaid, err := dao.NewRechargeDao(ctx).MarkPaid(orderNum, providerTradeNo, paidAt)
	if err != nil {
		return err
	}
	if !freshPaid {
		return nil
	}

	event := &types.RechargePaidEvent{
		OrderNum:        paidOrder.OrderNum,
		UserID:          paidOrder.UserID,
		Channel:         paidOrder.Channel,
		Amount:          paidOrder.Amount,
		ProviderTradeNo: providerTradeNo,
		PaidAt:          paidAt,
	}
	if publishErr := rabbitmq.PublishJSON(ctx, consts.RechargePaidQueue, event); publishErr != nil {
		log.LogrusObj.Error(publishErr)
	}
	if err = refreshPendingCredit(ctx, paidOrder.UserID); err != nil {
		log.LogrusObj.Error(err)
	}

	log.LogrusObj.Infof("充值回调成功 user=%d amount=%.2f orderNum=%s", paidOrder.UserID, paidOrder.Amount, orderNum)
	return nil
}

func refreshPendingCredit(ctx context.Context, userID uint) error {
	pending, err := dao.NewRechargeDao(ctx).SumPendingCreditByUser(userID)
	if err != nil {
		return err
	}
	key := pendingCreditKey(userID)
	if pending <= 0 {
		return cache.RedisClient.Del(ctx, key).Err()
	}
	return cache.RedisClient.Set(ctx, key, fmt.Sprintf("%f", pending), pendingCreditExpire).Err()
}

func (s *PayGatewaySrv) HandleWechatNotifyRequest(ctx context.Context, req *http.Request) error {
	wConf := conf.Config.WechatPay
	if wConf == nil {
		return errors.New("微信支付未配置")
	}

	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		if wConf.IsSandbox {
			orderNum := req.URL.Query().Get("order_num")
			if orderNum == "" {
				return err
			}
			return s.markRechargePaid(ctx, orderNum, "wechat-sandbox", time.Now(), mustRechargeAmount(ctx, orderNum))
		}
		return err
	}

	if wConf.PublicKey != "" && wConf.PublicKeyID != "" {
		client, clientErr := newWechatClient()
		if clientErr != nil {
			return clientErr
		}
		if err = notifyReq.VerifySignByPKMap(client.WxPublicKeyMap()); err != nil {
			return err
		}
	} else if !wConf.IsSandbox {
		return errors.New("微信回调验签公钥未配置")
	}

	payResult, err := notifyReq.DecryptPayCipherText(wConf.ApiV3Key)
	if err != nil {
		return err
	}
	if payResult.TradeState != "SUCCESS" {
		return nil
	}

	paidAt := time.Now()
	if payResult.SuccessTime != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, payResult.SuccessTime); parseErr == nil {
			paidAt = parsed
		}
	}
	amount := float64(0)
	if payResult.Amount != nil {
		amount = fenToYuan(payResult.Amount.PayerTotal)
	}

	return s.markRechargePaid(ctx, payResult.OutTradeNo, payResult.TransactionId, paidAt, amount)
}

func mustRechargeAmount(ctx context.Context, orderNum string) float64 {
	order, err := dao.NewRechargeDao(ctx).GetByOrderNum(orderNum)
	if err != nil {
		return 0
	}
	return order.Amount
}

func (s *PayGatewaySrv) HandleAlipayNotifyRequest(ctx context.Context, req *http.Request) error {
	aConf := conf.Config.Alipay
	if aConf == nil {
		return errors.New("支付宝未配置")
	}

	notifyMap, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return err
	}
	if aConf.AlipayPublicKey != "" {
		ok, verifyErr := alipay.VerifySign(aConf.AlipayPublicKey, notifyMap)
		if verifyErr != nil {
			return verifyErr
		}
		if !ok {
			return errors.New("支付宝回调验签失败")
		}
	} else if !aConf.IsSandbox {
		return errors.New("支付宝回调验签公钥未配置")
	}

	tradeStatus := notifyMap.Get("trade_status")
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		return nil
	}
	amount, _ := strconv.ParseFloat(notifyMap.Get("total_amount"), 64)
	paidAt := time.Now()
	if notifyMap.Get("gmt_payment") != "" {
		if parsed, parseErr := time.ParseInLocation("2006-01-02 15:04:05", notifyMap.Get("gmt_payment"), time.Local); parseErr == nil {
			paidAt = parsed
		}
	}

	return s.markRechargePaid(ctx, notifyMap.Get("out_trade_no"), notifyMap.Get("trade_no"), paidAt, amount)
}

func (s *PayGatewaySrv) WechatRefund(ctx context.Context, req *types.RechargeRefundReq) (*types.RechargeRefundResp, error) {
	order, err := dao.NewRechargeDao(ctx).GetByOrderNum(req.OrderNum)
	if err != nil {
		return nil, err
	}
	if order.Channel != model.RechargeChannelWechat {
		return nil, errors.New("充值订单不是微信渠道")
	}
	if order.Status != model.RechargeStatusPaid {
		return nil, errors.New("仅已支付且未入账的充值订单可退款")
	}
	if req.Amount <= 0 || math.Abs(req.Amount-order.Amount) > 0.001 {
		return nil, errors.New("充值退款需整笔退回")
	}

	client, err := newWechatClient()
	if err != nil {
		return nil, err
	}
	refundNo := genRefundNo(order.OrderNum)
	reason := req.Reason
	if reason == "" {
		reason = "钱包充值退款"
	}
	if err = dao.NewRechargeDao(ctx).MarkRefundProcessing(order.OrderNum, refundNo, reason, req.Amount); err != nil {
		return nil, err
	}
	if err = refreshPendingCredit(ctx, order.UserID); err != nil {
		log.LogrusObj.Error(err)
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", order.OrderNum).
		Set("out_refund_no", refundNo).
		Set("reason", reason).
		SetBodyMap("amount", func(b gopay.BodyMap) {
			b.Set("refund", amountToFen(req.Amount)).
				Set("total", amountToFen(order.Amount)).
				Set("currency", "CNY")
		})
	wxRsp, err := client.V3Refund(ctx, bm)
	if err != nil || wxRsp == nil || wxRsp.Response == nil || wxRsp.Code != wechat.Success {
		_ = dao.NewRechargeDao(ctx).MarkRefundResult(order.OrderNum, refundNo, model.RechargeRefundFailed, nil)
		_ = refreshPendingCredit(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		if wxRsp == nil {
			return nil, errors.New("微信退款失败")
		}
		return nil, fmt.Errorf("微信退款失败: %s", wxRsp.Error)
	}

	status := model.RechargeRefundProcessing
	if wxRsp.Response.Status == "SUCCESS" {
		status = model.RechargeRefundSuccess
		now := time.Now()
		_ = dao.NewRechargeDao(ctx).MarkRefundResult(order.OrderNum, refundNo, status, &now)
	} else {
		_ = dao.NewRechargeDao(ctx).MarkRefundResult(order.OrderNum, refundNo, status, nil)
	}
	if err = refreshPendingCredit(ctx, order.UserID); err != nil {
		log.LogrusObj.Error(err)
	}

	return &types.RechargeRefundResp{
		OrderNum:     order.OrderNum,
		RefundNo:     refundNo,
		Amount:       req.Amount,
		RefundStatus: status,
		ProviderID:   wxRsp.Response.RefundId,
	}, nil
}

func (s *PayGatewaySrv) AlipayRefund(ctx context.Context, req *types.RechargeRefundReq) (*types.RechargeRefundResp, error) {
	order, err := dao.NewRechargeDao(ctx).GetByOrderNum(req.OrderNum)
	if err != nil {
		return nil, err
	}
	if order.Channel != model.RechargeChannelAlipay {
		return nil, errors.New("充值订单不是支付宝渠道")
	}
	if order.Status != model.RechargeStatusPaid {
		return nil, errors.New("仅已支付且未入账的充值订单可退款")
	}
	if req.Amount <= 0 || math.Abs(req.Amount-order.Amount) > 0.001 {
		return nil, errors.New("充值退款需整笔退回")
	}

	client, err := newAlipayClient()
	if err != nil {
		return nil, err
	}
	refundNo := genRefundNo(order.OrderNum)
	reason := req.Reason
	if reason == "" {
		reason = "钱包充值退款"
	}
	if err = dao.NewRechargeDao(ctx).MarkRefundProcessing(order.OrderNum, refundNo, reason, req.Amount); err != nil {
		return nil, err
	}
	if err = refreshPendingCredit(ctx, order.UserID); err != nil {
		log.LogrusObj.Error(err)
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", order.OrderNum).
		Set("out_request_no", refundNo).
		Set("refund_amount", fmt.Sprintf("%.2f", req.Amount)).
		Set("refund_reason", reason)
	aliRsp, err := client.TradeRefund(ctx, bm)
	if err != nil || aliRsp == nil || aliRsp.Response == nil {
		_ = dao.NewRechargeDao(ctx).MarkRefundResult(order.OrderNum, refundNo, model.RechargeRefundFailed, nil)
		_ = refreshPendingCredit(ctx, order.UserID)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("支付宝退款失败")
	}

	now := time.Now()
	if aliRsp.Response.GmtRefundPay != "" {
		if parsed, parseErr := time.ParseInLocation("2006-01-02 15:04:05", aliRsp.Response.GmtRefundPay, time.Local); parseErr == nil {
			now = parsed
		}
	}
	_ = dao.NewRechargeDao(ctx).MarkRefundResult(order.OrderNum, refundNo, model.RechargeRefundSuccess, &now)
	if err = refreshPendingCredit(ctx, order.UserID); err != nil {
		log.LogrusObj.Error(err)
	}

	return &types.RechargeRefundResp{
		OrderNum:     order.OrderNum,
		RefundNo:     refundNo,
		Amount:       req.Amount,
		RefundStatus: model.RechargeRefundSuccess,
		ProviderID:   aliRsp.Response.TradeNo,
	}, nil
}
