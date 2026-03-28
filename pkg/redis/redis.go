// pkg/redis/redis.go
package redis

import (
	"context"
	"fmt"
	"ginchat/pkg/config"
	"ginchat/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client

func Init() {
	cfg := config.Cfg.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Redis 连接失败", zap.Error(err))
	}
}
