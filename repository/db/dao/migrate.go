package dao

import (
	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/utils/log"
)

func migrate() (err error) {
	err = _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Favorite{},
			&model.Order{}, &model.Admin{}, &model.Address{},
			&model.Cart{}, &model.Category{}, &model.Carousel{},
			&model.Notice{}, &model.Product{},
			&model.ProductImg{}, &model.ProductCertificate{}, &model.FlashSale{},
			&model.FlashSale2MQ{}, &model.Review{},
			&model.Coupon{}, &model.UserCoupon{},
			&model.RechargeOrder{}, &model.SellerProfile{},
			&model.CommissionConfig{}, &model.AccountFlow{},
			&model.Settlement{}, &model.SellerAccount{},
			&model.SellerWithdraw{}, &model.Notification{},
		)
	if err != nil {
		return
	}
	seedAdmin()
	return
}

// seedAdmin 首次启动时写入超级管理员账号（已存在则跳过）
func seedAdmin() {
	adminConf := conf.Config.Admin
	if adminConf == nil || adminConf.UserName == "" {
		return
	}

	var count int64
	_db.Model(&model.User{}).Where("user_name = ?", adminConf.UserName).Count(&count)
	if count > 0 {
		return // 已存在，跳过
	}

	admin := &model.User{
		UserName: adminConf.UserName,
		NickName: adminConf.NickName,
		Status:   model.Active,
		IsAdmin:  true,
	}
	if err := admin.SetPassword(adminConf.Password); err != nil {
		log.LogrusObj.Errorf("seedAdmin SetPassword error: %v", err)
		return
	}
	if err := _db.Create(admin).Error; err != nil {
		log.LogrusObj.Errorf("seedAdmin create error: %v", err)
		return
	}
	log.LogrusObj.Infof("seedAdmin: 管理员账号 [%s] 创建成功", adminConf.UserName)
}
