package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
	util "github.com/YasinDoyle/e-mall/utils/upload"
	"gorm.io/gorm"
)

var AdminSrvIns *AdminSrv
var AdminSrvOnce sync.Once

type AdminSrv struct{}

const adminStatsDateLayout = "2006-01-02"

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
	if err = validateCarouselProduct(ctx, req.ProductID); err != nil {
		return
	}
	if err = dao.NewAdminDao(ctx).CreateCarousel(req.ImgPath, req.ProductID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "创建成功"
	return
}

func (s *AdminSrv) CarouselUpload(ctx context.Context, file multipart.File, fileSize int64) (resp interface{}, err error) {
	fileName := fmt.Sprintf("carousel_%d", time.Now().UnixNano())
	path, err := util.UploadCarouselImage(file, fileSize, fileName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp = &types.AdminUploadResp{URL: util.ProductImageURL(path)}
	return
}

func (s *AdminSrv) CarouselList(ctx context.Context) (resp interface{}, err error) {
	return GetCarouselSrv().ListCarousel(ctx, &types.ListCarouselReq{})
}

func (s *AdminSrv) CarouselUpdate(ctx context.Context, req *types.AdminCarouselUpdateReq) (resp interface{}, err error) {
	if err = validateCarouselProduct(ctx, req.ProductID); err != nil {
		return
	}
	if err = dao.NewAdminDao(ctx).UpdateCarousel(req.ID, req.ImgPath, req.ProductID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "更新成功"
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

func validateCarouselProduct(ctx context.Context, productID uint) error {
	if productID == 0 {
		return e.NewBusinessError(e.ErrorCarouselProductRequired)
	}
	if _, err := dao.NewProductDao(ctx).GetProductById(productID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewBusinessError(e.ErrorCarouselProductNotExist)
		}
		log.LogrusObj.Error(err)
		return err
	}
	return nil
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
	list := make([]*types.AdminUserResp, 0, len(users))
	for _, user := range users {
		list = append(list, &types.AdminUserResp{
			ID:        user.ID,
			UserName:  user.UserName,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    user.Status,
			Avatar:    user.Avatar,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Unix(),
		})
	}
	resp = &types.DataListResp{Item: list, Total: total}
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

// ===== 商家管理 =====

func (s *AdminSrv) SellerList(ctx context.Context, req *types.AdminSellerListReq) (resp interface{}, err error) {
	if req.PageSize == 0 {
		req.PageSize = 15
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	profiles, total, err := dao.NewSellerDao(ctx).ListSellerProfiles(req.PageNum, req.PageSize, req.Status)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.AdminSellerResp, 0, len(profiles))
	for _, profile := range profiles {
		list = append(list, buildAdminSellerResp(profile))
	}
	resp = &types.DataListResp{Item: list, Total: total}
	return
}

func (s *AdminSrv) SellerAudit(ctx context.Context, req *types.AdminSellerAuditReq) (resp interface{}, err error) {
	if err = validateAdminSellerAuditReq(req); err != nil {
		return nil, err
	}
	if err = dao.NewSellerDao(ctx).AuditSellerProfile(req.ID, req.Status, strings.TrimSpace(req.RejectReason)); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "审核完成"
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
	if req.AuditStatus == nil && req.Status != nil {
		req.AuditStatus = req.Status
	}
	products, total, err := dao.NewAdminDao(ctx).ListProductsAdmin(req.PageNum, req.PageSize, req.AuditStatus)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.AdminProductResp, 0, len(products))
	for _, product := range products {
		item := &types.AdminProductResp{
			ID:            product.ID,
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			OnSale:        product.OnSale,
			Num:           product.Num,
			BossID:        product.BossID,
			BossName:      product.BossName,
			BossAvatar:    product.BossAvatar,
			AuditStatus:   product.AuditStatus,
			Status:        product.AuditStatus,
			CreatedAt:     product.CreatedAt.Unix(),
		}
		item.ImgPath = util.ProductImageURL(item.ImgPath)
		list = append(list, item)
	}
	resp = &types.DataListResp{Item: list, Total: total}
	return
}

func (s *AdminSrv) ProductAudit(ctx context.Context, req *types.AdminProductAuditReq) (resp interface{}, err error) {
	if req.AuditStatus == 0 && req.Status != 0 {
		req.AuditStatus = req.Status
	}
	if req.AuditStatus == consts.ProductAuditApproved {
		product, loadErr := dao.NewProductDao(ctx).ShowProductById(req.ID)
		if loadErr != nil {
			log.LogrusObj.Error(loadErr)
			return nil, loadErr
		}
		boss, loadErr := dao.NewUserDao(ctx).GetUserById(product.BossID)
		if loadErr != nil {
			log.LogrusObj.Error(loadErr)
			return nil, loadErr
		}
		sellerProfile, loadErr := dao.NewSellerDao(ctx).GetSellerProfileByUserID(product.BossID)
		if loadErr != nil {
			if errors.Is(loadErr, gorm.ErrRecordNotFound) {
				return nil, e.NewBusinessError(e.ErrorProductSellerNotApproved)
			}
			log.LogrusObj.Error(loadErr)
			return nil, loadErr
		}
		if err = ensureSellerProfileApproved(sellerProfile); err != nil {
			return
		}
		if err = ensureSellerCanEnableTrading(boss, true); err != nil {
			return
		}
	}
	if err = dao.NewAdminDao(ctx).AuditProduct(req.ID, req.AuditStatus); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "审核完成"
	return
}

func (s *AdminSrv) ProductDelete(ctx context.Context, req *types.AdminIDReq) (resp interface{}, err error) {
	if err = dao.NewAdminDao(ctx).DeleteProduct(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if syncErr := GetProductIndexSrv().DeleteProduct(ctx, req.ID); syncErr != nil {
		log.LogrusObj.Errorln(syncErr)
	}
	resp = "删除成功"
	return
}

// ===== 统计 =====

func (s *AdminSrv) StatsOverview(ctx context.Context) (resp interface{}, err error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrowStart := todayStart.AddDate(0, 0, 1)

	adminDao := dao.NewAdminDao(ctx)
	todayOrders, err := adminDao.CountTodayOrders(todayStart, tomorrowStart)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	totalSales, err := adminDao.SumTotalSales()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	registeredUsers, err := adminDao.CountRegisteredUsers()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	resp = &types.AdminStatsOverviewResp{
		TodayOrders:     todayOrders,
		TotalSales:      totalSales,
		RegisteredUsers: registeredUsers,
	}
	return
}

func (s *AdminSrv) StatsOrders(ctx context.Context, req *types.AdminStatsOrdersReq) (resp interface{}, err error) {
	start, endExclusive, err := parseAdminStatsRange(req)
	if err != nil {
		return nil, err
	}

	rows, err := dao.NewAdminDao(ctx).OrderTrend(start, endExclusive)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	countByDate := make(map[string]int64, len(rows))
	salesByDate := make(map[string]float64, len(rows))
	for _, row := range rows {
		dateKey := normalizeAdminStatsDateKey(row.Date)
		countByDate[dateKey] = row.OrderCount
		salesByDate[dateKey] = row.SalesAmount
	}

	dates := make([]string, 0)
	orderCounts := make([]int64, 0)
	salesAmounts := make([]float64, 0)
	for day := start; day.Before(endExclusive); day = day.AddDate(0, 0, 1) {
		date := day.Format(adminStatsDateLayout)
		dates = append(dates, date)
		orderCounts = append(orderCounts, countByDate[date])
		salesAmounts = append(salesAmounts, salesByDate[date])
	}

	resp = &types.AdminStatsOrdersResp{
		Dates:        dates,
		OrderCounts:  orderCounts,
		SalesAmounts: salesAmounts,
	}
	return
}

func parseAdminStatsRange(req *types.AdminStatsOrdersReq) (time.Time, time.Time, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := todayStart.AddDate(0, 0, -6)
	end := todayStart
	var err error

	if req.StartDate != "" {
		start, err = time.ParseInLocation(adminStatsDateLayout, req.StartDate, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("start_date格式应为YYYY-MM-DD")
		}
	}
	if req.EndDate != "" {
		end, err = time.ParseInLocation(adminStatsDateLayout, req.EndDate, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("end_date格式应为YYYY-MM-DD")
		}
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errors.New("end_date不能早于start_date")
	}

	return start, end.AddDate(0, 0, 1), nil
}

func normalizeAdminStatsDateKey(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if len(value) == len(adminStatsDateLayout) {
		return value
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.Format(adminStatsDateLayout)
	}
	if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local); err == nil {
		return parsed.Format(adminStatsDateLayout)
	}
	if len(value) >= len(adminStatsDateLayout) {
		return value[:len(adminStatsDateLayout)]
	}
	return value
}
