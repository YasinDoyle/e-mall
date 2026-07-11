package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var CouponSrvIns *CouponSrv
var CouponSrvOnce sync.Once

type CouponSrv struct{}

func GetCouponSrv() *CouponSrv {
	CouponSrvOnce.Do(func() { CouponSrvIns = &CouponSrv{} })
	return CouponSrvIns
}

// AdminCouponCreate 管理员创建优惠券
func (s *CouponSrv) AdminCouponCreate(ctx context.Context, req *types.AdminCouponCreateReq) (resp interface{}, err error) {
	if req.ExpireAt.Before(time.Now()) {
		return nil, errors.New("过期时间必须晚于当前时间")
	}
	if req.Stock < -1 {
		return nil, errors.New("库存不能小于 -1")
	}
	if req.CouponType == model.CouponTypeFixed && req.Discount <= 0 {
		return nil, errors.New("固定金额优惠必须大于 0")
	}
	if req.CouponType == model.CouponTypePercent && (req.Discount <= 0 || req.Discount > 1) {
		return nil, errors.New("折扣券折扣必须在 0 到 1 之间")
	}

	coupon := &model.Coupon{
		Name:       req.Name,
		CouponType: req.CouponType,
		Discount:   req.Discount,
		MinAmount:  req.MinAmount,
		Stock:      req.Stock,
		ExpireAt:   req.ExpireAt,
	}
	if coupon.Stock == 0 {
		coupon.Stock = -1 // 默认无限
	}
	if err = dao.NewCouponDao(ctx).CreateCoupon(coupon); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "创建成功"
	return
}

func (s *CouponSrv) AdminCouponList(ctx context.Context) (resp interface{}, err error) {
	coupons, err := dao.NewCouponDao(ctx).ListCouponsAdmin()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.CouponResp, 0, len(coupons))
	for _, coupon := range coupons {
		list = append(list, buildCouponResp(coupon))
	}
	resp = &types.DataListResp{Item: list, Total: int64(len(list))}
	return
}

func (s *CouponSrv) AdminCouponOffline(ctx context.Context, req *types.AdminCouponOfflineReq) (resp interface{}, err error) {
	if err = dao.NewCouponDao(ctx).OfflineCoupon(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "下线成功"
	return
}

// CouponList 查询可领取的优惠券列表（公开）
func (s *CouponSrv) CouponList(ctx context.Context) (resp interface{}, err error) {
	coupons, err := dao.NewCouponDao(ctx).ListCoupons()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.CouponResp, 0, len(coupons))
	for _, coupon := range coupons {
		list = append(list, buildCouponResp(coupon))
	}
	resp = &types.DataListResp{Item: list, Total: int64(len(list))}
	return
}

// CouponClaim 用户领券
func (s *CouponSrv) CouponClaim(ctx context.Context, req *types.CouponClaimReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	err = dao.NewCouponDao(ctx).ClaimCouponWithStock(u.Id, req.CouponID)
	if errors.Is(err, dao.ErrCouponAlreadyClaimed) {
		return nil, errors.New("您已领取过该优惠券")
	}
	if errors.Is(err, dao.ErrCouponUnavailable) {
		return nil, errors.New("优惠券不存在、已过期或已发放完毕")
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	resp = "领取成功"
	return
}

// UserCouponList 用户查看自己的优惠券
func (s *CouponSrv) UserCouponList(ctx context.Context) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	userCoupons, err := dao.NewCouponDao(ctx).ListUserCoupons(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	// 批量获取优惠券详情
	result := make([]*types.UserCouponResp, 0, len(userCoupons))
	couponDao := dao.NewCouponDao(ctx)
	for _, uc := range userCoupons {
		c, cErr := couponDao.GetCouponByID(uc.CouponID)
		if cErr != nil {
			continue
		}
		result = append(result, &types.UserCouponResp{
			ID:         uc.ID,
			CouponID:   uc.CouponID,
			Name:       c.Name,
			CouponType: c.CouponType,
			Discount:   c.Discount,
			MinAmount:  c.MinAmount,
			ExpireAt:   c.ExpireAt,
			Status:     uc.Status,
		})
	}
	resp = &types.DataListResp{Item: result, Total: int64(len(result))}
	return
}

// CalcDiscount 计算优惠后金额（供 OrderCreate 调用）
// 返回 finalAmount（优惠后金额）和 userCouponID（用于核销）
func CalcDiscount(ctx context.Context, userID, couponID uint, originalAmount float64) (finalAmount float64, userCouponID uint, err error) {
	return calcDiscount(dao.NewCouponDao(ctx), userID, couponID, originalAmount)
}

func calcDiscount(couponDao *dao.CouponDao, userID, couponID uint, originalAmount float64) (finalAmount float64, userCouponID uint, err error) {
	finalAmount = originalAmount
	if couponID == 0 {
		return
	}
	uc, err := couponDao.GetUserCoupon(userID, couponID)
	if err != nil {
		return 0, 0, errors.New("优惠券不存在或已使用")
	}
	userCouponID = uc.ID

	coupon, err := couponDao.GetCouponByID(couponID)
	if err != nil {
		return 0, 0, errors.New("优惠券信息异常")
	}
	if coupon.ExpireAt.Before(time.Now()) {
		return 0, 0, errors.New("优惠券已过期")
	}
	if originalAmount < coupon.MinAmount {
		return 0, 0, errors.New("未达到优惠券使用门槛")
	}

	switch coupon.CouponType {
	case model.CouponTypeFixed:
		finalAmount = originalAmount - coupon.Discount
	case model.CouponTypePercent:
		finalAmount = originalAmount * coupon.Discount
	}
	if finalAmount < 0 {
		finalAmount = 0
	}
	return
}

func buildCouponResp(coupon *model.Coupon) *types.CouponResp {
	return &types.CouponResp{
		ID:         coupon.ID,
		Name:       coupon.Name,
		CouponType: coupon.CouponType,
		Discount:   coupon.Discount,
		MinAmount:  coupon.MinAmount,
		Stock:      coupon.Stock,
		ExpireAt:   coupon.ExpireAt,
		CreatedAt:  coupon.CreatedAt.Unix(),
	}
}
