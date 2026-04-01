package v1

import (
	"net/http"
	"ops/internal/database"
	"ops/internal/handler"
	"ops/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	authHandler    *handler.AuthHandler
	fileHandler    *handler.FileHandler
	purchaseHandler *handler.PurchaseHandler
)

func init() {
	db = database.GetDB()
	authHandler = handler.NewAuthHandler(db)
	fileHandler = handler.NewFileHandler(db)
	purchaseHandler = handler.NewPurchaseHandler(db)
}

// RegisterRoutes 注册所有v1版本的API路由
func RegisterRoutes(r *gin.RouterGroup) {
	// API根路径测试
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "OPS API v1",
			"data":    nil,
		})
	})

	// 静态文件路由 - 保持兼容性
	r.StaticFS("/static", http.Dir("./dist"))
	
	// 用户认证相关路由
	userGroup := r.Group("/users")
	{
		// 用户认证
		userGroup.POST("/login", authHandler.UserLogin)
		userGroup.POST("/register", authHandler.UserRegister)
		userGroup.POST("/forgot-password", authHandler.UserForgotPassword)
		userGroup.POST("/reset-password", authHandler.UserResetPassword)

		// 用户信息 - 需要认证
		userGroup.PUT("/profile", middleware.AuthToken(), authHandler.UserUpdateProfile)
		userGroup.GET("/profile", middleware.AuthToken(), authHandler.UserProfile)
		userGroup.POST("/logout", middleware.AuthToken(), authHandler.UserLogout)

		// 用户管理（管理员） - TODO: 实现管理员功能
		userGroup.GET("/list", middleware.AuthToken(), adminMiddleware(), getUserList)
		userGroup.POST("/create", middleware.AuthToken(), adminMiddleware(), createUser)
		userGroup.PUT("/:id", middleware.AuthToken(), adminMiddleware(), updateUser)
		userGroup.DELETE("/:id", middleware.AuthToken(), adminMiddleware(), deleteUser)
	}

	// 文件上传管理 - v1 API
	// 注意：具体路由必须放在通配路由之前
	r.POST("/files/upload", middleware.AuthToken(), fileHandler.UploadFile)
	r.GET("/files/list", middleware.AuthToken(), fileHandler.GetFileList)
	r.GET("/files/:id", middleware.AuthToken(), fileHandler.GetFileByID)
	r.DELETE("/files/:id", middleware.AuthToken(), fileHandler.DeleteFile)
	r.GET("/files/download/:hash", fileHandler.DownloadFile)
	r.GET("/files/get/:hash", fileHandler.GetFile)

	// 采购订单管理
	purchaseGroup := r.Group("/purchase")
	{
		// 保持与前端兼容的POST路由（原始API使用POST）
		purchaseGroup.POST("/getorders", middleware.AuthToken(), purchaseHandler.GetOrders)
		purchaseGroup.POST("/addorder", middleware.AuthToken(), purchaseHandler.CreateOrder)
		
		// RESTful风格的新API
		purchaseGroup.GET("/orders", middleware.AuthToken(), purchaseHandler.GetOrders)
		purchaseGroup.POST("/orders", middleware.AuthToken(), purchaseHandler.CreateOrder)
		purchaseGroup.GET("/orders/:id", middleware.AuthToken(), purchaseHandler.GetOrderDetails)
		
		// TODO: 实现更新、删除和其他功能
		purchaseGroup.PUT("/orders/:id", middleware.AuthToken(), purchaseUpdateOrder)
		purchaseGroup.DELETE("/orders/:id", middleware.AuthToken(), purchaseDeleteOrder)
		purchaseGroup.POST("/orders/:id/costs", middleware.AuthToken(), purchaseAddCost)
		purchaseGroup.PUT("/orders/:id/costs/:costId", middleware.AuthToken(), purchaseUpdateCost)
		purchaseGroup.DELETE("/orders/:id/costs/:costId", middleware.AuthToken(), purchaseDeleteCost)
	}

	// 系统管理（管理员）
	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", systemStatus)
		systemGroup.GET("/config", middleware.AuthToken(), adminMiddleware(), systemGetConfig)
		systemGroup.PUT("/config", middleware.AuthToken(), adminMiddleware(), systemUpdateConfig)
	}
}

// 管理员中间件占位函数
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现管理员权限检查
		c.Next()
	}
}

// 占位函数 - 将在后续步骤中实现
func getUserList(c *gin.Context)         {}
func createUser(c *gin.Context)          {}
func updateUser(c *gin.Context)          {}
func deleteUser(c *gin.Context)          {}
func purchaseUpdateOrder(c *gin.Context) {}
func purchaseDeleteOrder(c *gin.Context) {}
func purchaseAddCost(c *gin.Context)     {}
func purchaseUpdateCost(c *gin.Context)  {}
func purchaseDeleteCost(c *gin.Context)  {}
func systemStatus(c *gin.Context)        {}
func systemGetConfig(c *gin.Context)     {}
func systemUpdateConfig(c *gin.Context)  {}