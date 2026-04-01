package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorCodeMap 错误码映射
var ErrorCodeMap map[string]string

// 加载错误码
func init() {
	ErrorCodeMap = map[string]string{
		"apiOK": "API正常",
		"-1":    "内部错误",
		"-2":    "参数错误",
		"-3":    "用户未登录",
		"-4":    "用户已存在",
		"-5":    "用户不存在",
		"-6":    "密码错误",
		"-7":    "权限不足",
		"-8":    "请求频率过高",
		"-9":    "文件上传失败",
		"-10":   "文件类型不支持",
		"-11":   "文件大小超过限制",
		"-42":   "用户名或密码错误",
	}
}

// Success 成功响应
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "Success",
		Data:    data,
	})
}

// Error 错误响应
func Error(ctx *gin.Context, code string, data interface{}) {
	message := ErrorCodeMap[code]
	if message == "" {
		message = "Unknown error"
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 参数错误
func BadRequest(ctx *gin.Context, message string) {
	if message == "" {
		message = "Bad request"
	}
	ctx.JSON(http.StatusBadRequest, Response{
		Code:    "-2",
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 未授权
func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:    "-3",
		Message: "Unauthorized",
		Data:    nil,
	})
}

// Forbidden 禁止访问
func Forbidden(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:    "-7",
		Message: "Forbidden",
		Data:    nil,
	})
}

// InternalError 内部错误
func InternalError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    "-1",
		Message: "Internal server error",
		Data:    nil,
	})
}