// middleware/ratelimit.go
package middleware

import (
	"ginchat/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter *rate.Limiter
}

var (
	limiters sync.Map
)

// 每个 IP 单独一个限流器
func getLimiter(ip string) *rate.Limiter {
	val, ok := limiters.Load(ip)
	if !ok {
		// 每秒 5 个请求，最多突发 10 个
		l := rate.NewLimiter(5, 10)
		limiters.Store(ip, &ipLimiter{limiter: l})
		return l
	}
	return val.(*ipLimiter).limiter
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getLimiter(c.ClientIP())
		if !limiter.Allow() {
			response.Fail(c, response.CodeRateLimitExceeded)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 针对登录/注册更严格的限流：每秒 1 个，最多突发 3 个
func StrictRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "strict:" + c.ClientIP()
		val, ok := limiters.Load(key)
		if !ok {
			l := rate.NewLimiter(1, 3)
			limiters.Store(key, &ipLimiter{limiter: l})
			if !l.Allow() {
				response.Fail(c, response.CodeRateLimitExceeded)
				c.Abort()
				return
			}
		} else {
			if !val.(*ipLimiter).limiter.Allow() {
				response.Fail(c, response.CodeRateLimitExceeded)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
