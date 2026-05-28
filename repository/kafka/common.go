package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	segmentio "github.com/segmentio/kafka-go"
	"gorm.io/gorm"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	log "github.com/YasinDoyle/e-mall/utils/log"
)

var (
	writer       *segmentio.Writer
	consumerOnce sync.Once
)

func InitKafka() {
	config := defaultConfig()
	if config == nil {
		return
	}

	if writer == nil {
		writer = &segmentio.Writer{
			Addr:     segmentio.TCP(strings.Split(config.Address, ",")...),
			Topic:    consts.FlashSaleQueues,
			Balancer: &segmentio.LeastBytes{},
		}
	}

	if config.DisableConsumer {
		return
	}

	consumerOnce.Do(func() {
		go consumeFlashSaleOrders(config)
	})
}

func PublishFlashSaleOrder(ctx context.Context, payload *model.FlashSale2MQ) error {
	InitKafka()
	if writer == nil {
		return fmt.Errorf("kafka unavailable")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return writer.WriteMessages(ctx, segmentio.Message{
		Key:   []byte(fmt.Sprintf("%d:%d", payload.FlashSaleId, payload.UserId)),
		Value: body,
	})
}

func defaultConfig() *conf.KafkaConfig {
	if conf.Config == nil || conf.Config.KafKa == nil {
		return nil
	}

	return conf.Config.KafKa["default"]
}

func consumeFlashSaleOrders(config *conf.KafkaConfig) {
	reader := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:  strings.Split(config.Address, ","),
		Topic:    consts.FlashSaleQueues,
		GroupID:  "e-mall-flash-sale",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			if log.LogrusObj != nil {
				log.LogrusObj.Error(err)
			}
			return
		}

		var payload model.FlashSale2MQ
		if err = json.Unmarshal(msg.Value, &payload); err != nil {
			if log.LogrusObj != nil {
				log.LogrusObj.Error(err)
			}
			continue
		}

		if err = handleFlashSaleOrder(&payload); err != nil {
			if log.LogrusObj != nil {
				log.LogrusObj.Error(err)
			}
			continue
		}

		if err = reader.CommitMessages(context.Background(), msg); err != nil && log.LogrusObj != nil {
			log.LogrusObj.Error(err)
		}
	}
}

func handleFlashSaleOrder(payload *model.FlashSale2MQ) error {
	ctx := context.Background()
	if _, err := dao.NewAddressDao(ctx).GetAddressByAid(payload.AddressId, payload.UserId); err != nil {
		return err
	}

	return dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		flashSaleDao := dao.NewFlashSaleDaoByDB(tx)
		exists, err := flashSaleDao.HasAsyncOrder(payload.FlashSaleId, payload.UserId)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		if err = flashSaleDao.CreateAsyncOrder(payload); err != nil {
			return err
		}

		order := &model.Order{
			UserID:    payload.UserId,
			ProductID: payload.ProductId,
			BossID:    payload.BossId,
			AddressID: payload.AddressId,
			Num:       1,
			OrderNum:  buildOrderNum(payload.ProductId, payload.UserId),
			Type:      consts.OrderTypeUnPaid,
			Money:     payload.Money,
		}

		return dao.NewOrderDaoByDB(tx).CreateOrder(order)
	})
}

func buildOrderNum(productID, userID uint) uint64 {
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	number = number + strconv.Itoa(int(productID)) + strconv.Itoa(int(userID))
	orderNum, _ := strconv.ParseUint(number, 10, 64)
	return orderNum
}
