// main.go
package main

import (
	"fmt"
	"ginchat/internal/ws"
	"ginchat/pkg/config"
	"ginchat/pkg/database"
	"ginchat/pkg/logger"
	"ginchat/pkg/redis"
	"ginchat/router"

	"go.uber.org/zap"
)

func main() {
	// 1. 基础配置初始化
	config.Init()
	database.Init()
	redis.Init()

	// 2. 日志初始化
	logger.Init()
	// 关键：在程序退出前，将缓冲区内的日志刷新到磁盘/控制台
	// 这里的 _ 是因为在某些环境下对标准输出 Sync 会报错，可以直接忽略
	defer func() {
		_ = logger.Log.Sync()
	}()

	// 3. 业务逻辑启动
	ws.StartWorkers(5) // 启动 5 个 Worker

	// 4. 路由初始化 (注意：请确保你的 router.Init() 内部使用了 logger.Log)
	r := router.Init()

	logger.Info("服务启动", zap.Int("port", config.Cfg.Server.Port))
	err := r.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port))
	if err != nil {
		return
	}
}
