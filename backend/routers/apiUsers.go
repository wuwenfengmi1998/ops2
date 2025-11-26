package routers

import (
	"fmt"
	"ops/models"
	"strconv"
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
		//fmt.Println("用户不存在")

		//对密码加盐
		user.Salt = models.RandStr32()
		user.Pass = "adminpassword"
		models.HashUserPass(&user)
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

type From_user_login struct {
	Username string `json:"username"`
	Userpass string `json:"userpass"`
	Remember bool   `json:"remember"`
}

type From_user_updateinfo struct {
	Username string `json:"username"`
	Remark   string `json:"remark"`
	Birthday string `json:"birthday"`
}

func AuthenticationAuthority(ctx *gin.Context) (bool, models.TabUser_, map[string]interface{}) {
	var user models.TabUser_

	data, cookieval := SeparateData(ctx)
	//fmt.Println("cookieis" + cookieval)
	if cookieval != "" {
		cookie := models.TabCookie_{
			Value: cookieval,
		}
		if models.DB.Where(&cookie).First(&cookie).Error == nil {
			//找到cookie，验证cookie有效性，以及更新cookie
			if models.CheckCookiesAndUpdate(&cookie) {
				//cookie有效
				//载入user
				user := models.TabUser_{
					ID: cookie.UserID,
				}
				models.DB.Where(&user).First(&user)

				return true, user, data

			} else {
				ReturnJson(ctx, "userCookieExpired", nil)
				return false, user, nil
			}

		} else {
			ReturnJson(ctx, "userCookieNotFund", nil)
			return false, user, nil
		}

	} else {
		ReturnJson(ctx, "userCookieError", nil)
		return false, user, nil
	}

	//return false, user
}

func ApiUser(r *gin.RouterGroup) {

	r.GET("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})
	r.POST("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})
	//更新用户info
	r.POST("/updateInfo", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var jsonData From_user_updateinfo

			if err := mapstructure.Decode(data, &jsonData); err == nil {
				// fmt.Println("updateinfo data is", jsonData)
				// fmt.Println(user)
				t, err := time.Parse("2006-01-02", jsonData.Birthday)
				if err == nil {
					var userinfo models.TabUserInfo_
					userinfo.UserID = user.ID

					var userinfoupdate models.TabUserInfo_
					userinfoupdate.UserID = user.ID
					userinfoupdate.CreatedAt = time.Now()
					userinfoupdate.Username = jsonData.Username
					userinfoupdate.Birthdate = t
					userinfoupdate.FirstName = jsonData.Remark

					//先查找是否有记录
					if models.DB.Where(&userinfo).First(&userinfo).Error == nil {
						//有记录，更新
						models.DB.Model(&userinfo).Updates(&userinfoupdate)
					} else {
						//无记录，创建
						models.DB.Create(&userinfoupdate) // 传入指针

					}

					ReturnJson(ctx, "apiOK", nil)

				} else {
					ReturnJson(ctx, "jsonErr", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		}

	})

	//通过cookie获取用户info
	r.POST("/getinfo", func(ctx *gin.Context) {
		isAuth, user, _ := AuthenticationAuthority(ctx)
		if isAuth {
			//载入用户info
			var userinfo models.TabUserInfo_
			userinfo.UserID = user.ID
			//fmt.Println(userInfo)
			var redata map[string]interface{} = make(map[string]interface{})
			if models.DB.Where(&userinfo).First(&userinfo).Error == nil {
				redata["userInfo"] = userinfo
			} else {
				redata["userInfo"] = nil
			}

			user.Pass = ""
			user.Salt = ""

			redata["user"] = user

			ReturnJson(ctx, "apiOK", redata)

		}
		// _, cookieval := SeparateData(ctx)
		// //fmt.Println("cookieis" + cookieval)
		// if cookieval != "" {
		// 	cookie := models.TabCookie_{
		// 		Value: cookieval,
		// 	}
		// 	if models.DB.Where(&cookie).First(&cookie).Error == nil {
		// 		//找到cookie，验证cookie有效性，以及更新cookie
		// 		if models.CheckCookiesAndUpdate(&cookie) {
		// 			//cookie有效
		// 			//返回最新cookie
		// 			redata := map[string]interface{}{
		// 				"cookie": cookie,
		// 			}
		// 			//载入用户info
		// 			userInfo := models.TabFileInfo_{
		// 				UserID: cookie.UserID,
		// 			}
		// 			if models.DB.Where(&userInfo).First(&userInfo).Error == nil {
		// 				redata["userInfo"] = userInfo
		// 			} else {
		// 				redata["userInfo"] = nil
		// 			}

		// 			//载入user
		// 			user := models.TabUser_{
		// 				ID: cookie.UserID,
		// 			}
		// 			models.DB.Where(&user).First(&user)
		// 			user.Pass = ""
		// 			user.Salt = ""

		// 			redata["user"] = user

		// 			ReturnJson(ctx, "apiOK", redata)

		// 		} else {
		// 			ReturnJson(ctx, "userCookieExpired", nil)
		// 		}

		// 	} else {
		// 		ReturnJson(ctx, "userCookieNotFund", nil)
		// 	}

		// } else {
		// 	ReturnJson(ctx, "userCookieError", nil)
		// }

	})
	//用户登陆
	r.POST("/login", func(ctx *gin.Context) {
		var loginuser From_user_login
		data, _ := SeparateData(ctx)
		if data != nil {
			if err := mapstructure.Decode(data, &loginuser); err == nil {
				if loginuser.Username != "" && loginuser.Userpass != "" {
					//传入的数据都ok，获取用户信息

					getuser := models.TabUser_{
						Name: loginuser.Username,
					}

					if models.DB.Where(&getuser).First(&getuser).Error == nil {
						//倒入数据
						user := models.TabUser_{
							Pass: loginuser.Userpass, //密码明文
							Salt: getuser.Salt,       //保存的盐制
						}
						//哈希密
						models.HashUserPass(&user)
						if user.Pass == getuser.Pass {
							//用户密码正确,生成cookie
							cookie := models.TabCookie_{
								UserID:    getuser.ID,
								Name:      "login",
								Value:     models.RandStr32(),
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								ExpiresAt: time.Now().Add(time.Duration(models.ConfigsUser.CookieTimeout) * time.Second), //计算过期时间,
								Remember:  loginuser.Remember,
							}
							models.DB.Create(&cookie) // 传入指针

							redata := map[string]interface{}{
								"cookie": cookie,
							}

							ReturnJson(ctx, "apiOK", redata)
						} else {
							ReturnJson(ctx, "userPassIncorrect", nil)
						}

					} else {
						//用户不存在
						ReturnJson(ctx, "userNameNoFund", nil)
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

	//用户注册
	r.POST("/register", func(ctx *gin.Context) {
		//转换传进来的数据
		var jsonData From_user_add

		data, _ := SeparateData(ctx)

		if data != nil {
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

					//用户名是唯一的，先读取是否有这个用户名
					var user models.TabUser_
					user.Name = newUser.Name

					if models.DB.Where(&user).First(&user).Error == nil {
						//fmt.Println("找到用户:", user.ID)
						ReturnJson(ctx, "userNameDup", nil)
					} else {
						//fmt.Println("用户不存在")

						//对密码加盐
						newUser.Salt = models.RandStr32()

						//对用户的密码进行哈希替换
						models.HashUserPass(&newUser)

						models.DB.Create(&newUser) // 传入指针

						//创建用户后写一个log

						models.LogAdd(ctx, "New user id:"+strconv.Itoa(int(newUser.ID)))

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
