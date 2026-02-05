package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/addorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			fmt.Println(isAuth)
			fmt.Println(user)
			fmt.Println(data)

		} else {
			ReturnJson(ctx, "jsonErr", nil)
		}

		ReturnJson(ctx, "apiErr", nil)
	})

}
