// main.go
package main

import (
	"fmt"
	"ginchat/internal/ws"
	"ginchat/pkg/config"
	"ginchat/pkg/database"
	"ginchat/pkg/redis"
	"ginchat/router"
)

func main() {
	config.Init()
	database.Init()
	redis.Init()

	ws.StartWorkers(5) // 启动 5 个 Worker

	r := router.Init()

	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
