package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type CostItem struct {
	Cost         float64 `json:"cost" binding:"required,min=0"`    // 必须，非负
	CostT        float64 `json:"cost_t" binding:"required,min=0"`  // 必须，非负
	CurrencyType string  `json:"currency_type" binding:"required"` // 必须
	Int          int     `json:"int" binding:"required"`           // 必须
	Type         string  `json:"type" binding:"required"`          // 必须
}
type From_purchase_addorder struct {
	Costs          []CostItem `json:"costs"`                           //
	Link           string     `json:"link"`                            // 可选
	OrderStatus    string     `json:"order_status" binding:"required"` //
	PartName       string     `json:"part_name"`                       // 可选
	Photos         []string   `json:"photos"`                          // 可选
	Remark         string     `json:"remark"`                          // 可选
	Styles         string     `json:"styles"`                          // 可选
	Title          string     `json:"title" binding:"required"`        // 必须
	TrackingNumber string     `json:"tracking_number"`                 // 可选
	UpdateTime     string     `json:"update_time"`                     // 可选
}

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/addorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {

			//需要处理提交的数据，接口有固定的数据格式，不允许乱搞
			//fmt.Println(isAuth)
			fmt.Println(user)
			//DebugPrintJson(data)
			var jsondata From_purchase_addorder
			if err := mapstructure.Decode(data, &jsondata); err == nil {

				fmt.Println("转换后数据:\n", jsondata)

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}

		} else {
			ReturnJson(ctx, "jsonErr", nil)
		}

		ReturnJson(ctx, "apiErr", nil)
	})

}
