package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	conf "github.com/YasinDoyle/e-mall/config"
	log "github.com/YasinDoyle/e-mall/utils/log"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.Mutex
)

func InitRabbitMQ() {
	mu.Lock()
	defer mu.Unlock()

	if channel != nil {
		return
	}
	if conf.Config == nil || conf.Config.RabbitMq == nil {
		return
	}

	rConfig := conf.Config.RabbitMq
	uri := fmt.Sprintf("%s://%s:%s@%s:%s/",
		rConfig.RabbitMQ,
		rConfig.RabbitMQUser,
		rConfig.RabbitMQPassWord,
		rConfig.RabbitMQHost,
		rConfig.RabbitMQPort,
	)

	conn, err := amqp.Dial(uri)
	if err != nil {
		if log.LogrusObj != nil {
			log.LogrusObj.Error(err)
		}
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		if log.LogrusObj != nil {
			log.LogrusObj.Error(err)
		}
		return
	}

	connection = conn
	channel = ch
}

func PublishJSON(queue string, payload interface{}) error {
	if channel == nil {
		InitRabbitMQ()
	}
	if channel == nil {
		return fmt.Errorf("rabbitmq unavailable")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	_, err = channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
