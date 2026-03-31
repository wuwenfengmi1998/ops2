package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthToken 认证令牌中间件
// 兼容现有的 userCookieValue 字段
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取认证
		authHeader := c.GetHeader("Authorization")
		
		// 如果没有Authorization头，尝试从POST数据中获取
		if authHeader == "" && c.Request.Method == http.MethodPost {
			var requestData map[string]interface{}
			
			// 尝试解析JSON body
			if c.Request.Body != nil && c.Request.ContentLength > 0 {
				// 先读取请求体内容
				requestBody, err := io.ReadAll(c.Request.Body)
				if err == nil {
					// 重置body以便后续使用
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
					
					// 尝试解析JSON
					if err := c.ShouldBindJSON(&requestData); err == nil {
						if cookieValue, ok := requestData["userCookieValue"].(string); ok && cookieValue != "" {
							c.Set("userCookieValue", cookieValue)
							c.Set("authMethod", "cookie_value")
							c.Set("authValid", true)
							c.Next()
							return
						}
					}
					// 如果JSON解析失败，重置body并继续
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				}
			}
			
			// 尝试从表单数据获取
			if cookieValue := c.PostForm("userCookieValue"); cookieValue != "" {
				c.Set("userCookieValue", cookieValue)
				c.Set("authMethod", "cookie_value")
				c.Set("authValid", true)
				c.Next()
				return
			}
		}

		// Bearer token 认证
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			c.Set("authToken", token)
			c.Set("authMethod", "bearer_token")
			c.Set("authValid", true)
			c.Next()
			return
		}

		// 检查URL查询参数中的cookie
		if cookieValue := c.Query("userCookieValue"); cookieValue != "" {
			c.Set("userCookieValue", cookieValue)
			c.Set("authMethod", "cookie_query")
			c.Set("authValid", true)
			c.Next()
			return
		}

		// 验证失败
		c.Set("authValid", false)
		c.Next()
	}
}

// AuthRequired 需要认证的中间件
// 如果用户未认证，返回401错误
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先运行认证中间件
		authMiddleware := AuthToken()
		authMiddleware(c)

		// 检查认证结果
		if authValid, exists := c.Get("authValid"); !exists || !authValid.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "401",
				"message": "Authentication required",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminRequired 需要管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行基础认证
		AuthRequired()(c)

		// 如果请求被中止（认证失败），直接返回
		if c.IsAborted() {
			return
		}

		// TODO: 检查用户是否为管理员
		// 暂时允许所有已认证用户
		c.Next()
	}
}

// GetAuthMethod 获取认证方法
func GetAuthMethod(c *gin.Context) string {
	if method, exists := c.Get("authMethod"); exists {
		return method.(string)
	}
	return ""
}

// GetCookieValue 获取用户cookie值
func GetCookieValue(c *gin.Context) string {
	if cookie, exists := c.Get("userCookieValue"); exists {
		return cookie.(string)
	}
	return ""
}

// GetAuthToken 获取Bearer token
func GetAuthToken(c *gin.Context) string {
	if token, exists := c.Get("authToken"); exists {
		return token.(string)
	}
	return ""
}