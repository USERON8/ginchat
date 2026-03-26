// pkg/redis/redis.go
package redis

import (
	"context"
	"fmt"
	"ginchat/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			config.Cfg.Redis.Host,
			config.Cfg.Redis.Port,
		),
	})

	// 测试连通性
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}
