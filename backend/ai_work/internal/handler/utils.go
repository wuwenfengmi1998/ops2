package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIntParam 获取整数参数
func GetIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetUintParam 获取uint参数
func GetUintParam(c *gin.Context, key string) uint {
	value := c.Param(key)
	if value == "" {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return uint(intValue)
}