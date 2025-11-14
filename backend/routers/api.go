package routers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var ErrorCode map[string]interface{}

func init() {
	//读取默认配置
	fmt.Println("尝试读取错误码文件")
	data, err := os.ReadFile("./defConfig/errorCodes.json")
	if err != nil {

		fmt.Println("读取错误码文件失败", err)
	}

	if err := json.Unmarshal(data, &ErrorCode); err != nil {
		fmt.Println("解析错误码文件失败", err)
	}

}

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

	ApiUser(r.Group("/users"))

	r.GET("/", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

}
