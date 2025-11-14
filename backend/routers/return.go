package routers

import "github.com/gin-gonic/gin"

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
