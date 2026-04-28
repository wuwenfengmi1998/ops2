package routers

import (
	"github.com/gin-gonic/gin"
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

	// TODO: 其他系统管理员接口可在此添加
	// 例如：用户管理、用户组管理、登录日志查询等
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
