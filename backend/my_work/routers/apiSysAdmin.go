package routers

import (
	"ops/models"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// InitSysAdminRouter 初始化系统管理员路由
func ApiSysAdmin(r *gin.RouterGroup) {
	// 获取系统管理员列表（仅系统管理员可访问）
	r.POST("/sysadmins", func(ctx *gin.Context) {
		isAuth, user, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var redata map[string]interface{} = make(map[string]interface{})
		redata["sysAdmins"] = sysAdmins
		ReturnJson(ctx, "apiOK", redata)
	})

	// 获取用户列表（仅系统管理员可访问）
	r.POST("/users", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		// 解析分页和搜索参数
		var params struct {
			Page     int    `json:"page"`
			PageSize int    `json:"page_size"`
			Search   string `json:"search"`
		}
		if err := mapstructure.Decode(data, &params); err != nil {
			params.Page = 1
			params.PageSize = 20
		}
		if params.Page < 1 {
			params.Page = 1
		}
		if params.PageSize < 1 || params.PageSize > 100 {
			params.PageSize = 20
		}

		offset := (params.Page - 1) * params.PageSize

		// 构建查询
		var users []TabUser
		var total int64
		query := models.DB.Model(&TabUser{})

		// 搜索条件
		if params.Search != "" {
			search := "%" + params.Search + "%"
			query = query.Where("name LIKE ? OR email LIKE ?", search, search)
		}

		// 获取总数
		query.Count(&total)

		// 获取分页数据
		query.Order("id DESC").Offset(offset).Limit(params.PageSize).Find(&users)

		// 获取用户详细信息（包括头像）
		var userList []map[string]interface{}
		for _, u := range users {
			userInfo := GetUserInfoFromUserID(u.ID)
			userData := map[string]interface{}{
				"id":         u.ID,
				"name":       u.Name,
				"email":      u.Email,
				"type":       u.Type,
				"date":       u.Date,
				"username":   "",
				"avatarPath": "",
			}
			if userInfo != nil {
				userData["username"] = userInfo.Username
				userData["avatarPath"] = userInfo.AvatarPath
			}
			userList = append(userList, userData)
		}

		var redata map[string]interface{} = make(map[string]interface{})
		redata["users"] = userList
		redata["total"] = total
		redata["page"] = params.Page
		redata["page_size"] = params.PageSize

		ReturnJson(ctx, "apiOK", redata)
	})

	// 获取所有用户组列表（仅系统管理员可访问）
	r.POST("/groups", func(ctx *gin.Context) {
		isAuth, _, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		var groups []TabUserGroups
		models.DB.Order("id ASC").Find(&groups)

		var list []map[string]interface{}
		for _, g := range groups {
			// 统计成员数量
			var memberCount int64
			models.DB.Model(&TabUserGroupBinds{}).Where("group_id = ?", g.ID).Count(&memberCount)

			// 获取部分成员（最多5个）
			var binds []TabUserGroupBinds
			models.DB.Where("group_id = ?", g.ID).Limit(5).Find(&binds)
			var memberIDs []uint
			for _, b := range binds {
				memberIDs = append(memberIDs, b.UserID)
			}

			list = append(list, map[string]interface{}{
				"id":          g.ID,
				"name":        g.Name,
				"email":       g.Email,
				"type":        g.Type,
				"date":        g.Date,
				"memberCount": memberCount,
				"memberIDs":   memberIDs,
			})
		}

		var redata map[string]interface{} = make(map[string]interface{})
		redata["groups"] = list
		ReturnJson(ctx, "apiOK", redata)
	})

	// 获取用户组成员列表（仅系统管理员可访问）
	r.POST("/group_members", func(ctx *gin.Context) {
		isAuth, authUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(authUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params struct {
			GroupID  float64 `json:"group_id" mapstructure:"group_id"`
			Page     float64 `json:"page" mapstructure:"page"`
			PageSize float64 `json:"page_size" mapstructure:"page_size"`
		}
		if err := mapstructure.Decode(data, &params); err != nil {
			params.Page = 1
			params.PageSize = 20
		}
		if params.Page < 1 {
			params.Page = 1
		}
		if params.PageSize < 1 || params.PageSize > 100 {
			params.PageSize = 20
		}

		// 验证用户组是否存在
		var group TabUserGroups
		if models.DB.First(&group, params.GroupID).Error != nil {
			ReturnJson(ctx, "groupNotFound", nil)
			return
		}

		groupID := uint(params.GroupID)
		page := int(params.Page)
		pageSize := int(params.PageSize)
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		offset := (page - 1) * pageSize

		var binds []TabUserGroupBinds
		var total int64
		models.DB.Model(&TabUserGroupBinds{}).Where("group_id = ?", groupID).Count(&total)
		models.DB.Where("group_id = ?", groupID).Order("id ASC").Offset(offset).Limit(pageSize).Find(&binds)

		// 获取成员用户信息
		var members []map[string]interface{}
		for _, b := range binds {
			var u TabUser
			if models.DB.First(&u, b.UserID).Error == nil {
				userInfo := GetUserInfoFromUserID(u.ID)
				member := map[string]interface{}{
					"id":         u.ID,
					"name":       u.Name,
					"email":      u.Email,
					"type":       u.Type,
					"avatarPath": "",
					"username":   "",
				}
				if userInfo != nil {
					member["username"] = userInfo.Username
					member["avatarPath"] = userInfo.AvatarPath
				}
				members = append(members, member)
			}
		}

		var redata map[string]interface{} = make(map[string]interface{})
		redata["group_id"] = groupID
		redata["group_name"] = group.Name
		redata["members"] = members
		redata["total"] = total
		redata["page"] = page
		redata["page_size"] = pageSize
		ReturnJson(ctx, "apiOK", redata)
	})

	// 获取用户详细信息（仅系统管理员可访问）
	r.POST("/user_detail", func(ctx *gin.Context) {
		isAuth, authUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(authUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params struct {
			UserID uint `json:"user_id"`
		}
		if err := mapstructure.Decode(data, &params); err != nil || params.UserID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 获取用户基本信息
		var user TabUser
		if models.DB.First(&user, params.UserID).Error != nil {
			ReturnJson(ctx, "userNotFound", nil)
			return
		}

		// 获取用户扩展信息
		userInfo := GetUserInfoFromUserID(user.ID)

		// 构建返回数据
		redata := map[string]interface{}{
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"type":  user.Type,
				"date":  user.Date,
			},
			"userinfo": userInfo,
		}

		ReturnJson(ctx, "apiOK", redata)
	})

	// 重置用户密码（仅系统管理员可访问）
	r.POST("/reset_user_password", func(ctx *gin.Context) {
		isAuth, adminUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(adminUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params struct {
			UserID   float64 `json:"user_id" mapstructure:"user_id"`
			Password string  `json:"password" mapstructure:"password"`
		}
		if err := mapstructure.Decode(data, &params); err != nil || params.UserID == 0 || params.Password == "" {
			ReturnJson(ctx, "parameErr", map[string]interface{}{"decode_err": err != nil, "user_id": params.UserID, "pass_empty": params.Password == ""})
			return
		}

		// 查找目标用户
		var targetUser TabUser
		if models.DB.First(&targetUser, uint(params.UserID)).Error != nil {
			ReturnJson(ctx, "userNotFound", nil)
			return
		}

		// 生成新盐值并哈希密码
		newSalt := models.RandStr32()
		tempUser := TabUser{
			Pass: params.Password,
			Salt: newSalt,
		}
		HashUserPass(&tempUser)

		// 更新密码和盐值
		updates := TabUser{
			Pass: tempUser.Pass,
			Salt: newSalt,
		}
		if err := models.DB.Model(&targetUser).Updates(&updates).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}

		// 注销该用户的所有 cookie（强制重新登录）
		if err := models.DB.Where("user_id = ?", targetUser.ID).Delete(&TabUserCookie{}).Error; err != nil {
			// 删除 cookie 失败不影响密码修改结果，仅记录
			//fmt.Println("删除用户 cookie 失败:", err)
		}

		ReturnJson(ctx, "apiOK", nil)
	})

	// 添加用户组成员（仅系统管理员可访问）
	r.POST("/add_group_member", func(ctx *gin.Context) {
		isAuth, adminUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(adminUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params struct {
			GroupID float64 `json:"group_id" mapstructure:"group_id"`
			UserID  float64 `json:"user_id" mapstructure:"user_id"`
		}
		if err := mapstructure.Decode(data, &params); err != nil || params.GroupID == 0 || params.UserID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 验证用户组是否存在
		var group TabUserGroups
		if models.DB.First(&group, uint(params.GroupID)).Error != nil {
			ReturnJson(ctx, "groupNotFound", nil)
			return
		}

		// 验证用户是否存在
		var user TabUser
		if models.DB.First(&user, uint(params.UserID)).Error != nil {
			ReturnJson(ctx, "userNotFound", nil)
			return
		}

		// 检查绑定是否已存在
		var existingBind TabUserGroupBinds
		if models.DB.Where("group_id = ? AND user_id = ?", uint(params.GroupID), uint(params.UserID)).First(&existingBind).Error == nil {
			ReturnJson(ctx, "userAlreadyInGroup", nil)
			return
		}

		// 创建绑定
		newBind := TabUserGroupBinds{
			UserID:  uint(params.UserID),
			GroupID: uint(params.GroupID),
		}
		if err := models.DB.Create(&newBind).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}

		// 根据组名刷新对应的权限缓存
		switch group.Name {
		case "admins":
			updateSysAdminsCash()
		case "schedule_admin":
			ScheduleUpdateAdminsCash()
		case "purchase_admin":
			PurchaseUpdateAdminsCash()
		case "work_order_admin":
			WorkOrderUpdateAdminsCash()
		case "warehouse_admin":
			WarehouseUpdateAdminsCash()
		case "customer_admin":
			CustomerUpdateAdminsCash()
		}

		ReturnJson(ctx, "apiOK", nil)
	})

	// 移除用户组成员（仅系统管理员可访问）
	r.POST("/remove_group_member", func(ctx *gin.Context) {
		isAuth, adminUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(adminUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params struct {
			GroupID float64 `json:"group_id" mapstructure:"group_id"`
			UserID  float64 `json:"user_id" mapstructure:"user_id"`
		}
		if err := mapstructure.Decode(data, &params); err != nil || params.GroupID == 0 || params.UserID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 验证用户组是否存在
		var group TabUserGroups
		if models.DB.First(&group, uint(params.GroupID)).Error != nil {
			ReturnJson(ctx, "groupNotFound", nil)
			return
		}

		// 删除绑定
		if err := models.DB.Where("group_id = ? AND user_id = ?", uint(params.GroupID), uint(params.UserID)).Delete(&TabUserGroupBinds{}).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}

		// 根据组名刷新对应的权限缓存
		switch group.Name {
		case "admins":
			updateSysAdminsCash()
		case "schedule_admin":
			ScheduleUpdateAdminsCash()
		case "purchase_admin":
			PurchaseUpdateAdminsCash()
		case "work_order_admin":
			WorkOrderUpdateAdminsCash()
		case "warehouse_admin":
			WarehouseUpdateAdminsCash()
		case "customer_admin":
			CustomerUpdateAdminsCash()
		case "calendar_admin":
			CalendarUpdateAdminsCash()
		}

		ReturnJson(ctx, "apiOK", nil)
	})

	// 获取操作日志（仅系统管理员可访问）
	r.POST("/operation_logs", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		// 解析参数
		var params struct {
			Page     int    `json:"page"`
			PageSize int    `json:"page_size"`
			Module   string `json:"module"` // 模块: all/customer/purchase/schedule/warehouse/work_order
		}
		if err := mapstructure.Decode(data, &params); err != nil {
			params.Page = 1
			params.PageSize = 20
			params.Module = "all"
		}
		if params.Page < 1 {
			params.Page = 1
		}
		if params.PageSize < 1 || params.PageSize > 100 {
			params.PageSize = 20
		}

		type LogEntry struct {
			ID         uint       `json:"id"`
			Module     string     `json:"module"`      // 模块名称
			EntityID   uint       `json:"entity_id"`   // 关联实体ID
			UserID     uint       `json:"user_id"`     // 操作人ID
			ActionType string     `json:"action_type"` // 操作类型
			IP         string     `json:"ip"`
			Remark     string     `json:"remark"`
			CreatedAt  *time.Time `json:"created_at"`
		}

		var allLogs []LogEntry

		// 根据模块筛选查询
		if params.Module == "all" || params.Module == "customer" {
			var logs []TabCustomerLog
			query := models.DB.Model(&TabCustomerLog{})
			if params.Module == "customer" {
				query.Order("created_at DESC").Find(&logs)
			} else {
				query.Order("created_at DESC").Limit(1000).Find(&logs)
			}
			for _, log := range logs {
				allLogs = append(allLogs, LogEntry{
					ID:         log.ID,
					Module:     "customer",
					EntityID:   log.CustomerID,
					UserID:     log.UserID,
					ActionType: log.ActionType,
					IP:         log.IP,
					Remark:     log.Remark,
					CreatedAt:  log.CreatedAt,
				})
			}
		}

		if params.Module == "all" || params.Module == "purchase" {
			var logs []TabPurchaseLog
			query := models.DB.Model(&TabPurchaseLog{})
			if params.Module == "purchase" {
				query.Order("created_at DESC").Find(&logs)
			} else {
				query.Order("created_at DESC").Limit(1000).Find(&logs)
			}
			for _, log := range logs {
				allLogs = append(allLogs, LogEntry{
					ID:         log.ID,
					Module:     "purchase",
					EntityID:   log.OrderID,
					UserID:     log.UserID,
					ActionType: log.ActionType,
					IP:         log.IP,
					Remark:     log.Remark,
					CreatedAt:  log.CreatedAt,
				})
			}
		}

		if params.Module == "all" || params.Module == "schedule" {
			var logs []TabScheduleLog
			query := models.DB.Model(&TabScheduleLog{})
			if params.Module == "schedule" {
				query.Order("created_at DESC").Find(&logs)
			} else {
				query.Order("created_at DESC").Limit(1000).Find(&logs)
			}
			for _, log := range logs {
				allLogs = append(allLogs, LogEntry{
					ID:         log.ID,
					Module:     "schedule",
					EntityID:   log.ScheduleID,
					UserID:     log.UserID,
					ActionType: log.ActionType,
					IP:         log.IP,
					Remark:     log.Remark,
					CreatedAt:  log.CreatedAt,
				})
			}
		}

		if params.Module == "all" || params.Module == "warehouse" {
			var logs []TabWarehouseLog
			query := models.DB.Model(&TabWarehouseLog{})
			if params.Module == "warehouse" {
				query.Order("created_at DESC").Find(&logs)
			} else {
				query.Order("created_at DESC").Limit(1000).Find(&logs)
			}
			for _, log := range logs {
				allLogs = append(allLogs, LogEntry{
					ID:         log.ID,
					Module:     "warehouse",
					EntityID:   log.EntityID,
					UserID:     log.UserID,
					ActionType: log.ActionType,
					IP:         log.IP,
					Remark:     log.Remark,
					CreatedAt:  &log.CreatedAt,
				})
			}
		}

		if params.Module == "all" || params.Module == "work_order" {
			var logs []TabWorkOrderLog
			query := models.DB.Model(&TabWorkOrderLog{})
			if params.Module == "work_order" {
				query.Order("created_at DESC").Find(&logs)
			} else {
				query.Order("created_at DESC").Limit(1000).Find(&logs)
			}
			for _, log := range logs {
				allLogs = append(allLogs, LogEntry{
					ID:         log.ID,
					Module:     "work_order",
					EntityID:   log.WorkOrderID,
					UserID:     log.UserID,
					ActionType: log.ActionType,
					IP:         log.IP,
					Remark:     log.Remark,
					CreatedAt:  log.CreatedAt,
				})
			}
		}

		// 按时间倒序排序
		sort.Slice(allLogs, func(i, j int) bool {
			if allLogs[i].CreatedAt == nil || allLogs[j].CreatedAt == nil {
				return allLogs[i].ID > allLogs[j].ID
			}
			return allLogs[i].CreatedAt.After(*allLogs[j].CreatedAt)
		})

		total := len(allLogs)
		offset := (params.Page - 1) * params.PageSize
		end := offset + params.PageSize
		if offset > total {
			offset = total
		}
		if end > total {
			end = total
		}

		var pagedLogs []LogEntry
		if offset < total {
			pagedLogs = allLogs[offset:end]
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"logs":      pagedLogs,
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
		})
	})

	// 获取登录失败日志（仅系统管理员可访问）
	r.POST("/login_fail_logs", func(ctx *gin.Context) {
		isAuth, authUser, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		// 检查是否为系统管理员
		if !SysAdminCheck(authUser.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		// 解析分页和搜索参数
		var params struct {
			Page     int    `json:"page"`
			PageSize int    `json:"page_size"`
			Search   string `json:"search"`
		}
		if err := mapstructure.Decode(data, &params); err != nil {
			params.Page = 1
			params.PageSize = 20
		}
		if params.Page < 1 {
			params.Page = 1
		}
		if params.PageSize < 1 || params.PageSize > 100 {
			params.PageSize = 20
		}

		offset := (params.Page - 1) * params.PageSize

		// 构建查询
		var logs []TabUserLoginFailLog
		var total int64
		query := models.DB.Model(&TabUserLoginFailLog{})

		// 搜索条件（用户名或IP）
		if params.Search != "" {
			search := "%" + params.Search + "%"
			query = query.Where("username LIKE ? OR ip LIKE ?", search, search)
		}

		// 获取总数
		query.Count(&total)

		// 获取分页数据，按时间倒序
		query.Order("updated_at DESC").Offset(offset).Limit(params.PageSize).Find(&logs)

		// 构建返回数据
		var logList []map[string]interface{}
		for _, log := range logs {
			logData := map[string]interface{}{
				"id":         log.ID,
				"username":   log.Username,
				"user_id":    log.UserID,
				"ip":         log.IP,
				"user_agent": log.UserAgent,
				"reason":     log.Reason,
				"count":      log.Count,
				"created_at": log.CreatedAt,
				"updated_at": log.UpdatedAt,
			}
			logList = append(logList, logData)
		}

		var redata map[string]interface{} = make(map[string]interface{})
		redata["logs"] = logList
		redata["total"] = total
		redata["page"] = params.Page
		redata["page_size"] = params.PageSize

		ReturnJson(ctx, "apiOK", redata)
	})
}

// SysAdminCheck 检查当前用户是否为系统管理员
func SysAdminCheck(userID uint) bool {
	for _, adminID := range sysAdmins {
		if adminID == userID {
			return true
		}
	}
	return false
}

// RefreshSysAdmins 刷新系统管理员缓存（可从数据库重新加载）
func RefreshSysAdmins() {
	updateSysAdminsCash()
}
