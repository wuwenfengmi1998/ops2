package routers

import (
	"ops/models"
	"path"

	"github.com/gin-gonic/gin"
)

//处理api的静态内容

func ApiStatic(r *gin.RouterGroup) {
	r.GET("/avatar/:filename", func(ctx *gin.Context) {
		filename := ctx.Param("filename")
		dst := path.Join(models.ConfigsFile.Pahts["avatar"], filename)
		if models.FileExists(dst) {
			ctx.File(dst)
		} else {
			//找不到文件
			ctx.String(404, "file not found")

		}

	})
}
