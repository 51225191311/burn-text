package middleware

import (
	"burn-text/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// IPwRateLimiter 基于IP的限流中间件
func IPwRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取客户端IP
		ip := c.ClientIP()

		//检查限流
		allowed, err := storage.AllowRequest(ip, 5, time.Minute)
		if err != nil {
			//如果Redis出错，选择放行，记录日志
			fmt.Println("限流器 Redis 错误：%v", err)
			c.Next()
			return
		}

		if !allowed {
			//如果限流直接拦截，返回429 Too Many Requests
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请一分钟后再试。",
			})
			c.Abort()
			return
		}

		//运行无误，放行
		c.Next()
	}
}
