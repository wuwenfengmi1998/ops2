package routers

import (
	"github.com/gin-gonic/gin"
)

// 版本信息，由 main.go 在启动前赋值（值来自 -ldflags 注入）
var (
	GitVersion = "dev"
	GitCommit  = "unknown"
	BuildTime  = "unknown"
)

// 把数据分离成cookie和json
func SeparateData(ctx *gin.Context) (map[string]interface{}, string) {
	var jsonData map[string]interface{}

	if err := ctx.ShouldBindJSON(&jsonData); err == nil {
		//分离数据
		cookie, ok := jsonData["userCookieValue"].(string)
		if !ok {
			cookie = ""
		}

		data, ok := jsonData["data"].(map[string]interface{})
		if !ok {
			data = nil
		}

		return data, cookie
	}

	return nil, ""

}

func ApiRoot(r *gin.RouterGroup) {

	ApiStatic(r.Group("/static"))
	ApiUser(r.Group("/users"))
	ApiFiles(r.Group("/files"))
	ApiPurchase(r.Group("/purchase"))
	ApiSchedule(r.Group("/schedule"))
	ApiWorkOrder(r.Group("/work_order"))
	ApiWarehouse(r.Group("/warehouse"))
	r.GET("/", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", gin.H{
			"isOpsApiRoot": true,
			"version":      GitVersion,
			"gitCommit":    GitCommit,
			"buildTime":    BuildTime,
		})
	})

}
