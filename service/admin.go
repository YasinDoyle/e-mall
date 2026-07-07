package service

import (
	"context"
	"sync"

	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var AdminSrvIns *AdminSrv
var AdminSrvOnce sync.Once

type AdminSrv struct{}

func GetAdminSrv() *AdminSrv {
	AdminSrvOnce.Do(func() { AdminSrvIns = &AdminSrv{} })
	return AdminSrvIns
}

// ===== 分类 =====

func (s *AdminSrv) CategoryCreate(ctx context.Context, req *types.AdminCategoryReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).CreateCategory(req.CategoryName); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "创建成功"
	return
}

func (s *AdminSrv) CategoryUpdate(ctx context.Context, req *types.AdminCategoryUpdateReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).UpdateCategory(req.ID, req.CategoryName); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "更新成功"
	return
}

func (s *AdminSrv) CategoryDelete(ctx context.Context, req *types.AdminIDReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).DeleteCategory(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "删除成功"
	return
}

// ===== 轮播图 =====

func (s *AdminSrv) CarouselCreate(ctx context.Context, req *types.AdminCarouselReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).CreateCarousel(req.ImgPath, req.ProductID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "创建成功"
	return
}

func (s *AdminSrv) CarouselDelete(ctx context.Context, req *types.AdminIDReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).DeleteCarousel(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "删除成功"
	return
}

// ===== 公告 =====

func (s *AdminSrv) NoticeList(ctx context.Context) (resp interface{}, err error) {
	notices, err := dao.NewAdminDao(ctx).ListNotice()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.DataListResp{Item: notices, Total: int64(len(notices))}
	return
}

func (s *AdminSrv) NoticeCreate(ctx context.Context, req *types.AdminNoticeReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).CreateNotice(req.Text); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "创建成功"
	return
}

func (s *AdminSrv) NoticeUpdate(ctx context.Context, req *types.AdminNoticeUpdateReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).UpdateNotice(req.ID, req.Text); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "更新成功"
	return
}

func (s *AdminSrv) NoticeDelete(ctx context.Context, req *types.AdminIDReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).DeleteNotice(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "删除成功"
	return
}

// ===== 用户管理 =====

func (s *AdminSrv) UserList(ctx context.Context, req *types.AdminListReq) (resp interface{}, err error) {
	if req.PageSize == 0 {
		req.PageSize = 15
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	users, total, err := dao.NewAdminDao(ctx).ListUsers(req.PageNum, req.PageSize)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.DataListResp{Item: users, Total: total}
	return
}

func (s *AdminSrv) UserBan(ctx context.Context, req *types.AdminUserBanReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).SetUserBan(req.ID, req.Banned); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "操作成功"
	return
}

// ===== 商品审核 =====

func (s *AdminSrv) ProductList(ctx context.Context, req *types.AdminProductListReq) (resp interface{}, err error) {
	if req.PageSize == 0 {
		req.PageSize = 15
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	products, total, err := dao.NewAdminDao(ctx).ListProductsAdmin(req.PageNum, req.PageSize, req.AuditStatus)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.DataListResp{Item: products, Total: total}
	return
}

func (s *AdminSrv) ProductAudit(ctx context.Context, req *types.AdminProductAuditReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).AuditProduct(req.ID, req.AuditStatus); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "审核完成"
	return
}
