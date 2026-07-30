package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	domainevent "github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	esrepo "github.com/YasinDoyle/e-mall/repository/es"
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
	resp = buildProductRespWithCertificates(ctx, p)

	return

}

// 创建商品
func (s *ProductSrv) ProductCreate(ctx context.Context, files []*multipart.FileHeader, certificateFiles []*multipart.FileHeader, req *types.ProductCreateReq) (resp interface{}, err error) {
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
	if len(files) == 0 {
		return nil, errors.New("请上传商品图片")
	}
	// 以第一张作为封面图
	tmp, _ := files[0].Open()
	path, err := util.UploadProductImage(tmp, files[0].Size, uId, req.Name)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	coverPath := path
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

		Brand:             strings.TrimSpace(req.Brand),
		Origin:            strings.TrimSpace(req.Origin),
		Specification:     strings.TrimSpace(req.Specification),
		ProductionDate:    strings.TrimSpace(req.ProductionDate),
		ShelfLife:         strings.TrimSpace(req.ShelfLife),
		ServiceGuarantees: strings.TrimSpace(req.ServiceGuarantees),
		CertificateMeta:   strings.TrimSpace(req.CertificateMeta),
	}
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		num := strconv.Itoa(index)
		if index == 0 {
			path = coverPath
		} else {
			tmp, _ = file.Open()
			path, err = util.UploadProductImage(tmp, file.Size, uId, req.Name+num)
			if err != nil {
				log.LogrusObj.Error(err)
				return
			}
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
	if err = createProductCertificates(productDao.DB, uId, product, certificateFiles, req.CertificateTypes, req.CertificateNames); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	domainevent.Publish(ctx, domainevent.ProductSubmitted{Product: product})
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
		pRespList = append(pRespList, buildProductResp(p))
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
	domainevent.Publish(ctx, domainevent.ProductDeleted{ProductID: req.ID})
	return
}

// 更新商品
func (s *ProductSrv) ProductUpdate(ctx context.Context, certificateFiles []*multipart.FileHeader, req *types.ProductUpdateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	sellerProfile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(u.Id)
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

	productDao := dao.NewProductDao(ctx)
	updates := map[string]interface{}{
		"name":               strings.TrimSpace(req.Name),
		"category_id":        req.CategoryID,
		"title":              strings.TrimSpace(req.Title),
		"info":               strings.TrimSpace(req.Info),
		"price":              req.Price,
		"discount_price":     req.DiscountPrice,
		"num":                req.Num,
		"brand":              strings.TrimSpace(req.Brand),
		"origin":             strings.TrimSpace(req.Origin),
		"specification":      strings.TrimSpace(req.Specification),
		"production_date":    strings.TrimSpace(req.ProductionDate),
		"shelf_life":         strings.TrimSpace(req.ShelfLife),
		"service_guarantees": strings.TrimSpace(req.ServiceGuarantees),
		"certificate_meta":   strings.TrimSpace(req.CertificateMeta),
		"audit_status":       consts.ProductAuditPending,
		"on_sale":            false,
	}
	if err = productDao.UpdateProductByBoss(req.ID, u.Id, updates); err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	product, err := productDao.ShowProductById(req.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if req.ReplaceCertificates {
		certificateDao := dao.NewProductCertificateDaoByDB(productDao.DB)
		if err = certificateDao.DeleteByProductID(product.ID); err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
		if err = createProductCertificates(productDao.DB, u.Id, product, certificateFiles, req.CertificateTypes, req.CertificateNames); err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
	}
	domainevent.Publish(ctx, domainevent.ProductChanged{Product: product})

	return
}

// 搜索商品 TODO 后续用脚本同步数据MySQL到ES，用ES进行搜索
func (s *ProductSrv) ProductSearch(ctx context.Context, req *types.ProductSearchReq) (resp interface{}, err error) {
	req.BasePage = normalizeProductPage(req.BasePage)
	if products, count, searchErr := esrepo.NewProductIndexRepo().SearchProducts(ctx, req.Info, req.BasePage); searchErr == nil {
		for _, p := range products {
			p.BossAvatar = util.AvatarURL(p.BossAvatar)
			p.ImgPath = util.ProductImageURL(p.ImgPath)
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
		pRespList = append(pRespList, buildProductResp(p))
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
		list = append(list, buildProductResp(product))
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

func createProductCertificates(db *gorm.DB, sellerID uint, product *model.Product, files []*multipart.FileHeader, certificateTypes []string, certificateNames []string) error {
	if len(files) == 0 {
		return nil
	}
	certificateDao := dao.NewProductCertificateDaoByDB(db)
	for index, file := range files {
		tmp, err := file.Open()
		if err != nil {
			return err
		}
		path, err := util.UploadProductImage(tmp, file.Size, sellerID, product.Name+"_certificate_"+strconv.Itoa(index))
		if err != nil {
			return err
		}
		certificate := &model.ProductCertificate{
			ProductID:       product.ID,
			CertificateType: valueAt(certificateTypes, index, "other"),
			Name:            valueAt(certificateNames, index, file.Filename),
			FilePath:        path,
		}
		if err = certificateDao.CreateProductCertificate(certificate); err != nil {
			return err
		}
	}
	return nil
}

func buildProductRespWithCertificates(ctx context.Context, product *model.Product) *types.ProductResp {
	resp := buildProductResp(product)
	if product == nil {
		return resp
	}
	certificates, err := dao.NewProductCertificateDao(ctx).ListByProductID(product.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return resp
	}
	resp.Certificates = buildProductCertificateRespList(certificates)
	return resp
}

func buildProductResp(product *model.Product) *types.ProductResp {
	if product == nil {
		return &types.ProductResp{}
	}
	resp := &types.ProductResp{
		ID:                product.ID,
		Name:              product.Name,
		CategoryID:        product.CategoryID,
		Title:             product.Title,
		Info:              product.Info,
		ImgPath:           product.ImgPath,
		Price:             product.Price,
		DiscountPrice:     product.DiscountPrice,
		View:              productView(product),
		CreatedAt:         product.CreatedAt.Unix(),
		Num:               product.Num,
		OnSale:            product.OnSale,
		BossID:            product.BossID,
		BossName:          product.BossName,
		BossAvatar:        product.BossAvatar,
		AuditStatus:       product.AuditStatus,
		Brand:             product.Brand,
		Origin:            product.Origin,
		Specification:     product.Specification,
		ProductionDate:    product.ProductionDate,
		ShelfLife:         product.ShelfLife,
		ServiceGuarantees: product.ServiceGuarantees,
		CertificateMeta:   product.CertificateMeta,
	}
	resp.BossAvatar = avatarURL(resp.BossAvatar)
	resp.ImgPath = productImageURL(resp.ImgPath)
	return resp
}

func buildProductCertificateRespList(certificates []*model.ProductCertificate) []types.ProductCertificateResp {
	list := make([]types.ProductCertificateResp, 0, len(certificates))
	for _, certificate := range certificates {
		list = append(list, buildProductCertificateResp(certificate))
	}
	return list
}

func buildProductCertificateResp(certificate *model.ProductCertificate) types.ProductCertificateResp {
	if certificate == nil {
		return types.ProductCertificateResp{}
	}
	return types.ProductCertificateResp{
		ID:              certificate.ID,
		ProductID:       certificate.ProductID,
		CertificateType: certificate.CertificateType,
		Name:            certificate.Name,
		FilePath:        productImageURL(certificate.FilePath),
		CreatedAt:       certificate.CreatedAt.Unix(),
	}
}

func valueAt(values []string, index int, fallback string) string {
	if index >= len(values) {
		return strings.TrimSpace(fallback)
	}
	value := strings.TrimSpace(values[index])
	if value == "" {
		return strings.TrimSpace(fallback)
	}
	return value
}

func productView(product *model.Product) (view uint64) {
	if product == nil || product.ID == 0 {
		return 0
	}
	defer func() {
		if recover() != nil {
			view = 0
		}
	}()
	return product.View()
}

func productImageURL(path string) (url string) {
	defer func() {
		if recover() != nil {
			url = path
		}
	}()
	return util.ProductImageURL(path)
}

func avatarURL(path string) (url string) {
	defer func() {
		if recover() != nil {
			url = path
		}
	}()
	return util.AvatarURL(path)
}

// ProductImgList 获取商品列表图片
func (s *ProductSrv) ProductImgList(ctx context.Context, req *types.ListProductImgReq) (resp interface{}, err error) {
	productImgs, _ := dao.NewProductImgDao(ctx).ListProductImgByProductId(req.ID)
	for i := range productImgs {
		productImgs[i].ImgPath = productImageURL(productImgs[i].ImgPath)
	}

	resp = &types.DataListResp{
		Item:  productImgs,
		Total: int64(len(productImgs)),
	}

	return
}
