package global

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port),
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
		PoolSize: Config.Redis.PoolSize,
	})

	// 测试连接
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	Log.Info("Redis初始化成功")
	return nil
}
