package models

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRealIP 获取真实IP（处理代理）
func GetRealIP(c *gin.Context) string {
	// 优先级顺序
	headers := []string{
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, header := range headers {
		if ip := c.GetHeader(header); ip != "" {
			// 处理多个IP的情况（如 X-Forwarded-For: client, proxy1, proxy2）
			if strings.Contains(ip, ",") {
				ips := strings.Split(ip, ",")
				ip = strings.TrimSpace(ips[0])
			}

			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// 最后使用Gin的ClientIP方法
	return c.ClientIP()
}

func LogAdd(c *gin.Context, msg string) {

	var logtemp APIRequestLog_

	logtemp.IPAddress = GetRealIP(c)
	logtemp.Path = c.Request.URL.Path
	logtemp.Method = c.Request.Method
	logtemp.Message = msg

	//fmt.Println(logtemp)
	DB.Create(&logtemp)

}
