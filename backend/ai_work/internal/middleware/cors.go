package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许所有来源（生产环境应指定具体域名）
		AllowOrigins: []string{"*"},
		
		// 允许的方法
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		},
		
		// 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"X-Request-ID",
			"X-Requested-With",
			"Accept",
			"Cache-Control",
			// 自定义头
			"User-Cookie-Value", // 兼容现有系统
		},
		
		// 暴露的响应头
		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
			"X-Request-ID",
			"Content-Disposition",
		},
		
		// 是否允许携带凭证
		AllowCredentials: true,
		
		// 预检请求缓存时间（秒）
		MaxAge: 12 * time.Hour,
		
		// 允许读取自定义头
		AllowPrivateNetwork: true,
	})
}

// CORSMiddleware 简化的CORS中间件（兼容老版本）
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		
		// 设置CORS头
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID, X-Requested-With, Accept, Cache-Control, User-Cookie-Value")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}