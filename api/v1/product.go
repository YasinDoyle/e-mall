package v1

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/service"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

// CreateProductHandler 创建商品
func CreateProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		form, _ := ctx.MultipartForm()
		files := []*multipart.FileHeader{}
		certificateFiles := []*multipart.FileHeader{}
		if form != nil {
			files = form.File["image"]
			certificateFiles = form.File["certificate"]
		}
		req.CertificateTypes = ctx.PostFormArray("certificate_type")
		req.CertificateNames = ctx.PostFormArray("certificate_name")
		l := service.GetProductSrv()
		resp, err := l.ProductCreate(ctx.Request.Context(), files, certificateFiles, &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ListProductsHandler 商品列表
func ListProductsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductListReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ShowProductHandler 商品详情
func ShowProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductShowReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductShow(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// DeleteProductHandler 删除商品
func DeleteProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductDeleteReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductDelete(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// UpdateProductHandler 更新商品
func UpdateProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		form, _ := ctx.MultipartForm()
		certificateFiles := []*multipart.FileHeader{}
		if form != nil {
			certificateFiles = form.File["certificate"]
		}
		req.ReplaceCertificates = ctx.PostForm("replace_certificates") == "true"
		req.CertificateTypes = ctx.PostFormArray("certificate_type")
		req.CertificateNames = ctx.PostFormArray("certificate_name")

		l := service.GetProductSrv()
		resp, err := l.ProductUpdate(ctx.Request.Context(), certificateFiles, &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// SearchProductsHandler 搜索商品
func SearchProductsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductSearchReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}
		if req.PageNum == 0 {
			req.PageNum = 1
		}

		l := service.GetProductSrv()
		resp, err := l.ProductSearch(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

func ListProductImgHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ListProductImgReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.ID == 0 {
			err := errors.New("参数错误,id不能为空")
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductImgList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// BossProductListHandler 卖家查看自己的商品列表
func BossProductListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.BossProductListReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		resp, err := service.GetProductSrv().BossProductList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// BossProductOnSaleHandler 卖家上架/下架商品
func BossProductOnSaleHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.BossProductOnSaleReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		resp, err := service.GetProductSrv().BossProductOnSale(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
