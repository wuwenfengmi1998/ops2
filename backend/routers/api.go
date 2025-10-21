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

func ApiRoot(r *gin.RouterGroup) {

	r.GET("/", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

}
