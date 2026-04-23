package routers

import (
	"encoding/json"
	"ops/models"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------- 数据表结构 ----------

type TabWarehouseContainer struct {
	ID         uint       `gorm:"primaryKey" json:"ID"`
	Title      string     `gorm:"size:255;not null;comment:容器名" json:"Title"`
	Remark     string     `gorm:"type:text;comment:描述" json:"Remark"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime" json:"CreatedAt"`
	UpdatedAt  *time.Time `gorm:"type:datetime;autoUpdateTime" json:"UpdatedAt"`
	CreatorID  uint       `gorm:"not null;index;comment:创建者id" json:"CreatorID"`
	ParentID   *uint      `gorm:"index;comment:父容器id，nil=顶级" json:"ParentID"`
	ItemCount  int        `gorm:"default:0;comment:直接子物品数量" json:"ItemCount"`
	ChildCount int        `gorm:"default:0;comment:子容器数量" json:"ChildCount"`
}

type TabWarehouseItem struct {
	ID           uint       `gorm:"primaryKey" json:"ID"`
	Name         string     `gorm:"size:255;not null;comment:物品名" json:"Name"`
	SerialNumber string     `gorm:"size:255;comment:序列号" json:"SerialNumber"`
	Remark       string     `gorm:"type:text;comment:描述" json:"Remark"`
	Quantity     int        `gorm:"default:1;comment:数量" json:"Quantity"`
	CreatedAt    *time.Time `gorm:"type:datetime;autoCreateTime" json:"CreatedAt"`
	UpdatedAt    *time.Time `gorm:"type:datetime;autoUpdateTime" json:"UpdatedAt"`
	CreatorID    uint       `gorm:"not null;index;comment:创建者id" json:"CreatorID"`
	ContainerID  *uint      `gorm:"index;comment:所属容器id，nil=未入库" json:"ContainerID"`
}

type TabWarehouseContainerFileBind struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerID uint      `gorm:"not null;index;comment:关联容器id"`
	FileID      uint      `gorm:"not null;comment:关联文件id"`
	CreatorID   uint      `gorm:"not null;comment:上传人id"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseItemFileBind struct {
	ID        uint      `gorm:"primaryKey"`
	ItemID    uint      `gorm:"not null;index;comment:关联物品id"`
	FileID    uint      `gorm:"not null;comment:关联文件id"`
	CreatorID uint      `gorm:"not null;comment:上传人id"`
	CreatedAt time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseItemCommit struct {
	ID           uint      `gorm:"primaryKey" json:"ID"`
	ItemID       uint      `gorm:"not null;index;comment:关联物品id" json:"ItemID"`
	UserID       uint      `gorm:"not null;comment:操作人id" json:"UserID"`
	OldContainer *uint    `gorm:"index;comment:原容器id" json:"OldContainer"`
	NewContainer *uint    `gorm:"index;comment:新容器id" json:"NewContainer"`
	Remark       string    `gorm:"type:text;comment:备注" json:"Remark"`
	IP           string    `gorm:"size:50;comment:操作IP" json:"IP"`
	CreatedAt    time.Time `gorm:"type:datetime;autoCreateTime" json:"CreatedAt"`
}

type TabWarehouseLog struct {
	ID         uint      `gorm:"primaryKey"`
	EntityType string    `gorm:"size:50;not null;index;comment:操作对象类型"`
	EntityID   uint      `gorm:"not null;index;comment:操作对象id"`
	UserID     uint      `gorm:"not null;index;comment:操作人id"`
	ActionType string    `gorm:"size:50;not null;comment:操作类型: create update delete move query"`
	OldContent string    `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent string    `gorm:"type:text;comment:修改后内容(JSON)"`
	IP         string    `gorm:"size:50;comment:操作IP"`
	Remark     string    `gorm:"size:500;comment:备注"`
	CreatedAt  time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseItemWorkOrderBind struct {
	ID          uint      `gorm:"primaryKey"`
	ItemID      uint      `gorm:"not null;index;comment:关联物品id"`
	WorkOrderID uint      `gorm:"not null;index;comment:关联工单id"`
	Remark      string    `gorm:"size:500;comment:备注"`
	CreatorID   uint      `gorm:"not null;comment:绑定人id"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
}

var (
	warehouseUserGroup models.TabUserGroups_
	warehouseAdmins   []uint
)

// updateWarehouseAdminsCash 刷新仓库管理员缓存
func updateWarehouseAdminsCash() {
	warehouseAdmins = nil
	warehouseAdmins = append(warehouseAdmins, 1) // id=1 超级管理员
	var binds []models.TabUserGroupBinds_
	models.DB.Where("group_id = ?", warehouseUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(warehouseAdmins, item.UserID) {
			warehouseAdmins = append(warehouseAdmins, item.UserID)
		}
	}
}

// canModifyWarehouse 判断是否有权限修改/删除仓库资源（创建者或仓库管理员）
func canModifyWarehouse(userID, creatorUserID uint) bool {
	if slices.Contains(warehouseAdmins, userID) {
		return true
	}
	return userID == creatorUserID
}

// ---------- 初始化 ----------

func ApiWarehouseInit() {
	models.DB.AutoMigrate(
		&TabWarehouseContainer{},
		&TabWarehouseItem{},
		&TabWarehouseContainerFileBind{},
		&TabWarehouseItemFileBind{},
		&TabWarehouseItemCommit{},
		&TabWarehouseLog{},
		&TabWarehouseItemWorkOrderBind{},
	)

	warehouseUserGroup.Name = "warehouse_admin"
	if models.DB.Where(&warehouseUserGroup).First(&warehouseUserGroup).Error == nil {
		updateWarehouseAdminsCash()
	} else {
		warehouseUserGroup.Type = "usergroup"
		models.DB.Create(&warehouseUserGroup)
	}
}

// ---------- 路由注册 ----------

func ApiWarehouse(r *gin.RouterGroup) {

	// 新增容器
	r.POST("/add_container", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromAdd struct {
			Title    string `json:"title"`
			Remark   string `json:"remark"`
			ParentID *uint  `json:"parent_id"`
			Photos   []string `json:"photos"`
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

		// 检查嵌套层级不超过5层
		if from.ParentID != nil {
			var parent TabWarehouseContainer
			if err := models.DB.First(&parent, *from.ParentID).Error; err != nil {
				ReturnJson(ctx, "parent_not_found", nil)
				return
			}
			depth := 0
			curID := *from.ParentID
			for depth < 100 {
				var c TabWarehouseContainer
				if err := models.DB.First(&c, curID).Error; err != nil {
					break
				}
				if c.ParentID == nil {
					break
				}
				curID = *c.ParentID
				depth++
			}
			if depth >= 4 { // 新容器将落在第5层
				ReturnJson(ctx, "max_depth_exceeded", nil)
				return
			}
		}

		c := TabWarehouseContainer{
			Title:     from.Title,
			Remark:    from.Remark,
			CreatorID: user.ID,
			ParentID: from.ParentID,
		}
		models.DB.Create(&c)

		// 绑定图片
		for _, hash := range from.Photos {
			var findFile TabFileInfo_
			if models.DB.Where(&TabFileInfo_{Sha256: hash, Type: "image"}).First(&findFile).Error == nil {
				models.DB.Create(&TabWarehouseContainerFileBind{
					ContainerID: c.ID,
					FileID:     findFile.ID,
					CreatorID:  user.ID,
				})
			}
		}

		// 父容器的 ChildCount +1
		if from.ParentID != nil {
			models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *from.ParentID).Update("child_count", models.DB.Raw("child_count + 1"))
		}

		// 写操作日志
		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabWarehouseLog{
			EntityType: "container",
			EntityID:   c.ID,
			UserID:     user.ID,
			ActionType: "create",
			NewContent: string(newContent),
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", gin.H{"id": c.ID})
	})

	// 编辑容器
	r.POST("/update_container", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromUpdate struct {
			ID       uint    `json:"id"`
			Title    string  `json:"title"`
			Remark   string  `json:"remark"`
			Photos   []string `json:"photos"`
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

		var c TabWarehouseContainer
		if err := models.DB.Where("id = ?", from.ID).First(&c).Error; err != nil {
			ReturnJson(ctx, "container_not_found", nil)
			return
		}

		if !canModifyWarehouse(user.ID, c.CreatorID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		oldContent, _ := json.Marshal(c)
		models.DB.Model(&c).Updates(map[string]interface{}{
			"title":  from.Title,
			"remark": from.Remark,
		})

		// 重建图片绑定
		models.DB.Where("container_id = ?", from.ID).Delete(&TabWarehouseContainerFileBind{})
		for _, hash := range from.Photos {
			var findFile TabFileInfo_
			if models.DB.Where(&TabFileInfo_{Sha256: hash, Type: "image"}).First(&findFile).Error == nil {
				models.DB.Create(&TabWarehouseContainerFileBind{
					ContainerID: from.ID,
					FileID:      findFile.ID,
					CreatorID:   user.ID,
				})
			}
		}

		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabWarehouseLog{
			EntityType: "container",
			EntityID:   from.ID,
			UserID:     user.ID,
			ActionType: "update",
			OldContent: string(oldContent),
			NewContent: string(newContent),
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 删除容器
	r.POST("/delete_container", func(ctx *gin.Context) {
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

		var c TabWarehouseContainer
		if err := models.DB.Where("id = ?", from.ID).First(&c).Error; err != nil {
			ReturnJson(ctx, "container_not_found", nil)
			return
		}

		if !canModifyWarehouse(user.ID, c.CreatorID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 检查是否有子容器或物品
		if c.ChildCount > 0 || c.ItemCount > 0 {
			ReturnJson(ctx, "container_not_empty", nil)
			return
		}

		models.DB.Where("container_id = ?", from.ID).Delete(&TabWarehouseContainerFileBind{})

		// 父容器 ChildCount -1
		if c.ParentID != nil {
			models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *c.ParentID).Update("child_count", models.DB.Raw("child_count - 1"))
		}

		oldContent, _ := json.Marshal(c)
		models.DB.Create(&TabWarehouseLog{
			EntityType: "container",
			EntityID:   from.ID,
			UserID:     user.ID,
			ActionType: "delete",
			OldContent: string(oldContent),
			IP:         ctx.ClientIP(),
		})

		models.DB.Delete(&c)

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取容器列表
	r.POST("/list_container", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromList struct {
			Search    string `json:"search"`
			ParentID  *uint  `json:"parent_id"`
			AllLevels bool   `json:"all_levels"`
			Entries   int    `json:"entries"`
			Page      int    `json:"page"`
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
		query := models.DB.Model(&TabWarehouseContainer{})
		if from.Search != "" {
			query = query.Where("title LIKE ?", "%"+from.Search+"%")
		}
		if from.ParentID != nil {
			query = query.Where("parent_id = ?", *from.ParentID)
		} else if from.Search == "" && !from.AllLevels {
			// 无搜索时默认只显示顶级容器（all_levels=true 时返回所有层级）
			query = query.Where("parent_id IS NULL")
		}
		query.Count(&count)

		var containers []TabWarehouseContainer
		query.Order("created_at DESC").
			Offset(from.Entries * (from.Page - 1)).
			Limit(from.Entries).
			Find(&containers)

		canModifyContainers := make([]bool, len(containers))
		for i, c := range containers {
			canModifyContainers[i] = canModifyWarehouse(user.ID, c.CreatorID)
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"all_count":          count,
			"containers":         containers,
			"canModifyContainers": canModifyContainers,
		})
	})

	// 获取容器详情
	r.POST("/get_container", func(ctx *gin.Context) {
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

		var c TabWarehouseContainer
		if err := models.DB.Where("id = ?", from.ID).First(&c).Error; err != nil {
			ReturnJson(ctx, "container_not_found", nil)
			return
		}

		// 构建父容器链（从根到当前）
		type ParentItem struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		}
		parentChain := []ParentItem{}
		if c.ParentID != nil {
			curID := *c.ParentID
			visited := map[uint]bool{}
			for curID != 0 {
				if visited[curID] {
					break
				}
				visited[curID] = true
				var parent TabWarehouseContainer
				if err := models.DB.Select("id, title, parent_id").Where("id = ?", curID).First(&parent).Error; err != nil {
					break
				}
				parentChain = append([]ParentItem{{ID: parent.ID, Title: parent.Title}}, parentChain...)
				if parent.ParentID == nil {
					break
				}
				curID = *parent.ParentID
			}
		}

		// 计算当前容器深度（0=顶级，4=已达最大层级）
		depth := len(parentChain)

		// 关联图片
		var binds []TabWarehouseContainerFileBind
		models.DB.Where("container_id = ?", from.ID).Find(&binds)
		var fileIDs []uint
		for _, b := range binds {
			fileIDs = append(fileIDs, b.FileID)
		}
		var files []TabFileInfo_
		if len(fileIDs) > 0 {
			models.DB.Where("id IN ?", fileIDs).Find(&files)
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"container":        c,
			"photos":           files,
			"parent_chain":     parentChain,
			"depth":            depth,
			"canModifyContainer": canModifyWarehouse(user.ID, c.CreatorID),
		})
	})

	// 新增物品（查重逻辑：Name+SerialNumber相同则更新容器）
	r.POST("/add_item", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromAdd struct {
			Name         string  `json:"name"`
			SerialNumber string  `json:"serial_number"`
			Remark       string  `json:"remark"`
			Quantity     int     `json:"quantity"`
			ContainerID  *uint   `json:"container_id"`
			Photos       []string `json:"photos"`
		}
		var from FromAdd
		if err := decodeJSON(data, &from); err != nil || from.Name == "" {
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

		quantity := from.Quantity
		if quantity <= 0 {
			quantity = 1
		}

		// 查重：Name + SerialNumber 相同则更新容器
		var existingItem TabWarehouseItem
		exists := models.DB.Where("name = ? AND serial_number = ?", from.Name, from.SerialNumber).First(&existingItem).Error == nil

		var itemID uint
		var oldContainer *uint

		if exists {
			// 已有记录：更新容器
			oldContainer = existingItem.ContainerID

			// 同一容器无需操作，但仍然记录 commit
			if !ptrEqUint(oldContainer, from.ContainerID) {
				// 旧容器 ItemCount -1
				if oldContainer != nil {
					models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *oldContainer).Update("item_count", models.DB.Raw("item_count - 1"))
				}
				// 新容器 ItemCount +1
				if from.ContainerID != nil {
					models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *from.ContainerID).Update("item_count", models.DB.Raw("item_count + 1"))
				}
				// 更新物品容器
				existingItem.ContainerID = from.ContainerID
				models.DB.Save(&existingItem)
			}

			itemID = existingItem.ID

			// 写操作日志
			newContent, _ := json.Marshal(from)
			models.DB.Create(&TabWarehouseLog{
				EntityType: "item",
				EntityID:   itemID,
				UserID:     user.ID,
				ActionType: "update",
				OldContent: ptrStrUint(oldContainer),
				NewContent: string(newContent),
				IP:         ctx.ClientIP(),
			})
		} else {
			// 无记录：新建物品
			item := TabWarehouseItem{
				Name:         from.Name,
				SerialNumber: from.SerialNumber,
				Remark:       from.Remark,
				Quantity:     quantity,
				CreatorID:    user.ID,
				ContainerID:  from.ContainerID,
			}
			models.DB.Create(&item)
			itemID = item.ID

			// 绑定图片
			for _, hash := range from.Photos {
				var findFile TabFileInfo_
				if models.DB.Where(&TabFileInfo_{Sha256: hash, Type: "image"}).First(&findFile).Error == nil {
					models.DB.Create(&TabWarehouseItemFileBind{
						ItemID:    item.ID,
						FileID:    findFile.ID,
						CreatorID: user.ID,
					})
				}
			}

			// 所属容器的 ItemCount +1
			if from.ContainerID != nil {
				models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *from.ContainerID).Update("item_count", models.DB.Raw("item_count + 1"))
			}

			// 写操作日志
			newContent, _ := json.Marshal(from)
			models.DB.Create(&TabWarehouseLog{
				EntityType: "item",
				EntityID:   itemID,
				UserID:     user.ID,
				ActionType: "create",
				NewContent: string(newContent),
				IP:         ctx.ClientIP(),
			})
		}

		// 新增/更新时绑定图片
		for _, hash := range from.Photos {
			var findFile TabFileInfo_
			if models.DB.Where(&TabFileInfo_{Sha256: hash, Type: "image"}).First(&findFile).Error == nil {
				// 检查是否已绑定，避免重复
				var count int64
				models.DB.Model(&TabWarehouseItemFileBind{}).Where("item_id = ? AND file_id = ?", itemID, findFile.ID).Count(&count)
				if count == 0 {
					models.DB.Create(&TabWarehouseItemFileBind{
						ItemID:    itemID,
						FileID:    findFile.ID,
						CreatorID: user.ID,
					})
				}
			}
		}

		// 记录 commit（无论新建还是更新容器都记录）
		models.DB.Create(&TabWarehouseItemCommit{
			ItemID:       itemID,
			UserID:       user.ID,
			OldContainer: oldContainer,
			NewContainer: from.ContainerID,
			Remark:       from.Remark,
			IP:           ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", gin.H{"id": itemID, "updated": exists})
	})

	// 编辑物品
	r.POST("/update_item", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromUpdate struct {
			ID           uint    `json:"id"`
			Name         string  `json:"name"`
			SerialNumber string  `json:"serial_number"`
			Remark       string  `json:"remark"`
			Quantity     int     `json:"quantity"`
			Photos       []string `json:"photos"`
		}
		var from FromUpdate
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 || from.Name == "" {
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

		var item TabWarehouseItem
		if err := models.DB.Where("id = ?", from.ID).First(&item).Error; err != nil {
			ReturnJson(ctx, "item_not_found", nil)
			return
		}

		if !canModifyWarehouse(user.ID, item.CreatorID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		oldContent, _ := json.Marshal(item)
		models.DB.Model(&item).Updates(map[string]interface{}{
			"name":          from.Name,
			"serial_number": from.SerialNumber,
			"remark":        from.Remark,
			"quantity":      from.Quantity,
		})

		// 重建图片绑定
		models.DB.Where("item_id = ?", from.ID).Delete(&TabWarehouseItemFileBind{})
		for _, hash := range from.Photos {
			var findFile TabFileInfo_
			if models.DB.Where(&TabFileInfo_{Sha256: hash, Type: "image"}).First(&findFile).Error == nil {
				models.DB.Create(&TabWarehouseItemFileBind{
					ItemID:    from.ID,
					FileID:    findFile.ID,
					CreatorID: user.ID,
				})
			}
		}

		newContent, _ := json.Marshal(from)
		models.DB.Create(&TabWarehouseLog{
			EntityType: "item",
			EntityID:   from.ID,
			UserID:     user.ID,
			ActionType: "update",
			OldContent: string(oldContent),
			NewContent: string(newContent),
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 删除物品
	r.POST("/delete_item", func(ctx *gin.Context) {
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

		var item TabWarehouseItem
		if err := models.DB.Where("id = ?", from.ID).First(&item).Error; err != nil {
			ReturnJson(ctx, "item_not_found", nil)
			return
		}

		if !canModifyWarehouse(user.ID, item.CreatorID) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 所属容器 ItemCount -1
		if item.ContainerID != nil {
			models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *item.ContainerID).Update("item_count", models.DB.Raw("item_count - 1"))
		}

		// 删除关联
		models.DB.Where("item_id = ?", from.ID).Delete(&TabWarehouseItemFileBind{})
		models.DB.Where("item_id = ?", from.ID).Delete(&TabWarehouseItemCommit{})
		models.DB.Where("item_id = ?", from.ID).Delete(&TabWarehouseItemWorkOrderBind{})

		oldContent, _ := json.Marshal(item)
		models.DB.Create(&TabWarehouseLog{
			EntityType: "item",
			EntityID:   from.ID,
			UserID:     user.ID,
			ActionType: "delete",
			OldContent: string(oldContent),
			IP:         ctx.ClientIP(),
		})

		models.DB.Delete(&item)

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取物品列表
	r.POST("/list_item", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromList struct {
			Search      string `json:"search"`
			ContainerID *uint  `json:"container_id"`
			Entries     int    `json:"entries"`
			Page        int    `json:"page"`
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
		query := models.DB.Model(&TabWarehouseItem{})
		if from.Search != "" {
			query = query.Where("name LIKE ? OR serial_number LIKE ? OR remark LIKE ?",
				"%"+from.Search+"%", "%"+from.Search+"%", "%"+from.Search+"%")
		}
		if from.ContainerID != nil {
			query = query.Where("container_id = ?", *from.ContainerID)
		}
		query.Count(&count)

		var items []TabWarehouseItem
		query.Order("created_at DESC").
			Offset(from.Entries * (from.Page - 1)).
			Limit(from.Entries).
			Find(&items)

		canModifyItems := make([]bool, len(items))
		for i, item := range items {
			canModifyItems[i] = canModifyWarehouse(user.ID, item.CreatorID)
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"all_count":       count,
			"items":           items,
			"canModifyItems": canModifyItems,
		})
	})

	// 获取物品详情
	r.POST("/get_item", func(ctx *gin.Context) {
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

		var item TabWarehouseItem
		if err := models.DB.Where("id = ?", from.ID).First(&item).Error; err != nil {
			ReturnJson(ctx, "item_not_found", nil)
			return
		}

		// 关联图片
		var binds []TabWarehouseItemFileBind
		models.DB.Where("item_id = ?", from.ID).Find(&binds)
		var fileIDs []uint
		for _, b := range binds {
			fileIDs = append(fileIDs, b.FileID)
		}
		var files []TabFileInfo_
		if len(fileIDs) > 0 {
			models.DB.Where("id IN ?", fileIDs).Find(&files)
		}

		// 移动历史
		var commits []TabWarehouseItemCommit
		models.DB.Where("item_id = ?", from.ID).Order("created_at DESC").Find(&commits)

		// 关联工单
		var woBinds []TabWarehouseItemWorkOrderBind
		models.DB.Where("item_id = ?", from.ID).Find(&woBinds)

		type WOInfo struct {
			ID     uint   `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		}
		var workOrders []WOInfo
		for _, b := range woBinds {
			var wo TabWorkOrder
			if models.DB.Where("id = ?", b.WorkOrderID).First(&wo).Error == nil {
				workOrders = append(workOrders, WOInfo{ID: wo.ID, Title: wo.Title, Status: wo.CurrentStatus})
			}
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"item":          item,
			"photos":        files,
			"commits":       commits,
			"work_orders":   workOrders,
			"canModifyItem": canModifyWarehouse(user.ID, item.CreatorID),
		})
	})

	// 移动物品到其他容器
	r.POST("/move_item", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type FromMove struct {
			ItemID       uint    `json:"item_id"`
			NewContainer *uint   `json:"new_container"`
			Remark       string  `json:"remark"`
		}
		var from FromMove
		if err := decodeJSON(data, &from); err != nil || from.ItemID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var item TabWarehouseItem
		if err := models.DB.Where("id = ?", from.ItemID).First(&item).Error; err != nil {
			ReturnJson(ctx, "item_not_found", nil)
			return
		}

		oldContainer := item.ContainerID

		// 同一容器无需操作
		if ptrEqUint(oldContainer, from.NewContainer) {
			ReturnJson(ctx, "apiOK", nil)
			return
		}

		// 旧容器 ItemCount -1
		if oldContainer != nil {
			models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *oldContainer).Update("item_count", models.DB.Raw("item_count - 1"))
		}

		// 新容器 ItemCount +1
		if from.NewContainer != nil {
			models.DB.Model(&TabWarehouseContainer{}).Where("id = ?", *from.NewContainer).Update("item_count", models.DB.Raw("item_count + 1"))
		}

		// 更新物品容器
		item.ContainerID = from.NewContainer
		models.DB.Save(&item)

		// 记录移动日志
		models.DB.Create(&TabWarehouseItemCommit{
			ItemID:       from.ItemID,
			UserID:       user.ID,
			OldContainer: oldContainer,
			NewContainer: from.NewContainer,
			Remark:       from.Remark,
			IP:           ctx.ClientIP(),
		})

		// 写通用操作日志
		models.DB.Create(&TabWarehouseLog{
			EntityType: "item",
			EntityID:   from.ItemID,
			UserID:     user.ID,
			ActionType: "move",
			OldContent: ptrStrUint(oldContainer),
			NewContent: ptrStrUint(from.NewContainer),
			IP:         ctx.ClientIP(),
			Remark:     from.Remark,
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取仓库统计
	r.POST("/count", func(ctx *gin.Context) {
		isAuth, _, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		type WCount struct {
			ContainerTotal int64 `json:"container_total"`
			ItemTotal      int64 `json:"item_total"`
			UnstoredItems  int64 `json:"unstored_items"`
		}
		var count WCount
		models.DB.Model(&TabWarehouseContainer{}).Count(&count.ContainerTotal)
		models.DB.Model(&TabWarehouseItem{}).Count(&count.ItemTotal)
		models.DB.Model(&TabWarehouseItem{}).Where("container_id IS NULL").Count(&count.UnstoredItems)

		ReturnJson(ctx, "apiOK", gin.H{
			"container_total": count.ContainerTotal,
			"item_total":      count.ItemTotal,
			"unstored_items":   count.UnstoredItems,
		})
	})
}

// ---------- 辅助函数 ----------

func ptrEqUint(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func ptrStrUint(p *uint) string {
	if p == nil {
		return "nil"
	}
	return strconv.FormatUint(uint64(*p), 10)
}
