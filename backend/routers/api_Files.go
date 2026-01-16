package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func file_save() {

}

func ApiFiles(r *gin.RouterGroup) {
	r.POST("/upload", func(ctx *gin.Context) {

		cookie := ctx.PostForm("cookie")
		file, _ := ctx.FormFile("file")
		//通过cookie获取用户信息
		_, err := AuthenticationAuthorityFromCookie(cookie)
		if err == nil {

		}

		fmt.Println(file.Filename)
		fmt.Println(cookie)
		ReturnJson(ctx, "apiOK", nil)
	})
	r.GET("/upload", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

}
