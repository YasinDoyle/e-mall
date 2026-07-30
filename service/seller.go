package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	domainevent "github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var SellerSrvIns *SellerSrv
var SellerSrvOnce sync.Once

type SellerSrv struct{}

func GetSellerSrv() *SellerSrv {
	SellerSrvOnce.Do(func() {
		SellerSrvIns = &SellerSrv{}
	})
	return SellerSrvIns
}

func (s *SellerSrv) Apply(ctx context.Context, req *types.SellerApplyReq) (resp interface{}, err error) {
	if err = validateSellerApplyReq(req); err != nil {
		return nil, err
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	shopName := strings.TrimSpace(req.ShopName)
	description := strings.TrimSpace(req.Description)
	sellerDao := dao.NewSellerDao(ctx)
	profile, err := sellerDao.GetSellerProfileByUserID(u.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = &model.SellerProfile{
			UserID:      u.Id,
			ShopName:    shopName,
			Description: description,
			Status:      consts.SellerStatusPending,
		}
		if err = sellerDao.CreateSellerProfile(profile); err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
		notifySellerApplicationSubmitted(ctx, profile)
		resp = buildSellerProfileResp(profile)
		return
	}

	switch profile.Status {
	case consts.SellerStatusPending:
		return nil, e.NewBusinessError(e.ErrorSellerAuditPending)
	case consts.SellerStatusApproved:
		return nil, e.NewBusinessError(e.ErrorSellerAlreadyApproved)
	case consts.SellerStatusBanned:
		return nil, e.NewBusinessError(e.ErrorSellerBanned)
	case consts.SellerStatusRejected:
		if err = sellerDao.UpdateSellerApplication(u.Id, shopName, description); err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
		profile.ShopName = shopName
		profile.Description = description
		profile.Status = consts.SellerStatusPending
		profile.RejectReason = ""
		profile.ApprovedAt = nil
		notifySellerApplicationSubmitted(ctx, profile)
		resp = buildSellerProfileResp(profile)
		return
	default:
		return nil, e.NewBusinessError(e.ErrorSellerInvalidStatus)
	}
}

func notifySellerApplicationSubmitted(ctx context.Context, profile *model.SellerProfile) {
	if profile == nil {
		return
	}
	domainevent.Publish(ctx, domainevent.SellerApplied{SellerID: profile.UserID, ShopName: profile.ShopName})
}

func (s *SellerSrv) Profile(ctx context.Context) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	profile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(u.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSellerNotApplied)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp = buildSellerProfileResp(profile)
	return
}

func validateSellerApplyReq(req *types.SellerApplyReq) error {
	if req == nil {
		return e.NewBusinessError(e.ErrorSellerInvalidApplication)
	}
	if strings.TrimSpace(req.ShopName) == "" {
		return e.NewBusinessError(e.ErrorSellerShopNameRequired)
	}
	if len([]rune(strings.TrimSpace(req.ShopName))) > 80 {
		return e.NewBusinessError(e.ErrorSellerShopNameTooLong)
	}
	if len([]rune(strings.TrimSpace(req.Description))) > 500 {
		return e.NewBusinessError(e.ErrorSellerDescriptionTooLong)
	}
	return nil
}

func validateAdminSellerAuditReq(req *types.AdminSellerAuditReq) error {
	if req == nil || req.ID == 0 {
		return e.NewBusinessError(e.ErrorSellerInvalidApplication)
	}
	switch req.Status {
	case consts.SellerStatusApproved, consts.SellerStatusRejected, consts.SellerStatusBanned:
	default:
		return e.NewBusinessError(e.ErrorSellerAuditStatusInvalid)
	}
	if req.Status == consts.SellerStatusRejected && strings.TrimSpace(req.RejectReason) == "" {
		return e.NewBusinessError(e.ErrorSellerRejectReasonMissing)
	}
	return nil
}

func buildSellerProfileResp(profile *model.SellerProfile) *types.SellerProfileResp {
	if profile == nil {
		return nil
	}
	var approvedAt int64
	if profile.ApprovedAt != nil {
		approvedAt = profile.ApprovedAt.Unix()
	}
	return &types.SellerProfileResp{
		ID:           profile.ID,
		UserID:       profile.UserID,
		UserName:     profile.User.UserName,
		ShopName:     profile.ShopName,
		Description:  profile.Description,
		Status:       profile.Status,
		StatusText:   consts.SellerStatusMap[profile.Status],
		RejectReason: profile.RejectReason,
		ApprovedAt:   approvedAt,
		CreatedAt:    profile.CreatedAt.Unix(),
	}
}

func buildAdminSellerResp(profile *model.SellerProfile) *types.AdminSellerResp {
	if profile == nil {
		return nil
	}
	var approvedAt int64
	if profile.ApprovedAt != nil {
		approvedAt = profile.ApprovedAt.Unix()
	}
	return &types.AdminSellerResp{
		ID:           profile.ID,
		UserID:       profile.UserID,
		UserName:     profile.User.UserName,
		NickName:     profile.User.NickName,
		Email:        profile.User.Email,
		ShopName:     profile.ShopName,
		Description:  profile.Description,
		Status:       profile.Status,
		StatusText:   consts.SellerStatusMap[profile.Status],
		RejectReason: profile.RejectReason,
		ApprovedAt:   approvedAt,
		CreatedAt:    profile.CreatedAt.Unix(),
	}
}
