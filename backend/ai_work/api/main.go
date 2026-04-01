package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ops/api/v1"
	"ops/routers"
)

// RegisterAllRoutes 注册所有路由，包括兼容性路由
func RegisterAllRoutes(r *gin.Engine) {
	// API v1路由（RESTful风格）
	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1)

	// 兼容性API路由（保持原有路径结构）
	api := r.Group("/api")
	registerCompatibilityRoutes(api)

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/index.html")
	})
}

// registerCompatibilityRoutes 注册兼容性路由
func registerCompatibilityRoutes(api *gin.RouterGroup) {
	// 健康检查
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "API is healthy",
			"data":    nil,
		})
	})

	// 测试端点
	api.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "API test successful",
			"data":    nil,
		})
	})

	api.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "API test successful (POST)",
			"data":    nil,
		})
	})

	// 注册原有路由模块
	api.Static("/static", "./dist")
	routers.ApiStatic(api.Group("/static"))
	routers.ApiUser(api.Group("/users"))
	routers.ApiFiles(api.Group("/files"))
	routers.ApiPurchase(api.Group("/purchase"))

	// 根API路径
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "OPS API",
			"data": gin.H{
				"version": "1.0",
				"routes": []string{
					"/api/users/*",
					"/api/files/*",
					"/api/purchase/*",
					"/api/v1/* (RESTful API)",
				},
			},
		})
	})
}

// CreateRouter 创建完整路由引擎
func CreateRouter() *gin.Engine {
	r := gin.New()

	// 设置信任代理
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 注册所有路由
	RegisterAllRoutes(r)

	// 最后注册404处理
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回JSON 404
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(404, gin.H{
				"code":    "404",
				"message": "API endpoint not found",
				"data":    nil,
			})
			return
		}

		// 否则尝试提供静态文件
		fs := http.FileServer(http.Dir("./dist"))
		fs.ServeHTTP(c.Writer, c.Request)
	})

	return r
}

