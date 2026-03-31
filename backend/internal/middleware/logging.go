package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogResponseWriter 自定义ResponseWriter以捕获响应内容
type LogResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *LogResponseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *LogResponseWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

// Logger 请求日志中间件
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		
		// 请求方法
		httpMethod := c.Request.Method
		
		// 请求路径
		reqUri := c.Request.RequestURI
		
		// 客户端IP
		clientIP := c.ClientIP()
		
		// 用户代理
		userAgent := c.Request.UserAgent()
		
		// 请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Set("requestID", requestID)
		} else {
			c.Set("requestID", requestID)
		}
		
		// 记录原始请求体（如果不是文件上传等大请求）
		var requestBody []byte
		if c.Request.ContentLength > 0 && c.Request.ContentLength < 1024*1024 && // 1MB限制
			c.Request.Header.Get("Content-Type") != "multipart/form-data" {
			// 读取请求体
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = bodyBytes
				// 重置请求体以便后续使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				// 尝试解析JSON
				var jsonBody interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					// 敏感信息过滤（如密码）
					if m, ok := jsonBody.(map[string]interface{}); ok {
						if _, exists := m["password"]; exists {
							m["password"] = "***REDACTED***"
						}
						if _, exists := m["oldPassword"]; exists {
							m["oldPassword"] = "***REDACTED***"
						}
						if _, exists := m["newPassword"]; exists {
							m["newPassword"] = "***REDACTED***"
						}
						if _, exists := m["confirmPassword"]; exists {
							m["confirmPassword"] = "***REDACTED***"
						}
					}
				}
			}
		}
		
		// 包装ResponseWriter以捕获响应
		blw := &LogResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw
		
		// 处理请求
		c.Next()
		
		// 结束时间
		endTime := time.Now()
		
		// 执行时间
		latency := endTime.Sub(startTime)
		
		// 响应状态码
		statusCode := c.Writer.Status()
		
		// 错误信息
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if errors == "" {
			errors = c.Errors.ByType(gin.ErrorTypePublic).String()
		}
		
		// 响应体（如果不是文件等大型响应）
		var responseBody interface{}
		var responseMap map[string]interface{}
		if blw.body != nil && blw.body.Len() > 0 && blw.body.Len() < 10000 { // 10KB限制
			bodyBytes := blw.body.Bytes()
			if err := json.Unmarshal(bodyBytes, &responseMap); err == nil {
				responseBody = responseMap
			} else {
				responseBody = string(bodyBytes)
			}
		}
		
		// 根据状态码决定日志级别
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", httpMethod),
			zap.String("uri", reqUri),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
		}
		
		// 添加请求体（如果存在且不是太大）
		if len(requestBody) > 0 && len(requestBody) < 10000 {
			var reqBody interface{}
			if err := json.Unmarshal(requestBody, &reqBody); err == nil {
				fields = append(fields, zap.Any("request_body", reqBody))
			}
		}
		
		// 添加响应体（如果存在且不是太大）
		if responseBody != nil {
			fields = append(fields, zap.Any("response_body", responseBody))
		}
		
		// 添加错误信息
		if errors != "" {
			fields = append(fields, zap.String("error", errors))
		}
		
		// 获取用户标识（如果有）
		if cookieValue := GetCookieValue(c); cookieValue != "" {
			fields = append(fields, zap.String("auth_cookie_truncated", truncateString(cookieValue, 8)))
		}
		if authToken := GetAuthToken(c); authToken != "" {
			fields = append(fields, zap.String("auth_token_truncated", truncateString(authToken, 8)))
		}
		
		// 记录日志
		logFunc := logger.Info
		if statusCode >= 400 && statusCode < 500 {
			logFunc = logger.Warn
		} else if statusCode >= 500 {
			logFunc = logger.Error
		}
		
		logFunc("HTTP request", fields...)
	}
}

// Recovery 恢复中间件
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求ID
				requestID, _ := c.Get("requestID")
				
				// 记录Panic
				logger.Error("HTTP panic recovered",
					zap.Any("error", err),
					zap.String("request_id", requestID.(string)),
					zap.String("method", c.Request.Method),
					zap.String("uri", c.Request.RequestURI),
					zap.String("client_ip", c.ClientIP()),
				)
				
				// 返回500错误
				c.JSON(500, gin.H{
					"code":    "500",
					"message": "Internal server error",
					"data":    nil,
				})
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// 辅助函数
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + shortRandString()
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func shortRandString() string {
	// 简化的随机字符串生成
	return time.Now().Format("150405")
}

// SimpleLogger 简易日志中间件（用于开发和测试）
func SimpleLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 记录请求信息
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		
		// 输出到控制台
		fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}