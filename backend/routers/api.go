package routers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
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

	r.Use(func(ctx *gin.Context) {
		//转换传进来的数据
		var jsonData map[string]interface{}
		if err := ctx.ShouldBindJSON(&jsonData); err == nil {
			//分离数据

			if jsonData["cookie"] != "" && jsonData["cookie"] != nil {
				ctx.Set("cookie_value", jsonData["cookie"])
			}

			if jsonData["data"] != nil {
				//fmt.Println(jsonData["data"])
				var data_t map[string]interface{}
				if err = mapstructure.Decode(jsonData["data"], &data_t); err == nil {
					ctx.Set("data", &data_t)
				}
			}

		}
	})

	ApiUser(r.Group("/users"))

	r.GET("/", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

}
