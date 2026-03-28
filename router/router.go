// router/router.go
package router

import (
	"ginchat/internal/handler"
	"ginchat/middleware"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.New()                // 不用 Default()
	r.Use(middleware.Logger())    // 用自己的日志中间件
	r.Use(gin.Recovery())         // 保留 panic 恢复
	r.Use(middleware.RateLimit()) // 全局限流
	api := r.Group("/api")
	{
		user := api.Group("/user")
		user.POST("/register", handler.Register)
		user.POST("/login", handler.Login)
		user.POST("/user/register", middleware.StrictRateLimit(), handler.Register)
		user.POST("/user/login", middleware.StrictRateLimit(), handler.Login)
		user.POST("/user/refresh", handler.RefreshToken) // 刷新 token

	}

	// 需要登录的路由
	auth := api.Group("").Use(middleware.Auth())
	{
		// 用户
		auth.GET("/user/info", handler.GetUserInfo)
		auth.PUT("/user/info", handler.UpdateUserInfo)
		auth.GET("/user/search", handler.SearchUser)
		auth.POST("/user/logout", handler.Logout)
		auth.PUT("/user/password", handler.ChangePassword) // 修改密码

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
		auth.POST("/group", handler.CreateGroup)
		auth.GET("/group/list", handler.GetGroupList)
		auth.GET("/group/:id", handler.GetGroup)
		auth.POST("/group/:id/member", handler.AddGroupMember)
		auth.DELETE("/group/:id/member", handler.RemoveGroupMember)
		auth.GET("/group/:id/members", handler.GetGroupMembers)
		auth.GET("/group/:id/messages", handler.GetGroupMessageHistory)

		auth.PUT("/message/read", handler.MarkRead) // 改成单独接口
	}

	return r
}
