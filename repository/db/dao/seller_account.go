package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type SellerAccountDao struct {
	*gorm.DB
}

func NewSellerAccountDao(ctx context.Context) *SellerAccountDao {
	return &SellerAccountDao{NewDBClient(ctx)}
}

func NewSellerAccountDaoByDB(db *gorm.DB) *SellerAccountDao {
	return &SellerAccountDao{db}
}

func (dao *SellerAccountDao) GetBySellerID(sellerID uint) (*model.SellerAccount, error) {
	var account model.SellerAccount
	err := dao.DB.Where("seller_id = ?", sellerID).First(&account).Error
	return &account, err
}

func (dao *SellerAccountDao) GetOrCreateBySellerID(sellerID uint) (*model.SellerAccount, error) {
	account := &model.SellerAccount{SellerID: sellerID}
	err := dao.DB.Where("seller_id = ?", sellerID).
		Attrs(&model.SellerAccount{
			AvailableBalance: 0,
			FrozenBalance:    0,
			TotalIncome:      0,
			TotalWithdrawn:   0,
		}).
		FirstOrCreate(account).Error
	return account, err
}

func (dao *SellerAccountDao) Save(account *model.SellerAccount) error {
	return dao.DB.Save(account).Error
}
