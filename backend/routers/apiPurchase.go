package routers

import (
	"encoding/json"
	"fmt"
	"ops/models"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"gorm.io/datatypes"
)

type CostItem struct {
	Cost         int    `json:"cost"`         // 必须，非负
	CostT        int    `json:"costt"`        // 必须，非负
	CurrencyType string `json:"currencytype"` // 必须
	Int          int    `json:"int"`          // 必须
	Type         string `json:"type"`         // 必须
}
type From_purchase_addorder struct {
	Costs          []CostItem `json:"costs"`           //
	Link           string     `json:"link"`            // 可选
	OrderStatus    string     `json:"order_status"`    //
	PartName       string     `json:"partname"`        // 可选
	Photos         []string   `json:"photos"`          // 可选
	Remark         string     `json:"remark"`          // 可选
	Styles         string     `json:"styles"`          // 可选
	Title          string     `json:"title"`           // 必须
	TrackingNumber string     `json:"tracking_number"` // 可选
	UpdateTime     string     `json:"update_time"`     // 可选
}

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/getorders", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			fmt.Println(user)
			// DebugPrintJson(data)

			type From_purchase_getorders struct {
				Search  string
				Entries int
				Page    int
			}

			var jsondata From_purchase_getorders
			if err := mapstructure.Decode(data, &jsondata); err == nil {
				//fmt.Println(jsondata)

				is_data_ok := true

				if jsondata.Entries <= 0 {
					is_data_ok = false
				}
				if jsondata.Page <= 0 {
					is_data_ok = false
				}

				if is_data_ok {

					//读取有多少条目
					var count int64
					models.DB.Model(&models.TabPurchaseOrder{}).Count(&count)
					fmt.Println(count)

					//读取条目

				} else {
					ReturnJson(ctx, "jsonErr", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}

		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	r.POST("/addorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {

			//需要处理提交的数据，接口有固定的数据格式，不允许乱搞
			//fmt.Println(isAuth)
			//fmt.Println(user)
			//DebugPrintJson(data)
			var jsondata From_purchase_addorder
			if err := mapstructure.Decode(data, &jsondata); err == nil {

				//fmt.Println("转换后数据:\n", jsondata)

				//数据比较混乱 在这里校验

				//判断标题不为空
				is_data_ok := true
				if jsondata.Title == "" {
					is_data_ok = false
				}

				//判断数量与价格是否为负数
				for i := 0; i < len(jsondata.Costs); i++ {
					if jsondata.Costs[i].Cost <= 0 {
						is_data_ok = false
					}
					if jsondata.Costs[i].Int <= 0 {
						is_data_ok = false
					}
				}

				//判断图片是否为哈希值
				for i := 0; i < len(jsondata.Photos); i++ {
					//判断字符串是否包含标点符号
					if models.IsContainsSpecialChar(jsondata.Photos[i]) {
						is_data_ok = false
					}

				}

				//判断时间字符串是否合法
				uptime, e := models.StringToTimePtr(jsondata.UpdateTime)
				if e != nil {
					is_data_ok = false
				}

				if is_data_ok {
					//校验通过
					//载入数据库

					photos, _ := json.Marshal(jsondata.Photos) //把图片数组转换成字符串
					new_data := models.TabPurchaseOrder{
						UserID:         user.ID,
						Title:          jsondata.Title,
						Remark:         jsondata.Remark,
						Photos:         datatypes.JSON(photos),
						Link:           jsondata.Link,
						PartName:       jsondata.PartName,
						Styles:         jsondata.Styles,
						UpdateTime:     uptime,
						TrackingNumber: jsondata.TrackingNumber,
						OrderStatus:    jsondata.OrderStatus,
					}
					models.DB.Create(&new_data)

					for i := 0; i < len(jsondata.Costs); i++ {
						new_cost_data := models.TabPurchaseCosts{
							Price:    jsondata.Costs[i].Cost,
							Quantity: jsondata.Costs[i].Int,
							UserID:   user.ID,
							OrderID:  new_data.ID,
						}
						models.DB.Create(&new_cost_data)
					}

				} else {
					ReturnJson(ctx, "jsonErr", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}

		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

		ReturnJson(ctx, "apiErr", nil)
	})

}
