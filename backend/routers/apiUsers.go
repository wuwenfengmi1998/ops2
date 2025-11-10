package routers

import (
	"fmt"
	"ops/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func ApiInit() {
	//用户模块初始化init
	fmt.Println("users init")

	//创建admin用户
	var user models.TabUser_
	user.Name = "admin"

	if models.DB.Where(&user).First(&user).Error == nil {

	} else {
		fmt.Println("用户不存在")
		user.Pass = models.HashUserPass("adminpassword")
		models.DB.Create(&user) // 传入指针
	}

	//创建admin group
	var usergroup models.TabUserGroups_
	usergroup.Name = "admins"
	if models.DB.Where(&usergroup).First(&usergroup).Error == nil {

	} else {
		fmt.Println("用户组不存在")
		models.DB.Create(&usergroup) // 传入指针
	}

	//创建用户与用户组绑定
	var usergroupbind models.TabUserGroupBinds_
	usergroupbind.UserID = user.ID
	usergroupbind.GroupID = usergroup.ID

	if models.DB.Where(&usergroupbind).First(&usergroupbind).Error == nil {

	} else {
		models.DB.Create(&usergroupbind) // 传入指针
	}

}

type From_user_add struct {
	Useremail string `json:"useremail"`
	Username  string `json:"username"`
	Userpass  string `json:"userpass"`
}

func ApiUser(r *gin.RouterGroup) {

	r.GET("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})
	r.POST("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})
	r.POST("/register", func(ctx *gin.Context) {
		//转换传进来的数据
		var jsonData From_user_add
		data, isHaveData := ctx.Get("data")

		if isHaveData {
			if err := mapstructure.Decode(data, &jsonData); err == nil {
				//转换字段
				newUser := models.TabUser_{
					Name:  jsonData.Username,
					Email: jsonData.Useremail,
					Pass:  jsonData.Userpass, // 实际应替换为哈希值
					Date:  time.Now(),
					// Date 字段无需赋值，数据库会自动填充默认值
				}
				if newUser.Name != "" && newUser.Pass != "" && newUser.Email != "" {
					//对用户的密码进行哈希替换
					newUser.Pass = models.HashUserPass(newUser.Pass)
					//用户名是唯一的，先读取是否有这个用户名
					var user models.TabUser_
					user.Name = newUser.Name

					if models.DB.Where(&user).First(&user).Error == nil {
						//fmt.Println("找到用户:", user.ID)
						ReturnJson(ctx, "userNameDup", nil)
					} else {
						//fmt.Println("用户不存在")
						models.DB.Create(&newUser) // 传入指针

						// //创建info
						// var user_info models.TabUserInfo_
						// user_info.AvatarPath = models.ConfigsUser.AvatarPath
						// user_info.UserID = newUser.ID
						// models.DB.Create(&user_info) // 传入指针

						ReturnJson(ctx, "apiOK", nil)
					}

				} else {
					ReturnJson(ctx, "jsonErr", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)

			}
		} else {
			ReturnJson(ctx, "postErr", nil)

		}

	})
}
