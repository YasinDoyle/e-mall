package cache

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/redis/go-redis/v9"
	logging "github.com/sirupsen/logrus"

	conf "github.com/YasinDoyle/e-mall/config"
	trackutil "github.com/YasinDoyle/e-mall/utils/track"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

// InitCache 在中间件中初始化redis链接
func InitCache() {
	rConfig := conf.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	client.AddHook(redisTracingHook{})
	RedisClient = client
}

type redisTracingHook struct{}

func (redisTracingHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (redisTracingHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		span, spanCtx := trackutil.WithSpan(ctx, fmt.Sprintf("redis.%s", cmd.FullName()))
		defer span.Finish()

		err := next(spanCtx, cmd)
		if err != nil && err != redis.Nil {
			ext.Error.Set(span, true)
			span.SetTag("error.message", err.Error())
		}

		return err
	}
}

func (redisTracingHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		span, spanCtx := trackutil.WithSpan(ctx, "redis.pipeline")
		defer span.Finish()

		err := next(spanCtx, cmds)
		if err != nil && err != redis.Nil {
			ext.Error.Set(span, true)
			span.SetTag("error.message", err.Error())
		}

		return err
	}
}
