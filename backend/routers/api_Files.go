package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ApiFiles(r *gin.RouterGroup) {
	r.POST("/upload", func(ctx *gin.Context) {
		fmt.Print(ctx.FormFile("file"))
		ReturnJson(ctx, "apiOK", nil)
	})
	r.GET("/upload", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

}
