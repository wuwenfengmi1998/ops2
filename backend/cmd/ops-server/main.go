package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"ops/api"
	"ops/internal/config"
	"ops/internal/database"
	"ops/internal/middleware"
)

func main() {
	// 创建日志记录器
	logger, err := createLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// 加载配置
	configPath := "./data/config.yaml"
	if err := config.Load(configPath); err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer database.Close()

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		logger.Warn("Auto migration failed", zap.Error(err))
	}

	// 设置Gin模式
	if config.Current.Web.Host == "127.0.0.1" || config.Current.Web.Host == "localhost" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin实例（使用自定义logger）
	r := gin.New()

	// 注册中间件
	r.Use(middleware.CORS())
	
	// 根据环境选择日志中间件
	if config.Current.Web.Host == "127.0.0.1" || config.Current.Web.Host == "localhost" {
		// 开发环境使用简易日志
		r.Use(middleware.SimpleLogger())
	} else {
		// 生产环境使用高级日志
		r.Use(middleware.Logger(logger))
	}
	
	r.Use(middleware.Recovery(logger))

	// 注册API路由
	registerRoutes(r, logger)

	// 静态文件服务中间件（由api.CreateRouter处理）
	// 这里仅创建dist目录（如果不存在）
	ensureDistDirectory(logger)

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%s", config.Current.Web.Host, config.Current.Web.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// 优化服务器配置
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// 注册API路由
func registerRoutes(r *gin.Engine, logger *zap.Logger) {
	// 使用我们的路由配置
	api.RegisterAllRoutes(r)
	
	// 健康检查端点（额外添加）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "Server is healthy",
			"data": gin.H{
				"timestamp": time.Now().Unix(),
				"status":    "running",
				"version":   "1.0.0",
			},
		})
	})
}

// 兼容性路由，保持原有API结构
// 注意：此函数现在由api包统一处理，此函数保留作为参考
func compatRoutes(api *gin.RouterGroup, logger *zap.Logger) {
	logger.Info("兼容性路由由api包统一管理")
}

// 静态文件服务（已迁移到api包）

// 创建日志记录器
func createLogger() (*zap.Logger, error) {
	// 开发环境使用开发配置
	if gin.Mode() == gin.DebugMode {
		return zap.NewDevelopment()
	}
	// 生产环境使用生产配置
	return zap.NewProduction()
}

// ensureDistDirectory 确保dist目录存在
func ensureDistDirectory(logger *zap.Logger) {
	if _, err := os.Stat("./dist"); os.IsNotExist(err) {
		if err := os.MkdirAll("./dist", 0755); err != nil {
			logger.Warn("Failed to create dist directory", zap.Error(err))
		} else {
			logger.Info("Created empty dist directory for static files")
		}
	}
}

