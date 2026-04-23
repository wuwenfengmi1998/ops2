package routers

import (
	"encoding/json"
	"ops/models"
	parsefmt "fmt"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	workOrderUserGroup models.TabUserGroups_
	workOrderAdmins    []uint
)

// updateWorkOrderAdminsCash 刷新工单管理员缓存
func updateWorkOrderAdminsCash() {
	workOrderAdmins = nil
	workOrderAdmins = append(workOrderAdmins, 1) // id=1 超级管理员
	var binds []models.TabUserGroupBinds_
	models.DB.Where("group_id = ?", workOrderUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(workOrderAdmins, item.UserID) {
			workOrderAdmins = append(workOrderAdmins, item.UserID)
		}
	}
}

// canModifyWorkOrder 判断是否有权限修改/删除工单（创建者或管理员）
func canModifyWorkOrder(userID, creatorUserID uint) bool {
	if slices.Contains(workOrderAdmins, userID) {
		return true
	}
	return userID == creatorUserID
}

// ---------- 数据表结构 ----------

type TabWorkOrder struct {
	ID            uint           `gorm:"primarykey"`
	UserID        uint           `gorm:"not null;comment:创建人ID"`
	Title         string         `gorm:"size:200;not null;comment:工单标题"`
	Description   string         `gorm:"type:text;comment:问题描述"`
	CurrentStatus string         `gorm:"size:50;default:pending;comment:当前状态: pending-待处理 checked-已检查 parts_ordered-已下单零件 repaired-已维修 returned-已送还 unrepairable-无法维修"`
	CreatedAt     *time.Time     `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"type:datetime;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type TabWorkOrderFileBind struct {
	ID          uint       `gorm:"primarykey"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	FileID      uint       `gorm:"not null;comment:关联文件ID"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWorkOrderCommit struct {
	ID          uint       `gorm:"primarykey"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	UserID      uint       `gorm:"not null;comment:操作人ID"`
	Action      string     `gorm:"size:50;not null;comment:操作类型: create-创建 create_status-状态变更"`
	Status      string     `gorm:"size:50;comment:变更后的状态"`
	OldStatus   string     `gorm:"size:50;comment:变更前的状态"`
	Comment     string     `gorm:"type:text;comment:备注"`
	IP          string     `gorm:"size:50;comment:操作IP"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWorkOrderLog struct {
	ID          uint       `gorm:"primarykey"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	UserID      uint       `gorm:"not null;comment:操作人ID"`
	ActionType  string     `gorm:"size:50;not null;comment:操作类型: create update delete query"`
	OldContent  string     `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent  string     `gorm:"type:text;comment:修改后内容(JSON)"`
	IP          string     `gorm:"size:50;comment:操作IP"`
	Remark      string     `gorm:"size:500;comment:备注"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWorkOrderCommitFileBind struct {
	ID           uint       `gorm:"primarykey"`
	CommitID     uint       `gorm:"not null;index;comment:关联进度ID"`
	FileID       uint       `gorm:"not null;comment:关联文件ID"`
	WorkOrderID  uint       `gorm:"not null;index;comment:关联工单ID"`
	CreatedAt    *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabWorkOrderPurchaseOrderBind 工单与采购订单的关联表
type TabWorkOrderPurchaseOrderBind struct {
	ID                uint       `gorm:"primarykey"`
	WorkOrderID       uint       `gorm:"not null;index;comment:关联工单ID"`
	CommitID          uint       `gorm:"not null;index;comment:关联进度ID"`
	PurchaseOrderID   uint       `gorm:"not null;comment:关联采购订单ID"`
	CreatedAt         *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// PurchaseOrderInfo 采购订单简要信息
type PurchaseOrderInfo struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// ---------- 初始化 ----------

func ApiWorkOrderInit() {
	models.DB.AutoMigrate(&TabWorkOrder{})
	models.DB.AutoMigrate(&TabWorkOrderFileBind{})
	models.DB.AutoMigrate(&TabWorkOrderCommit{})
	models.DB.AutoMigrate(&TabWorkOrderLog{})
	models.DB.AutoMigrate(&TabWorkOrderCommitFileBind{})
	models.DB.AutoMigrate(&TabWorkOrderPurchaseOrderBind{})

	workOrderUserGroup.Name = "work_order_admin"
	if models.DB.Where(&workOrderUserGroup).First(&workOrderUserGroup).Error == nil {
		updateWorkOrderAdminsCash()
	} else {
		workOrderUserGroup.Type = "usergroup"
		models.DB.Create(&workOrderUserGroup)
	}
}

// ---------- 路由注册 ----------

func ApiWorkOrder(r *gin.RouterGroup) {

	// 新增工单
	r.POST("/add", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromAdd struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Photos      []string `json:"photos"`
		}
		var from FromAdd
		if err := decodeJSON(data, &from); err != nil || from.Title == "" {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		// 校验图片哈希
		for _, hash := range from.Photos {
			if models.IsContainsSpecialChar(hash) {
				ReturnJson(ctx, "photo_hash_invalid", nil)
				return
			}
		}

		order := TabWorkOrder{
			UserID:        user.ID,
			Title:         from.Title,
			Description:   from.Description,
			CurrentStatus: "pending",
		}
		models.DB.Create(&order)

		// 绑定图片
		for _, hash := range from.Photos {
			findFile := TabFileInfo_{Sha256: hash, Type: "image"}
			if models.DB.Where(&findFile).First(&findFile).Error == nil {
				models.DB.Create(&TabWorkOrderFileBind{
					WorkOrderID: order.ID,
					FileID:      findFile.ID,
				})
			}
		}

		// 写创建 commit
		models.DB.Create(&TabWorkOrderCommit{
			WorkOrderID: order.ID,
			UserID:      user.ID,
			Action:      "create",
			Status:      "pending",
			OldStatus:   "",
			Comment:     "工单创建",
			IP:          ctx.ClientIP(),
		})

		// 写操作日志
		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabWorkOrderLog{
			WorkOrderID: order.ID,
			UserID:      user.ID,
			ActionType:  "create",
			NewContent:  string(newContent),
			IP:          ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", gin.H{"id": order.ID})
	})

	// 编辑工单
	r.POST("/update", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromUpdate struct {
			ID          uint     `json:"id"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Photos      []string `json:"photos"`
		}
		var from FromUpdate
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 || from.Title == "" {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		// 校验图片哈希
		for _, hash := range from.Photos {
			if models.IsContainsSpecialChar(hash) {
				ReturnJson(ctx, "photo_hash_invalid", nil)
				return
			}
		}

		var order TabWorkOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		if !canModifyWorkOrder(user.ID, order.UserID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		oldContent, _ := json.Marshal(order)

		models.DB.Model(&order).Updates(map[string]interface{}{
			"title":       from.Title,
			"description": from.Description,
		})

		// 重建图片绑定
		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderFileBind{})
		for _, hash := range from.Photos {
			findFile := TabFileInfo_{Sha256: hash, Type: "image"}
			if models.DB.Where(&findFile).First(&findFile).Error == nil {
				models.DB.Create(&TabWorkOrderFileBind{
					WorkOrderID: from.ID,
					FileID:      findFile.ID,
				})
			}
		}

		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabWorkOrderLog{
			WorkOrderID: from.ID,
			UserID:      user.ID,
			ActionType:  "update",
			OldContent:  string(oldContent),
			NewContent:  string(newContent),
			IP:          ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取工单列表
	r.POST("/list", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromList struct {
			Search  string `json:"search"`
			Status  string `json:"status"`
			Entries int    `json:"entries"`
			Page    int    `json:"page"`
		}
		var from FromList
		if err := decodeJSON(data, &from); err != nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		if from.Entries <= 0 || from.Entries > 300 {
			from.Entries = 10
		}
		if from.Page <= 0 {
			from.Page = 1
		}

		var count int64
		query := models.DB.Model(&TabWorkOrder{})
		if from.Search != "" {
			query = query.Where("title LIKE ?", "%"+from.Search+"%")
		}
		if from.Status != "" {
			query = query.Where("current_status = ?", from.Status)
		}
		query.Count(&count)

		var orders []TabWorkOrder
		query.Order("created_at DESC").
			Offset(from.Entries * (from.Page - 1)).
			Limit(from.Entries).
			Find(&orders)

		ReturnJson(ctx, "apiOK", gin.H{
			"all_count":  count,
			"all_orders": orders,
		})
	})

	// 获取工单详情
	r.POST("/get", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromGet struct {
			ID uint `json:"id"`
		}
		var from FromGet
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var order TabWorkOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		// 关联图片
		var binds []TabWorkOrderFileBind
		models.DB.Where("work_order_id = ?", from.ID).Find(&binds)
		var fileIDs []uint
		for _, b := range binds {
			fileIDs = append(fileIDs, b.FileID)
		}
		var files []TabFileInfo_
		if len(fileIDs) > 0 {
			models.DB.Where("id IN ?", fileIDs).Find(&files)
		}

		// commits
		var commits []TabWorkOrderCommit
		models.DB.Where("work_order_id = ?", from.ID).Order("created_at DESC").Find(&commits)

		// 为每条 commit 附加图片和采购订单
		type CommitWithPhotos struct {
			TabWorkOrderCommit
			Photos         []TabFileInfo_      `json:"photos"`
			PurchaseOrders []PurchaseOrderInfo `json:"purchaseOrders"`
		}
		var commitsWithPhotos []CommitWithPhotos
		for _, c := range commits {
			item := CommitWithPhotos{TabWorkOrderCommit: c, Photos: []TabFileInfo_{}, PurchaseOrders: []PurchaseOrderInfo{}}

			// 附加图片
			var fileBinds []TabWorkOrderCommitFileBind
			models.DB.Where("commit_id = ?", c.ID).Find(&fileBinds)
			if len(fileBinds) > 0 {
				var fileIDs []uint
				for _, fb := range fileBinds {
					fileIDs = append(fileIDs, fb.FileID)
				}
				models.DB.Where("id IN ?", fileIDs).Find(&item.Photos)
			}

			// 附加采购订单
			var poBinds []TabWorkOrderPurchaseOrderBind
			models.DB.Where("commit_id = ?", c.ID).Find(&poBinds)
			if len(poBinds) > 0 {
				var poIDs []uint
				for _, pb := range poBinds {
					poIDs = append(poIDs, pb.PurchaseOrderID)
				}
				var pos []TabPurchaseOrder
				models.DB.Where("id IN ?", poIDs).Find(&pos)
				for _, po := range pos {
					item.PurchaseOrders = append(item.PurchaseOrders, PurchaseOrderInfo{
						ID:     po.ID,
						Title:  po.Title,
						Status: po.OrderStatus,
					})
				}
			}

			commitsWithPhotos = append(commitsWithPhotos, item)
		}

		canModify := canModifyWorkOrder(user.ID, order.UserID)
		// 所有登录用户都可以提交进度
		canCommit := true

		ReturnJson(ctx, "apiOK", gin.H{
			"order":     order,
			"canModify": canModify,
			"canCommit": canCommit,
			"photos":    files,
			"commits":   commitsWithPhotos,
		})
	})

	// 新增进度 commit（状态推进）
	r.POST("/commit", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromCommit struct {
			ID               uint     `json:"id"`
			Status           string   `json:"status"`
			Comment          string   `json:"comment"`
			Photos           []string `json:"photos"`
			PurchaseOrderIDs []uint   `json:"purchaseOrderIds"` // 关联的采购订单ID列表
		}
		var from FromCommit
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		validStatuses := map[string]bool{
			"pending":       true,
			"checked":       true,
			"parts_ordered": true,
			"repaired":      true,
			"returned":      true,
			"unrepairable":  true,
		}
		if !validStatuses[from.Status] {
			ReturnJson(ctx, "invalid_status", nil)
			return
		}

		// 校验图片哈希
		for _, hash := range from.Photos {
			if models.IsContainsSpecialChar(hash) {
				ReturnJson(ctx, "photo_hash_invalid", nil)
				return
			}
		}

		var order TabWorkOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		oldStatus := order.CurrentStatus
		models.DB.Model(&order).Update("current_status", from.Status)

		comment := from.Comment
		if comment == "" {
			comment = "状态变更为: " + from.Status
		}

		commit := TabWorkOrderCommit{
			WorkOrderID: order.ID,
			UserID:      user.ID,
			Action:      "create_status",
			Status:      from.Status,
			OldStatus:   oldStatus,
			Comment:     comment,
			IP:          ctx.ClientIP(),
		}
		models.DB.Create(&commit)

		// 绑定进度图片
		for _, hash := range from.Photos {
			findFile := TabFileInfo_{Sha256: hash, Type: "image"}
			if models.DB.Where(&findFile).First(&findFile).Error == nil {
				models.DB.Create(&TabWorkOrderCommitFileBind{
					CommitID:    commit.ID,
					WorkOrderID: order.ID,
					FileID:      findFile.ID,
				})
			}
		}

		// 绑定采购订单（去重）
		if len(from.PurchaseOrderIDs) > 0 {
			seen := make(map[uint]bool)
			for _, pid := range from.PurchaseOrderIDs {
				if !seen[pid] {
					seen[pid] = true
					models.DB.Create(&TabWorkOrderPurchaseOrderBind{
						WorkOrderID:     order.ID,
						CommitID:        commit.ID,
						PurchaseOrderID: pid,
					})
				}
			}
		}

		newContent, _ := json.Marshal(map[string]string{
			"status":  from.Status,
			"comment": comment,
		})
		oldContent, _ := json.Marshal(map[string]string{
			"status": oldStatus,
		})
		models.DB.Create(&TabWorkOrderLog{
			WorkOrderID: from.ID,
			UserID:      user.ID,
			ActionType:  "update_status",
			OldContent:  string(oldContent),
			NewContent:  string(newContent),
			IP:          ctx.ClientIP(),
			Remark:      comment,
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 删除工单
	r.POST("/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromDelete struct {
			ID uint `json:"id"`
		}
		var from FromDelete
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var order TabWorkOrder
		if err := models.DB.Where("id = ?", from.ID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		if !canModifyWorkOrder(user.ID, order.UserID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderFileBind{})
		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderCommitFileBind{})
		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderPurchaseOrderBind{})
		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderCommit{})
		models.DB.Where("work_order_id = ?", from.ID).Delete(&TabWorkOrderLog{})
		models.DB.Delete(&order)

		ReturnJson(ctx, "apiOK", nil)
	})

	// 删除进度
	r.POST("/delete_commit", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromDeleteCommit struct {
			WorkOrderID uint `json:"workOrderId"`
			CommitID    uint `json:"commitId"`
		}
		var from FromDeleteCommit
		if err := decodeJSON(data, &from); err != nil || from.WorkOrderID == 0 || from.CommitID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		// 获取工单信息
		var order TabWorkOrder
		if err := models.DB.Where("id = ?", from.WorkOrderID).First(&order).Error; err != nil {
			ReturnJson(ctx, "order_not_found", nil)
			return
		}

		// 获取进度信息
		var commit TabWorkOrderCommit
		if err := models.DB.Where("id = ? AND work_order_id = ?", from.CommitID, from.WorkOrderID).First(&commit).Error; err != nil {
			ReturnJson(ctx, "commit_not_found", nil)
			return
		}

		// 权限判断：工单创建者 或 进度创建者 或 管理员
		isOrderCreator := user.ID == order.UserID
		isCommitCreator := user.ID == commit.UserID
		isAdmin := slices.Contains(workOrderAdmins, user.ID)

		if !isOrderCreator && !isCommitCreator && !isAdmin {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 删除关联的采购订单绑定
		models.DB.Where("commit_id = ?", from.CommitID).Delete(&TabWorkOrderPurchaseOrderBind{})
		// 删除关联的图片
		models.DB.Where("commit_id = ?", from.CommitID).Delete(&TabWorkOrderCommitFileBind{})
		// 删除进度记录
		models.DB.Where("id = ?", from.CommitID).Delete(&commit)

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取工单数量统计
	r.POST("/count", func(ctx *gin.Context) {
		isAuth, _, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type WOCount struct {
			Pending      int64 `json:"pending"`
			Checked      int64 `json:"checked"`
			PartsOrdered int64 `json:"parts_ordered"`
			Repaired     int64 `json:"repaired"`
			Returned     int64 `json:"returned"`
			Unrepairable int64 `json:"unrepairable"`
			Total        int64 `json:"total"`
		}
		var count WOCount
		models.DB.Model(&TabWorkOrder{}).Count(&count.Total)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "pending").Count(&count.Pending)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "checked").Count(&count.Checked)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "parts_ordered").Count(&count.PartsOrdered)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "repaired").Count(&count.Repaired)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "returned").Count(&count.Returned)
		models.DB.Model(&TabWorkOrder{}).Where("current_status = ?", "unrepairable").Count(&count.Unrepairable)

		ReturnJson(ctx, "apiOK", gin.H{
			"pending":       count.Pending,
			"checked":       count.Checked,
			"parts_ordered": count.PartsOrdered,
			"repaired":      count.Repaired,
			"returned":      count.Returned,
			"unrepairable":  count.Unrepairable,
			"total":         count.Total,
		})
	})

	// 搜索采购订单（用于工单关联）
	r.POST("/search_purchase_orders", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromSearch struct {
			Search string `json:"search"`
			Limit  int    `json:"limit"`
		}
		var from FromSearch
		if err := decodeJSON(data, &from); err != nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		if from.Limit <= 0 || from.Limit > 20 {
			from.Limit = 5
		}

		query := models.DB.Model(&TabPurchaseOrder{})

		// 如果搜索词为空，返回最新的 N 条
		if from.Search != "" {
			// 尝试精确匹配 ID
			var id uint
			if _, err := parsefmt.Sscanf(from.Search, "%d", &id); err == nil && id > 0 {
				query = query.Where("id = ?", id)
			} else {
				// 模糊匹配标题或备注
				query = query.Where("title LIKE ? OR remark LIKE ?",
					"%"+from.Search+"%", "%"+from.Search+"%")
			}
		}

		var orders []TabPurchaseOrder
		query.Order("created_at DESC").Limit(from.Limit).Find(&orders)

		type OrderInfo struct {
			ID     uint   `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		}
		var result []OrderInfo
		for _, o := range orders {
			result = append(result, OrderInfo{
				ID:     o.ID,
				Title:  o.Title,
				Status: o.OrderStatus,
			})
		}

		ReturnJson(ctx, "apiOK", gin.H{"orders": result})
	})
}
