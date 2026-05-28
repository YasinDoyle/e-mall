package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/repository/kafka"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

const flashSaleReserveScript = `
local userKey = KEYS[1]
local stockKey = KEYS[2]

if redis.call("EXISTS", userKey) == 1 then
    return -1
end

local stock = tonumber(redis.call("GET", stockKey) or "0")
if stock <= 0 then
    return -2
end

stock = redis.call("DECR", stockKey)
redis.call("SET", userKey, "1")
return stock
`

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

func flashSaleUserProductKey(userID, productID uint) string {
	return fmt.Sprintf(cache.FlashSaleUserKey, fmt.Sprintf("%d:%d", userID, productID))
}

func (s *FlashSaleSrv) warmFlashSaleList(ctx context.Context, flashSales []*model.FlashSale) error {
	rc := cache.RedisClient
	if err := rc.Del(ctx, cache.FlashSaleListKey).Err(); err != nil {
		return err
	}
	for _, flashSale := range flashSales {
		jsonBytes, err := json.Marshal(flashSale)
		if err != nil {
			return err
		}
		if err := rc.RPush(ctx, cache.FlashSaleListKey, string(jsonBytes)).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (s *FlashSaleSrv) warmFlashSaleDetail(ctx context.Context, flashSale *model.FlashSale) error {
	rc := cache.RedisClient
	jsonBytes, err := json.Marshal(flashSale)
	if err != nil {
		return err
	}
	jsonString := string(jsonBytes)
	if err := rc.Set(ctx, fmt.Sprintf(cache.FlashSaleKey, flashSale.ProductId), jsonString, 0).Err(); err != nil {
		return err
	}
	if err := rc.Set(ctx, fmt.Sprintf(cache.FlashSaleStockKey, flashSale.ProductId), strconv.Itoa(flashSale.Num), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (s *FlashSaleSrv) warmFlashSaleCache(ctx context.Context, flashSales []*model.FlashSale) error {
	if err := s.warmFlashSaleList(ctx, flashSales); err != nil {
		return err
	}

	for _, flashSale := range flashSales {
		if err := s.warmFlashSaleDetail(ctx, flashSale); err != nil {
			return err
		}
	}

	return nil
}

func (s *FlashSaleSrv) loadFlashSaleByProductID(ctx context.Context, productID uint) (*model.FlashSale, error) {
	flashSale, err := dao.NewFlashSaleDao(ctx).GetByProductID(productID)
	if err != nil {
		return nil, err
	}
	if err := s.warmFlashSaleDetail(ctx, flashSale); err != nil {
		return nil, err
	}

	return flashSale, nil
}

func (s *FlashSaleSrv) reserveFlashSaleStock(ctx context.Context, userID, productID uint) (int64, error) {
	rc := cache.RedisClient
	userKey := flashSaleUserProductKey(userID, productID)
	stockKey := fmt.Sprintf(cache.FlashSaleStockKey, productID)
	result, err := rc.Eval(ctx, flashSaleReserveScript, []string{userKey, stockKey}).Int64()
	if err != nil {
		return 0, err
	}

	switch result {
	case -1:
		return 0, errors.New("请勿重复秒杀")
	case -2:
		return 0, errors.New("秒杀商品已售罄")
	default:
		return result, nil
	}
}

func (s *FlashSaleSrv) syncFlashSaleDetailStock(ctx context.Context, productID uint, remainingStock int64) error {
	rc := cache.RedisClient
	key := fmt.Sprintf(cache.FlashSaleKey, productID)
	result, err := rc.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	var flashSale model.FlashSale
	if err := json.Unmarshal([]byte(result), &flashSale); err != nil {
		return err
	}
	flashSale.Num = int(remainingStock)

	return s.warmFlashSaleDetail(ctx, &flashSale)
}

func (s *FlashSaleSrv) rollbackFlashSaleReserve(ctx context.Context, userID, productID uint, remainingStock int64) error {
	rc := cache.RedisClient
	if err := rc.Incr(ctx, fmt.Sprintf(cache.FlashSaleStockKey, productID)).Err(); err != nil {
		return err
	}
	if err := rc.Del(ctx, flashSaleUserProductKey(userID, productID)).Err(); err != nil {
		return err
	}

	return s.syncFlashSaleDetailStock(ctx, productID, remainingStock+1)
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

	if err = s.warmFlashSaleCache(ctx, flashSales); err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
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
		if err = s.warmFlashSaleCache(ctx, flashSales); err != nil {
			log.LogrusObj.Infoln(err)
			return nil, err
		}
		resp = flashSales
	} else {
		flashSales := make([]*model.FlashSale, 0, len(flashSaleList))
		for _, item := range flashSaleList {
			var flashSale model.FlashSale
			if err = json.Unmarshal([]byte(item), &flashSale); err != nil {
				log.LogrusObj.Infoln(err)
				return nil, err
			}
			flashSales = append(flashSales, &flashSale)
		}
		resp = flashSales
	}

	return
}

// GetFlashSale 详情展示
func (s *FlashSaleSrv) GetFlashSale(ctx context.Context, req *types.GetFlashSaleReq) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取列表
	result, err := rc.Get(ctx, fmt.Sprintf(cache.FlashSaleKey, req.ProductId)).Result()
	if err != nil {
		if err != redis.Nil {
			log.LogrusObj.Infoln(err)
			return nil, err
		}
		flashSale, loadErr := s.loadFlashSaleByProductID(ctx, req.ProductId)
		if loadErr != nil {
			log.LogrusObj.Infoln(loadErr)
			return nil, loadErr
		}
		return flashSale, nil
	}
	var flashSale model.FlashSale
	if err = json.Unmarshal([]byte(result), &flashSale); err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
	}

	return &flashSale, nil
}

// FlashSale 秒杀商品
func (s *FlashSaleSrv) FlashSale(ctx context.Context, req *types.FlashSaleReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
	}

	flashSaleResp, err := s.GetFlashSale(ctx, &types.GetFlashSaleReq{ProductId: req.ProductId})
	if err != nil {
		return nil, err
	}
	flashSale, ok := flashSaleResp.(*model.FlashSale)
	if !ok {
		return nil, errors.New("秒杀商品信息错误")
	}

	remainingStock, err := s.reserveFlashSaleStock(ctx, u.Id, req.ProductId)
	if err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
	}
	if err = s.syncFlashSaleDetailStock(ctx, req.ProductId, remainingStock); err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
	}

	message := &model.FlashSale2MQ{
		FlashSaleId: flashSale.Id,
		ProductId:   flashSale.ProductId,
		BossId:      flashSale.BossId,
		UserId:      u.Id,
		Money:       flashSale.Money,
		AddressId:   req.AddressId,
		Key:         req.Key,
	}
	if err = kafka.PublishFlashSaleOrder(ctx, message); err != nil {
		log.LogrusObj.Infoln(err)
		rollbackErr := s.rollbackFlashSaleReserve(ctx, u.Id, req.ProductId, remainingStock)
		if rollbackErr != nil {
			log.LogrusObj.Infoln(rollbackErr)
		}
		return nil, err
	}

	return &types.FlashSaleResp{
		ProductId:      req.ProductId,
		UserId:         u.Id,
		RemainingStock: remainingStock,
		Status:         "秒杀成功，等待下单",
	}, nil
}
