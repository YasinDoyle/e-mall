package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var FlashSaleSrvIns *FlashSaleSrv
var FlashSaleSrvOnce sync.Once

type FlashSaleSrv struct {
}

func GetFlashSaleSrv() *FlashSaleSrv {
	FlashSaleSrvOnce.Do(func() {
		FlashSaleSrvIns = &FlashSaleSrv{}
	})
	return FlashSaleSrvIns
}

// InitFlashSale 初始化秒杀商品信息
func (s *FlashSaleSrv) InitFlashSale(ctx context.Context) (resp interface{}, err error) {
	flashSales := make([]*model.FlashSale, 0)
	for i := 1; i < 10; i++ {
		flashSales = append(flashSales, &model.FlashSale{
			ProductId: uint(i),
			BossId:    2,
			Title:     "秒杀商品测试使用",
			Money:     200,
			Num:       10,
		})
	}
	err = dao.NewFlashSaleDao(ctx).BatchCreate(flashSales)
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	// 导入数据库的同时，初始化缓存
	for i := range flashSales {
		jsonBytes, errx := json.Marshal(flashSales[i])
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return
		}
		jsonString := string(jsonBytes)
		_, errx = cache.RedisClient.LPush(ctx, cache.FlashSaleListKey, jsonString).Result()
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return nil, errx
		}
	}

	return
}

// ListFlashSales 列表展示
func (s *FlashSaleSrv) ListFlashSales(ctx context.Context) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取列表
	flashSaleList, err := rc.LRange(ctx, cache.FlashSaleListKey, 0, -1).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	if len(flashSaleList) == 0 {
		flashSales, errx := dao.NewFlashSaleDao(ctx).ListFlashSales()
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return nil, errx
		}

		for i := range flashSales {
			// 将结构体转换为JSON格式的字符串
			jsonBytes, errx := json.Marshal(flashSales[i])
			if errx != nil {
				log.LogrusObj.Infoln(errx)
				return
			}
			// 将字节数组转换为字符串
			jsonString := string(jsonBytes)
			_, errx = rc.LPush(ctx, cache.FlashSaleListKey, jsonString).Result()
			if errx != nil {
				log.LogrusObj.Infoln(errx)
				return nil, errx
			}
		}
		resp = flashSales
	} else {
		resp = flashSaleList
	}

	return
}

// GetFlashSale 详情展示
func (s *FlashSaleSrv) GetFlashSale(ctx context.Context, req *types.GetFlashSaleReq) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取列表
	resp, err = rc.Get(ctx,
		fmt.Sprintf(cache.FlashSaleKey, req.ProductId)).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	return
}

// FlashSale 秒杀商品
func (s *FlashSaleSrv) FlashSale(ctx context.Context, req *types.FlashSaleReq) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取数据
	resp, err = rc.Get(ctx,
		fmt.Sprintf(cache.FlashSaleKey, req.ProductId)).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	return
}
