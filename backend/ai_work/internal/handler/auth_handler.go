package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"ops/internal/service"
	"ops/pkg/response"
)

// AuthHandler 用户认证处理器
type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	DeviceID string `json:"deviceID"`
	IP       string `json:"ip"`
	Remember string `json:"remember"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	UserID           uint   `json:"userID"`
	Name             string `json:"name"`
	AvatarURL        string `json:"avatarURL"`
	CookieValue      string `json:"cookieValue"`
	CookieExpireDate string `json:"cookieExpireDate"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// RegisterResponse 注册响应结构
type RegisterResponse struct {
	UserID      uint   `json:"userID"`
	Name        string `json:"name"`
	CookieValue string `json:"cookieValue"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone" binding:"omitempty,len=11"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=50"`
}

// LogoutRequest 退出登录请求
type LogoutRequest struct {
	CookieValue string `json:"cookieValue" binding:"required"`
	DeviceID    string `json:"deviceID"`
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(db),
		validate:    validator.New(),
	}
}

// UserLogin 用户登录
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req LoginRequest
	
	// 绑定和验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	user, cookie, err := h.authService.Login(req.Name, req.Password, req.DeviceID, req.IP, req.Remember)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, "-5", "User not found")
			return
		}
		if strings.Contains(err.Error(), "password") {
			response.Error(c, "-42", "Invalid password")
			return
		}
		response.InternalError(c, err)
		return
	}

	// 构建响应
	resp := LoginResponse{
		UserID:           user.UserID,
		Name:             user.Name,
		AvatarURL:        user.AvatarURL,
		CookieValue:      cookie.Value,
		CookieExpireDate: cookie.ExpireDate.Format("2006-01-02 15:04:05"),
	}

	response.Success(c, resp)
}

// UserRegister 用户注册
func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req RegisterRequest
	
	// 绑定和验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	user, cookie, err := h.authService.Register(req.Name, req.Password, req.Email, req.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			response.Error(c, "-4", "Username already exists")
			return
		}
		response.InternalError(c, err)
		return
	}

	// 构建响应
	resp := RegisterResponse{
		UserID:      user.UserID,
		Name:        user.Name,
		CookieValue: cookie.Value,
	}

	response.Success(c, resp)
}

// UserForgotPassword 忘记密码
func (h *AuthHandler) UserForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	
	// 绑定和验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 至少需要邮箱或手机号之一
	if req.Email == "" && req.Phone == "" {
		response.BadRequest(c, "Email or phone number is required")
		return
	}

	// 验证请求参数
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	token, err := h.authService.ForgotPassword(req.Name, req.Email, req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, "-5", "User not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	// 构建响应
	response.Success(c, gin.H{
		"resetToken": token,
		"message":    "Password reset instructions have been sent",
	})
}

// UserResetPassword 重置密码
func (h *AuthHandler) UserResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	
	// 绑定和验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	err := h.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") {
			response.Error(c, "-2", "Reset token is invalid or expired")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, "-5", "User not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Password has been reset successfully",
	})
}

// UserLogout 用户退出登录
func (h *AuthHandler) UserLogout(c *gin.Context) {
	var req LogoutRequest
	
	// 从认证中间件获取cookie值
	cookieValue := getCookieFromContext(c)
	if cookieValue == "" {
		// 尝试从请求body获取
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "Invalid request format")
			return
		}
		cookieValue = req.CookieValue
	}

	if cookieValue == "" {
		response.BadRequest(c, "Cookie value is required")
		return
	}

	// 从请求中获取设备ID
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" && req.DeviceID != "" {
		deviceID = req.DeviceID
	}

	// 调用服务层
	err := h.authService.Logout(cookieValue, deviceID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Logged out successfully",
	})
}

// UserProfile 获取用户信息
func (h *AuthHandler) UserProfile(c *gin.Context) {
	// 从认证中间件获取用户ID或名称
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c)
		return
	}

	// 调用服务层
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, "-5", "User not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, user)
}

// UserUpdateProfile 更新用户信息
func (h *AuthHandler) UserUpdateProfile(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c)
		return
	}

	// 解析更新请求
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 禁止更新某些字段
	delete(updateData, "id")
	delete(updateData, "name")
	delete(updateData, "password")
	delete(updateData, "createdAt")

	// 调用服务层
	user, err := h.authService.UpdateProfile(userID, updateData)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, user)
}

// 辅助函数
func getCookieFromContext(c *gin.Context) string {
	if cookie, exists := c.Get("userCookieValue"); exists && cookie != "" {
		return cookie.(string)
	}
	return ""
}

func getUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}