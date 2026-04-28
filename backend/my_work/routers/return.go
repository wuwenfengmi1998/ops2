package routers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var ErrorCode map[string]interface{}

func ReturnInit() {
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

func DebugPrintJson(data map[string]interface{}) {
	p, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println("\n", string(p))

}

func ReturnJson(ctx *gin.Context, errMsg string, data map[string]interface{}) {
	var errCode = ErrorCode[errMsg]
	returnData := map[string]interface{}{}

	// cookie, have_cookie := ctx.Get("cookie")
	// if have_cookie {
	// 	returnData["cookie"] = cookie
	// }

	returnData["err_code"] = errCode
	returnData["err_msg"] = errMsg
	if data != nil {
		returnData["return"] = data
	}

	ctx.JSON(200, &returnData)

	//ctx.Abort()

}

func ReturnFile(ctx *gin.Context, file_info *TabFileInfo, preview bool) {
	if preview {
		ctx.File(file_info.Path)
	} else {
		//需要从数据库拉取原始文件名
		ctx.FileAttachment(file_info.Path, file_info.Name)
	}

}
