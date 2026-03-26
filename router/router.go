// router/router.go
package router

import (
	"ginchat/internal/handler"
	"ginchat/middleware"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		user := api.Group("/user")
		user.POST("/register", handler.Register)
		user.POST("/login", handler.Login)
	}

	// 需要登录的路由
	auth := api.Group("").Use(middleware.Auth())
	{
		// 用户
		auth.GET("/user/info", handler.GetUserInfo)
		auth.PUT("/user/info", handler.UpdateUserInfo)
		auth.GET("/user/search", handler.SearchUser)
		auth.POST("/user/logout", handler.Logout)

		// 好友
		auth.POST("/friend/apply/:id", handler.ApplyFriend)
		auth.PUT("/friend/apply", handler.HandleApply)
		auth.GET("/friend/applies", handler.GetApplyList)
		auth.GET("/friend/list", handler.GetFriendList)
		auth.DELETE("/friend", handler.DeleteFriend)

		// 消息
		auth.GET("/message/history", handler.GetMessageHistory)
		auth.GET("/message/unread", handler.GetUnreadCount)

		// WebSocket
		auth.GET("/ws", handler.WsHandler)
	}

	return r
}
