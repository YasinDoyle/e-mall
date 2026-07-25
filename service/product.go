package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"
	"sync"

	"gorm.io/gorm"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
	util "github.com/YasinDoyle/e-mall/utils/upload"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once

type ProductSrv struct {
}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{}
	})
	return ProductSrvIns
}

func (s *ProductSrv) ProductShow(ctx context.Context, req *types.ProductShowReq) (resp interface{}, err error) {
	p, err := dao.NewProductDao(ctx).ShowProductById(req.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	pResp := &types.ProductResp{
		ID:            p.ID,
		Name:          p.Name,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Info:          p.Info,
		ImgPath:       p.ImgPath,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		View:          p.View(),
		CreatedAt:     p.CreatedAt.Unix(),
		Num:           p.Num,
		OnSale:        p.OnSale,
		BossID:        p.BossID,
		BossName:      p.BossName,
		BossAvatar:    p.BossAvatar,
	}
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
		pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
	}

	resp = pResp

	return

}

// 创建商品
func (s *ProductSrv) ProductCreate(ctx context.Context, files []*multipart.FileHeader, req *types.ProductCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	uId := u.Id
	boss, _ := dao.NewUserDao(ctx).GetUserById(uId)
	sellerProfile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(uId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSellerNotApproved)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	if err = ensureSellerProfileApproved(sellerProfile); err != nil {
		return nil, err
	}
	// 以第一张作为封面图
	tmp, _ := files[0].Open()
	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name)
	} else {
		path, err = util.UploadToQiNiu(tmp, files[0].Size)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       path,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        false, // 发布后待管理员审核，审核通过后才上架
		AuditStatus:   consts.ProductAuditPending,
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if syncErr := GetProductIndexSrv().SyncProduct(ctx, product); syncErr != nil {
		log.LogrusObj.Errorln(syncErr)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		num := strconv.Itoa(index)
		tmp, _ = file.Open()
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name+num)
		} else {
			path, err = util.UploadToQiNiu(tmp, file.Size)
		}
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		err = dao.NewProductImgDaoByDB(productDao.DB).CreateProductImg(productImg)
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		wg.Done()
	}

	wg.Wait()

	return
}

func (s *ProductSrv) ProductList(ctx context.Context, req *types.ProductListReq) (resp interface{}, err error) {
	var total int64
	condition := make(map[string]interface{})
	if req.CategoryID != 0 {
		condition["category_id"] = req.CategoryID
	}
	productDao := dao.NewProductDao(ctx)
	products, _ := productDao.ListProductByCondition(condition, req.BasePage)
	total, err = productDao.CountProductByCondition(condition)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
			pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp)
	}

	resp = &types.DataListResp{
		Item:  pRespList,
		Total: total,
	}

	return
}

// ProductDelete 删除商品
func (s *ProductSrv) ProductDelete(ctx context.Context, req *types.ProductDeleteReq) (resp interface{}, err error) {
	u, _ := ctl.GetUserInfo(ctx)
	err = dao.NewProductDao(ctx).DeleteProduct(req.ID, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if syncErr := GetProductIndexSrv().DeleteProduct(ctx, req.ID); syncErr != nil {
		log.LogrusObj.Errorln(syncErr)
	}
	return
}

// 更新商品
func (s *ProductSrv) ProductUpdate(ctx context.Context, req *types.ProductUpdateReq) (resp interface{}, err error) {
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Info:       req.Info,
		// ImgPath:       service.ImgPath,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        req.OnSale,
	}
	err = dao.NewProductDao(ctx).UpdateProduct(req.ID, product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	product, err = dao.NewProductDao(ctx).ShowProductById(req.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if syncErr := GetProductIndexSrv().SyncProduct(ctx, product); syncErr != nil {
		log.LogrusObj.Errorln(syncErr)
	}

	return
}

// 搜索商品 TODO 后续用脚本同步数据MySQL到ES，用ES进行搜索
func (s *ProductSrv) ProductSearch(ctx context.Context, req *types.ProductSearchReq) (resp interface{}, err error) {
	req.BasePage = normalizeProductPage(req.BasePage)
	if products, count, searchErr := GetProductIndexSrv().SearchProducts(ctx, req.Info, req.BasePage); searchErr == nil {
		for _, p := range products {
			if conf.Config.System.UploadModel == consts.UploadModelLocal {
				p.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + p.BossAvatar
				p.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + p.ImgPath
			}
		}
		resp = &types.DataListResp{
			Item:  products,
			Total: count,
		}
		return
	} else {
		log.LogrusObj.Errorln(searchErr)
	}

	products, count, err := dao.NewProductDao(ctx).SearchProduct(req.Info, req.BasePage)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
			pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp)
	}

	resp = &types.DataListResp{
		Item:  pRespList,
		Total: count,
	}

	return
}

func normalizeProductPage(page types.BasePage) types.BasePage {
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = consts.BasePageSize
	}
	return page
}

// ===== 卖家中心 =====

// BossProductList 卖家查看自己发布的商品列表
func (s *ProductSrv) BossProductList(ctx context.Context, req *types.BossProductListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	req.BasePage = normalizeProductPage(req.BasePage)
	products, total, err := dao.NewProductDao(ctx).ListProductByBoss(u.Id, req.BasePage)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.ProductResp, 0, len(products))
	for _, product := range products {
		item := &types.ProductResp{
			ID:            product.ID,
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			View:          product.View(),
			CreatedAt:     product.CreatedAt.Unix(),
			Num:           product.Num,
			OnSale:        product.OnSale,
			BossID:        product.BossID,
			BossName:      product.BossName,
			BossAvatar:    product.BossAvatar,
			AuditStatus:   product.AuditStatus,
		}
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			item.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + product.BossAvatar
			item.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + product.ImgPath
		}
		list = append(list, item)
	}
	resp = &types.DataListResp{Item: list, Total: total}
	return
}

// BossProductOnSale 卖家上架/下架自己的商品（仅审核通过的商品可上架）
func (s *ProductSrv) BossProductOnSale(ctx context.Context, req *types.BossProductOnSaleReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	user, err := dao.NewUserDao(ctx).GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if req.OnSale {
		sellerProfile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(u.Id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, e.NewBusinessError(e.ErrorSellerNotApproved)
			}
			log.LogrusObj.Error(err)
			return nil, err
		}
		if err = ensureSellerCanChangeSaleStatus(sellerProfile, req.OnSale); err != nil {
			return nil, err
		}
	}
	if err = ensureSellerCanEnableTrading(user, req.OnSale); err != nil {
		return
	}
	err = dao.NewProductDao(ctx).SetProductOnSale(req.ID, u.Id, req.OnSale)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "操作成功"
	return
}

func ensureSellerCanChangeSaleStatus(profile *model.SellerProfile, onSale bool) error {
	if !onSale {
		return nil
	}
	return ensureSellerProfileApproved(profile)
}

func ensureSellerProfileApproved(profile *model.SellerProfile) error {
	if profile == nil || !profile.IsApproved() {
		return e.NewBusinessError(e.ErrorSellerNotApproved)
	}
	return nil
}

func ensureSellerCanEnableTrading(user *model.User, onSale bool) error {
	if !onSale {
		return nil
	}
	if user == nil || !user.HasPayKey() {
		return e.NewBusinessError(e.ErrorSellerPayKeyRequired)
	}
	return nil
}

// ProductImgList 获取商品列表图片
func (s *ProductSrv) ProductImgList(ctx context.Context, req *types.ListProductImgReq) (resp interface{}, err error) {
	productImgs, _ := dao.NewProductImgDao(ctx).ListProductImgByProductId(req.ID)
	for i := range productImgs {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			productImgs[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + productImgs[i].ImgPath
		}
	}

	resp = &types.DataListResp{
		Item:  productImgs,
		Total: int64(len(productImgs)),
	}

	return
}
