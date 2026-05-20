package service

import (
	"context"
	"sync"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	esrepo "github.com/YasinDoyle/e-mall/repository/es"
	"github.com/YasinDoyle/e-mall/types"
)

var ProductIndexSrvIns *ProductIndexSrv
var ProductIndexSrvOnce sync.Once

type ProductIndexSrv struct{}

func GetProductIndexSrv() *ProductIndexSrv {
	ProductIndexSrvOnce.Do(func() {
		ProductIndexSrvIns = &ProductIndexSrv{}
	})
	return ProductIndexSrvIns
}

func (s *ProductIndexSrv) SyncProduct(ctx context.Context, product *model.Product) error {
	return esrepo.NewProductIndexRepo().IndexProduct(ctx, product)
}

func (s *ProductIndexSrv) DeleteProduct(ctx context.Context, productID uint) error {
	return esrepo.NewProductIndexRepo().DeleteProduct(ctx, productID)
}

func (s *ProductIndexSrv) SearchProducts(ctx context.Context, keyword string, page types.BasePage) ([]*types.ProductResp, int64, error) {
	return esrepo.NewProductIndexRepo().SearchProducts(ctx, keyword, page)
}
