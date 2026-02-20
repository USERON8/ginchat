package main

import (
	"fmt"
	"ginchat/global"
	"ginchat/models"
	"ginchat/router"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化配置（必须第一个）
	if err := global.InitConfig(); err != nil {
		panic(fmt.Sprintf("初始化配置失败: %v", err))
	}

	// 2. 初始化日志
	if err := global.InitLogger(); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {

		}
	}(global.Log) // 程序退出时刷新日志

	// 3. 初始化数据库
	if err := global.InitDB(); err != nil {
		global.Log.Fatal("初始化数据库失败", zap.Error(err))
	}

	// 4. 自动迁移表结构（新增模型后在这里添加）
	if err := global.DB.AutoMigrate(
		&models.UserBasic{},
		// 后续添加：&models.Message{}, &models.Group{}...
	); err != nil {
		global.Log.Fatal("表结构迁移失败", zap.Error(err))
	}

	// 5. 初始化Redis
	if err := global.InitRedis(); err != nil {
		global.Log.Fatal("初始化Redis失败", zap.Error(err))
	}

	// 6. 初始化Gin路由
	r := router.InitRouter()

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Log.Info(fmt.Sprintf("IM服务启动成功，监听端口: %d", global.Config.Server.Port))
	if err := r.Run(addr); err != nil {
		global.Log.Fatal("服务启动失败", zap.Error(err))
	}
}
