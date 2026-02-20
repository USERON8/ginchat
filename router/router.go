package router

import (
	"ginchat/api"
	"ginchat/global"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化Gin路由
func InitRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(global.Config.Server.Mode)
	r := gin.Default()

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 用户路由组
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/info", api.GetUserInfo)
	}

	return r
}
