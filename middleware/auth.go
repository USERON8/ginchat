package middleware

import (
	jwtpkg "ginchat/pkg/jwt"
	"ginchat/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			response.Fail(c, "未登录")
			c.Abort()
			return
		}

		claims, err := jwtpkg.Parse(strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			response.Fail(c, "token无效")
			c.Abort()
			return
		}

		// 把 userID 存入 context，后续 handler 直接取
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
