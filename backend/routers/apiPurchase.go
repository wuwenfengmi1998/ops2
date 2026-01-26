package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/addorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {

		} else {
			ReturnJson(ctx, "jsonErr", nil)
		}
		fmt.Println(isAuth)
		fmt.Println(user)
		fmt.Println(data)
		ReturnJson(ctx, "apiErr", nil)
	})

}
