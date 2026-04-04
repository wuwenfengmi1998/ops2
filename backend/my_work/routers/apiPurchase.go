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
	Cost         int    `json:"cost"`         // 费用
	CostT        int    `json:"costt"`        // 总价
	CurrencyType string `json:"currencytype"` // 货币类型
	Int          int    `json:"int"`          // 数量
	Type         string `json:"type"`         // 费用类型
}
type From_purchase_addorder struct {
	Costs          []CostItem `json:"costs"`           //  成本
	Link           string     `json:"link"`            //  链接
	OrderStatus    string     `json:"order_status"`    //  订单状态
	PartName       string     `json:"partname"`        //  物件名称
	Photos         []string   `json:"photos"`          //  图片备注
	Remark         string     `json:"remark"`          //  备注 
	Styles         string     `json:"styles"`          //  样式备注
	Title          string     `json:"title"`           //  标题
	TrackingNumber string     `json:"tracking_number"` //  快递单号
	UpdateTime     string     `json:"update_time"`     //  更新时间
}

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/getorders", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if isAuth {
			//fmt.Println(user)
			// DebugPrintJson(data)

			type From_purchase_getorders struct {
				Search  string
				Entries int
				Page    int
			}

			var jsondata From_purchase_getorders
			if err := mapstructure.Decode(data, &jsondata); err == nil {
				fmt.Println(jsondata)

				is_data_ok := true

				if jsondata.Entries <= 0 || jsondata.Entries > 300 {
					is_data_ok = false
				}
				if jsondata.Page <= 0 {
					is_data_ok = false
				}

				if is_data_ok {

					//读取有多少条目
					var count int64
					models.DB.Model(&models.TabPurchaseOrder{}).Count(&count)
					//fmt.Println(count)

					//读取条目
					var getorders []models.TabPurchaseOrder
					models.DB.Order("created_at DESC").Offset(jsondata.Entries * (jsondata.Page - 1)).Limit(jsondata.Entries).Find(&getorders)

					ReturnJson(ctx, "apiOK", map[string]interface{}{
						"all_count":  count,
						"all_orders": getorders,
					})

				} else {
					ReturnJson(ctx, "jsonErr_1", nil)
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

				jsonStr, _ := json.MarshalIndent(jsondata, "", "  ")
				fmt.Println("转换后数据:\n", string(jsonStr))

				//数据比较混乱 在这里校验

				//判断标题不为空
				is_data_ok := true
				if jsondata.Title == "" {
					is_data_ok = false

					fmt.Println("err1")
				}

				//判断数量与价格是否为负数
				for i := 0; i < len(jsondata.Costs); i++ {
					if jsondata.Costs[i].Cost <= 0 {
						is_data_ok = false
						fmt.Println("err2")
					}
					if jsondata.Costs[i].Int <= 0 {
						is_data_ok = false
						fmt.Println("err3")
					}
				}

				//判断图片是否为哈希值
				for i := 0; i < len(jsondata.Photos); i++ {
					//判断字符串是否包含标点符号
					if models.IsContainsSpecialChar(jsondata.Photos[i]) {
						is_data_ok = false
						fmt.Println("err4")
					}

				}

				//判断时间字符串是否合法
				uptime, e := models.StringToTimePtr(jsondata.UpdateTime)
				if e != nil {
					is_data_ok = false
					fmt.Println("err5")
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
					ReturnJson(ctx, "jsonErr_1", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}

		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

		//ReturnJson(ctx, "apiErr", nil)
	})

}
