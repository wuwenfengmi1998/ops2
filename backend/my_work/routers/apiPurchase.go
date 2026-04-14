package routers

import (
	"encoding/json"
	"ops/models"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	purchaseUserGroup models.TabUserGroups_
	purchaseAdmins    []uint
)

// 更新管理员成员缓存
func updatePurchaseAdminsCash() {
	purchaseAdmins = nil
	// id 1 是系统管理员
	purchaseAdmins = append(purchaseAdmins, 1)
	var binds []models.TabUserGroupBinds_
	models.DB.Where("group_id = ?", purchaseUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(purchaseAdmins, item.UserID) {
			purchaseAdmins = append(purchaseAdmins, item.UserID)
		}
	}
}

// 判断用户是否有权限修改/删除订单（创建者或管理员）
func canModifyPurchase(userID, creatorUserID uint) bool {
	if slices.Contains(purchaseAdmins, userID) {
		return true
	}
	if userID == creatorUserID {
		return true
	}
	return false
}

// decodeJSON 将 map 通过 JSON 中转解码到目标结构体，绕过 mapstructure 的字段名匹配问题
func decodeJSON(data map[string]interface{}, out interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, out)
}

type CostItem struct {
	Cost         int `json:"cost"`         // 费用（分）
	CostT        int `json:"costt"`        // 总价
	CurrencyType int `json:"currencytype"` // 货币类型: 1-CNY 2-MOP 3-HKD 4-USD
	Int          int `json:"int"`          // 数量
	Type         int `json:"type"`         // 费用类型: 1-单价 2-运费
}
type From_purchase_addorder struct {
	Costs       []CostItem `json:"costs"`        //  成本
	Link        string     `json:"link"`         //  链接
	OrderStatus string     `json:"order_status"` //  订单状态
	Photos      []string   `json:"photos"`       //  图片备注
	Remark      string     `json:"remark"`       //  备注
	Styles      string     `json:"styles"`       //  样式备注
	Title       string     `json:"title"`        //  标题
}

type TabPurchaseOrder struct {
	ID          uint           `gorm:"primarykey"`
	UserID      uint           `gorm:"not null"`
	Title       string         `gorm:"size:200;comment:标题"`
	Remark      string         `gorm:"type:text;comment:备注"`
	Link        string         `gorm:"size:1000;comment:链接"`
	Styles      string         `gorm:"type:text;comment:样式数组"`
	OrderStatus string         `gorm:"size:50;default:pending;comment:订单状态: pending-待处理 ordered-已下单 arrived-已到达 received-已收件 lost-丢件 returned-退件"`
	CreatedAt   *time.Time     `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"type:datetime;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TabPurchaseCosts struct {
	ID           uint       `gorm:"primarykey"`
	OrderID      uint       `gorm:"not null"`
	UserID       uint       `gorm:"not null"`
	Price        int        `gorm:"not null"`
	Quantity     int        `gorm:"not null"`
	CurrencyType int        `gorm:"default:1;comment:货币类型: 1-CNY 2-MOP 3-HKD 4-USD"`
	CostType     int        `gorm:"default:1;comment:费用类型: 1-单价 2-运费"`
	CreatedAt    *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabPurchaseFileBind struct {
	ID        uint       `gorm:"primarykey"`
	OrderID   uint       `gorm:"not null"`
	FileID    uint       `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabPurchaseCommit 记录订单状态变更及评论
type TabPurchaseCommit struct {
	ID        uint       `gorm:"primarykey"`
	OrderID   uint       `gorm:"not null;index;comment:关联订单ID"`
	UserID    uint       `gorm:"not null;comment:操作人ID"`
	Action    string     `gorm:"size:50;not null;comment:操作类型: create-创建 create_status-状态变更"`
	Status    string     `gorm:"size:50;comment:变更后的状态"`
	OldStatus string     `gorm:"size:50;comment:变更前的状态"`
	Comment   string     `gorm:"type:text;comment:评论/备注"`
	Photos    string     `gorm:"type:text;comment:变更图片(JSON数组，存放sha256哈希)"`
	IP        string     `gorm:"size:50;comment:操作IP"`
	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabPurchaseLog struct {
	ID         uint   `gorm:"primarykey"`
	OrderID    uint   `gorm:"not null;index;comment:关联OrderID"`
	UserID     uint   `gorm:"not null;comment:操作人ID"`
	ActionType string `gorm:"size:50;not null;comment:操作类型: create-创建 update-修改 delete-删除 query-查询"`
	OldContent string `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent string `gorm:"type:text;comment:修改后内容(JSON)"`
	IP         string `gorm:"size:50;comment:操作IP"`
	Remark     string `gorm:"size:500;comment:备注/操作描述"`

	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime;comment:操作时间"`
}

func ApiPurchaseInit() {

	models.DB.AutoMigrate(&TabPurchaseOrder{})
	models.DB.AutoMigrate(&TabPurchaseCosts{})
	models.DB.AutoMigrate(&TabPurchaseFileBind{})
	models.DB.AutoMigrate(&TabPurchaseLog{})
	models.DB.AutoMigrate(&TabPurchaseCommit{})

	//先检查用户组有没有这个key  purchase
	purchaseUserGroup.Name = "purchase_admin"
	if models.DB.Where(&purchaseUserGroup).First(&purchaseUserGroup).Error == nil {
		updatePurchaseAdminsCash()
	} else {
		purchaseUserGroup.Type = "usergroup"
		models.DB.Create(&purchaseUserGroup)
	}

}

func ApiPurchase(r *gin.RouterGroup) {

	r.POST("/getorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromGetOrder struct {
			ID uint `json:"id"`
		}
		var from FromGetOrder
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var order TabPurchaseOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		// 查询关联费用
		var costs []TabPurchaseCosts
		models.DB.Where("order_id = ?", from.ID).Find(&costs)

		// 查询关联图片
		var binds []TabPurchaseFileBind
		models.DB.Where("order_id = ?", from.ID).Find(&binds)
		var fileIDs []uint
		for _, b := range binds {
			fileIDs = append(fileIDs, b.FileID)
		}
		var files []TabFileInfo_
		if len(fileIDs) > 0 {
			models.DB.Where("id IN ?", fileIDs).Find(&files)
		}

		// 查询状态变更记录
		var commits []TabPurchaseCommit
		models.DB.Where("order_id = ?", from.ID).Order("created_at DESC").Find(&commits)

		// 解析每条 commit 的 Photos JSON 字段为数组
		type CommitResponse struct {
			ID        uint      `json:"id"`
			OrderID   uint      `json:"orderId"`
			UserID    uint      `json:"userId"`
			Action    string    `json:"action"`
			Status    string    `json:"status"`
			OldStatus string    `json:"oldStatus"`
			Comment   string    `json:"comment"`
			IP        string    `json:"ip"`
			CreatedAt time.Time `json:"createdAt"`
			Photos    []string  `json:"photos"`
		}
		var commitResps []CommitResponse
		for _, c := range commits {
			// Status 优先用数据库字段；若为空（历史旧数据），从 Comment 备注中截取状态
			status := c.Status
			if status == "" {
				status = strings.TrimPrefix(c.Comment, "状态变更为: ")
				status = strings.TrimPrefix(status, "变更状态为: ")
				// 如果截取后跟原文一样，说明不是"状态变更为"格式，取原文作为展示
				if status == c.Comment {
					status = ""
				}
			}
			resp := CommitResponse{
				ID:        c.ID,
				OrderID:   c.OrderID,
				UserID:    c.UserID,
				Action:    c.Action,
				Status:    status,
				OldStatus: c.OldStatus,
				Comment:   c.Comment,
				IP:        c.IP,
				CreatedAt: time.Time{},
			}
			if c.CreatedAt != nil {
				resp.CreatedAt = *c.CreatedAt
			}
			if c.Photos != "" {
				json.Unmarshal([]byte(c.Photos), &resp.Photos)
			}
			commitResps = append(commitResps, resp)
		}

		// 判断当前用户是否可以修改
		canModify := canModifyPurchase(user.ID, order.UserID)

		ReturnJson(ctx, "apiOK", gin.H{
			"order":     order,
			"canModify": canModify,
			"costs":     costs,
			"photos":    files,
			"commits":   commitResps,
		})
	})

	// 更新订单状态（可附带评论）
	r.POST("/updatestatus", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromUpdateStatus struct {
			ID      uint     `json:"id"`
			Status  string   `json:"status" binding:"required"`
			Comment string   `json:"comment"`
			Photos  []string `json:"photos"` // 变更附带的图片哈希
		}
		var from FromUpdateStatus
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		// 校验图片哈希（不包含标点符号的哈希值）
		for _, hash := range from.Photos {
			if models.IsContainsSpecialChar(hash) {
				ReturnJson(ctx, "photo_hash_invalid", nil)
				return
			}
		}

		// 校验状态值
		validStatuses := map[string]bool{
			"pending":  true,
			"ordered":  true,
			"arrived":  true,
			"received": true,
			"lost":     true,
			"returned": true,
		}
		if !validStatuses[from.Status] {
			ReturnJson(ctx, "invalid_status", nil)
			return
		}

		var order TabPurchaseOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		oldStatus := order.OrderStatus
		if oldStatus == from.Status {
			ReturnJson(ctx, "status_no_change", nil)
			return
		}

		// 更新状态
		updates := map[string]interface{}{
			"order_status": from.Status,
		}
		if err := models.DB.Model(&order).Updates(updates).Error; err != nil {
			ReturnJson(ctx, "apiErr", nil)
			return
		}

		// 写状态变更 commit
		comment := from.Comment
		if comment == "" && len(from.Photos) == 0 {
			comment = "状态变更为: " + from.Status
		}
		photosJSON := ""
		if len(from.Photos) > 0 {
			if pj, err := json.Marshal(from.Photos); err == nil {
				photosJSON = string(pj)
			}
		}
		commit := TabPurchaseCommit{
			OrderID:   order.ID,
			UserID:    user.ID,
			Action:    "create_status",
			Status:    from.Status,
			OldStatus: oldStatus,
			Comment:   comment,
			Photos:    photosJSON,
			IP:        ctx.ClientIP(),
		}
		models.DB.Create(&commit)

		// 写操作日志
		newContent, _ := json.Marshal(map[string]string{
			"status":  from.Status,
			"comment": comment,
		})
		oldContent, _ := json.Marshal(map[string]string{
			"status": oldStatus,
		})
		tosqllog := TabPurchaseLog{
			UserID:     user.ID,
			OrderID:    order.ID,
			ActionType: "update_status",
			NewContent: string(newContent),
			OldContent: string(oldContent),
			IP:         ctx.ClientIP(),
			Remark:     comment,
		}
		models.DB.Create(&tosqllog)

		ReturnJson(ctx, "apiOK", nil)
	})

	r.POST("/getorders", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if isAuth {
			//fmt.Println(user)
			// DebugPrintJson(data)

			type From_purchase_getorders struct {
				Search  string
				Status  string
				Entries int
				Page    int
			}

			var jsondata From_purchase_getorders
			if err := decodeJSON(data, &jsondata); err == nil {
				//fmt.Println(jsondata)

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
					query := models.DB.Model(TabPurchaseOrder{})
					if jsondata.Search != "" {
						query = query.Where("title LIKE ?", "%"+jsondata.Search+"%")
					}
					if jsondata.Status != "" {
						query = query.Where("order_status = ?", jsondata.Status)
					}
					query.Count(&count)

					//读取条目
					var getorders []TabPurchaseOrder
					query.Order("created_at DESC").Offset(jsondata.Entries * (jsondata.Page - 1)).Limit(jsondata.Entries).Find(&getorders)

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
			//fmt.Println(data)
			if err := decodeJSON(data, &jsondata); err == nil {

				//jsonStr, _ := json.MarshalIndent(jsondata, "", "  ")
				//fmt.Println("转换后数据:\n", string(jsonStr))

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
						//fmt.Println("err4")
					}

				}

				//判断时间字符串是否合法
				// uptime, e := models.StringToTimePtr(jsondata.UpdateTime)
				// if e != nil {
				// 	is_data_ok = false
				// 	fmt.Println("err5")
				// }

				if is_data_ok {
					//校验通过
					//photos, _ := json.Marshal(jsondata.Photos)
					new_data := TabPurchaseOrder{
						UserID:      user.ID,
						Title:       jsondata.Title,
						Remark:      jsondata.Remark,
						Link:        jsondata.Link,
						Styles:      jsondata.Styles,
						OrderStatus: "pending", // 默认待处理
					}
					models.DB.Create(&new_data)

					for i := 0; i < len(jsondata.Costs); i++ {
						currencyType := jsondata.Costs[i].CurrencyType
						if currencyType <= 0 {
							currencyType = 1 // 默认 CNY
						}
						costType := jsondata.Costs[i].Type
						if costType <= 0 {
							costType = 1 // 默认单价
						}
						new_cost_data := TabPurchaseCosts{
							Price:        jsondata.Costs[i].Cost,
							Quantity:     jsondata.Costs[i].Int,
							UserID:       user.ID,
							OrderID:      new_data.ID,
							CurrencyType: currencyType,
							CostType:     costType,
						}
						models.DB.Create(&new_cost_data)
					}

					//绑定文件
					for i := 0; i < len(jsondata.Photos); i++ {
						findFile := TabFileInfo_{
							Sha256: jsondata.Photos[i],
							Type:   "image",
						}
						if models.DB.Where(&findFile).First(&findFile).Error == nil {
							bind := TabPurchaseFileBind{
								OrderID: new_data.ID,
								FileID:  findFile.ID,
							}
							models.DB.Create(&bind)
						}
					}

					// 写创建日志
					newContent, _ := json.Marshal(jsondata)
					tosqllog := TabPurchaseLog{
						UserID:     user.ID,
						OrderID:    new_data.ID,
						ActionType: "create",
						NewContent: string(newContent),
						OldContent: "",
						IP:         ctx.ClientIP(),
					}
					models.DB.Create(&tosqllog)

					// 写状态创建 commit
					commitLog := TabPurchaseCommit{
						OrderID:   new_data.ID,
						UserID:    user.ID,
						Action:    "create",
						Status:    "pending",
						OldStatus: "",
						Comment:   "订单创建",
						IP:        ctx.ClientIP(),
					}
					models.DB.Create(&commitLog)

					ReturnJson(ctx, "apiOK", nil)

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

	// 编辑订单基本信息
	r.POST("/updateorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromUpdateOrder struct {
			ID     uint       `json:"id"`
			Title  string     `json:"title"`
			Remark string     `json:"remark"`
			Link   string     `json:"link"`
			Styles string     `json:"styles"`
			Photos []string   `json:"photos"`
			Costs  []CostItem `json:"costs"`
		}
		var from FromUpdateOrder
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		if from.Title == "" {
			ReturnJson(ctx, "jsonErr_1", nil)
			return
		}

		// 校验图片哈希
		for _, hash := range from.Photos {
			if models.IsContainsSpecialChar(hash) {
				ReturnJson(ctx, "photo_hash_invalid", nil)
				return
			}
		}
		// 校验费用
		for _, c := range from.Costs {
			if c.Cost <= 0 || c.Int <= 0 {
				ReturnJson(ctx, "jsonErr_1", nil)
				return
			}
		}

		var order TabPurchaseOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		// 权限校验：只有创建者或管理员可以修改
		if !canModifyPurchase(user.ID, order.UserID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 记录旧数据
		oldContent, _ := json.Marshal(order)

		// 更新基本字段
		if err := models.DB.Model(&order).Updates(map[string]interface{}{
			"title":  from.Title,
			"remark": from.Remark,
			"link":   from.Link,
			"styles": from.Styles,
		}).Error; err != nil {
			ReturnJson(ctx, "apiErr", nil)
			return
		}

		// 重建费用明细：先删旧，再插新
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseCosts{})
		for _, c := range from.Costs {
			ct := c.CurrencyType
			if ct <= 0 {
				ct = 1
			}
			ctype := c.Type
			if ctype <= 0 {
				ctype = 1
			}
			models.DB.Create(&TabPurchaseCosts{
				Price:        c.Cost,
				Quantity:     c.Int,
				UserID:       user.ID,
				OrderID:      from.ID,
				CurrencyType: ct,
				CostType:     ctype,
			})
		}

		// 重建图片绑定：先删旧，再插新
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseFileBind{})
		for _, hash := range from.Photos {
			findFile := TabFileInfo_{Sha256: hash, Type: "image"}
			if models.DB.Where(&findFile).First(&findFile).Error == nil {
				models.DB.Create(&TabPurchaseFileBind{
					OrderID: from.ID,
					FileID:  findFile.ID,
				})
			}
		}

		// 写操作日志
		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabPurchaseLog{
			UserID:     user.ID,
			OrderID:    from.ID,
			ActionType: "update",
			OldContent: string(oldContent),
			NewContent: string(newContent),
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 删除订单
	r.POST("/deleteorder", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromDeleteOrder struct {
			ID uint `json:"id"`
		}
		var from FromDeleteOrder
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var order TabPurchaseOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		// 权限校验：只有创建者或管理员可以删除
		if !canModifyPurchase(user.ID, order.UserID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 关联删除（硬删，不保留）
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseCosts{})
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseFileBind{})
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseCommit{})
		models.DB.Where("order_id = ?", from.ID).Delete(&TabPurchaseLog{})

		// 软删除订单本身
		models.DB.Delete(&order)

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取订单数量统计
	r.POST("/getordercount", func(ctx *gin.Context) {
		isAuth, _, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type OrderCount struct {
			Pending  int64 `json:"pending"`  // 待处理
			Ordered  int64 `json:"ordered"`  // 已下单
			Arrived  int64 `json:"arrived"`  // 已到达
			Received int64 `json:"received"` // 已收件
			Lost     int64 `json:"lost"`     // 丢件
			Returned int64 `json:"returned"` // 退件
			Total    int64 `json:"total"`    // 总数
		}

		var count OrderCount
		models.DB.Model(&TabPurchaseOrder{}).Count(&count.Total)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "pending").Count(&count.Pending)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "ordered").Count(&count.Ordered)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "arrived").Count(&count.Arrived)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "received").Count(&count.Received)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "lost").Count(&count.Lost)
		models.DB.Model(&TabPurchaseOrder{}).Where("order_status = ?", "returned").Count(&count.Returned)

		ReturnJson(ctx, "apiOK", map[string]interface{}{
			"pending":  count.Pending,
			"ordered":  count.Ordered,
			"arrived": count.Arrived,
			"received": count.Received,
			"lost":    count.Lost,
			"returned": count.Returned,
			"total":   count.Total,
		})
	})

}
