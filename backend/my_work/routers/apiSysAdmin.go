package routers

import (
	"ops/models"

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
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		var params struct {
			GroupID uint `json:"group_id"`
			Page     int  `json:"page"`
			PageSize int  `json:"page_size"`
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

		offset := (params.Page - 1) * params.PageSize

		var binds []TabUserGroupBinds
		var total int64
		models.DB.Model(&TabUserGroupBinds{}).Where("group_id = ?", params.GroupID).Count(&total)
		models.DB.Where("group_id = ?", params.GroupID).Order("id ASC").Offset(offset).Limit(params.PageSize).Find(&binds)

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
	redata["group_id"] = params.GroupID
	redata["group_name"] = group.Name
	redata["members"] = members
	redata["total"] = total
	redata["page"] = params.Page
	redata["page_size"] = params.PageSize
	ReturnJson(ctx, "apiOK", redata)
})

	// 获取登录失败日志（仅系统管理员可访问）
	r.POST("/login_fail_logs", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
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
