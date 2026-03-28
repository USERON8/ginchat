package middleware

import (
	jwtpkg "ginchat/pkg/jwt"
	"ginchat/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		// 先从 Header 取
		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}

		// Header 没有，从 URL 参数取（WebSocket 用）
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtpkg.Parse(tokenStr)
		if err != nil {
			response.Fail(c, response.CodeTokenInvalid)
			c.Abort()
			return
		}
		// 把 userID 存入 context，后续 handler 直接取
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
