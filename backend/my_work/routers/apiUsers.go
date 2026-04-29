package routers

import (
	"errors"
	"fmt"
	"ops/models"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type TabUser struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"` // 自增主键
	Name  string    `gorm:"size:100;uniqueIndex"`     // 唯一约束索引
	Email string    `gorm:"size:255;index"`           // 字符串长度限制100 索引
	Pass  string    `gorm:"size:128"`                 // 建议存储哈希后的密码
	Type  string    `gorm:"size:64;default:user"`     //
	Salt  string    `gorm:"size:64;"`
	Date  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserGroups struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`                // 自增主键
	Name  string    `gorm:"size:100;uniqueIndex"`                    // 唯一约束索引
	Email string    `gorm:"size:255;index"`                          // 字符串长度限制100 索引
	Type  string    `gorm:"size:64;default:usergroup"`               //
	Date  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserGroupBinds struct {
	ID      uint      `gorm:"primaryKey;autoIncrement"` // 自增主键
	UserID  uint      `gorm:"index"`
	GroupID uint      `gorm:"index"`
	Date    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserInfo struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null;uniqueIndex"`
	FirstName  string    `gorm:"size:50;null"`
	Username   string    `gorm:"size:30;null"`
	Birthdate  time.Time `gorm:"type:datetime;null"`
	Gender     string    `gorm:"type:char(1);check:gender IN ('M', 'F', 'U');default:'U'"`
	AvatarPath string    `gorm:"size:255"`
	Region     string    `gorm:"size:50"`
	Language   string    `gorm:"size:10;default:'zh-CN'"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:created_at"`
}

// var def_user_info = User_info{
// 	ID:0,
// 	UserID:0,
// }

type TabUserCookie struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	Name      string    `gorm:"size:255;not null;index"`
	Value     string    `gorm:"size:255;not null;index"`
	ExpiresAt time.Time `gorm:"type:datetime;index"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;index;not null;default:CURRENT_TIMESTAMP"`
	Remember  bool      `gorm:"default:false"`
}

// TabUserLoginFailLog 用户登录失败日志表
type TabUserLoginFailLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"size:100;index"`                                // 尝试登录的用户名
	UserID    uint      `gorm:"index;default:0"`                               // 用户ID（如果用户存在）
	IP        string    `gorm:"size:64"`                                       // 登录IP
	UserAgent string    `gorm:"size:512"`                                      // User-Agent
	Reason    string    `gorm:"size:64"`                                       // 失败原因：password_error / user_not_found
	Count     int       `gorm:"default:1"`                                     // 连续失败次数
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`       // 首次失败时间
	UpdatedAt time.Time `gorm:"type:datetime;index;default:CURRENT_TIMESTAMP"` // 最后失败时间
}

var (
	sysUserGroup TabUserGroups
	sysAdmins    []uint
)

func updateSysAdminsCash() {
	// 查询 admins 用户组的 ID
	var adminGroup TabUserGroups
	if models.DB.Where("name = ?", "admins").First(&adminGroup).Error != nil {
		// admins 组不存在，清空缓存
		sysAdmins = []uint{}
		return
	}

	// 查询 admins 组的所有成员
	var binds []TabUserGroupBinds
	if models.DB.Where("group_id = ?", adminGroup.ID).Find(&binds).Error != nil {
		sysAdmins = []uint{}
		return
	}

	// 更新缓存
	newAdmins := make([]uint, 0, len(binds))
	for _, bind := range binds {
		newAdmins = append(newAdmins, bind.UserID)
	}
	sysAdmins = newAdmins
}

// recordLoginFail 记录登录失败日志，更新连续失败次数
func recordLoginFail(ctx *gin.Context, username string, userID uint, reason string) {
	// 获取客户端IP和UserAgent
	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	// 查找是否有该用户最近的失败记录（24小时内）
	var existingLog TabUserLoginFailLog
	err := models.DB.Where("username = ? AND reason = ? AND updated_at > ?",
		username, reason, time.Now().Add(-24*time.Hour)).
		Order("id DESC").First(&existingLog).Error

	if err == nil {
		// 存在最近失败的记录，更新次数
		models.DB.Model(&existingLog).Updates(map[string]interface{}{
			"count":      existingLog.Count + 1,
			"ip":         ip,
			"user_agent": userAgent,
			"updated_at": time.Now(),
		})
	} else {
		// 不存在，创建新记录
		newLog := TabUserLoginFailLog{
			Username:  username,
			UserID:    userID,
			IP:        ip,
			UserAgent: userAgent,
			Reason:    reason,
			Count:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		models.DB.Create(&newLog)
	}
}

func HashUserPass(user *TabUser) {
	switch models.ConfigsUser.PassHashType {
	case "text":
		break
	case "md5":
		user.Pass = models.Md5Str(user.Pass)

	case "md5salt":
		if user.Salt == "" {
			user.Salt = models.RandStr32()
		}
		user.Pass = models.Md5Str(models.Md5Str(user.Pass) + user.Salt)

	}

}

func CheckCookiesAndUpdate(cookie *TabUserCookie) bool {
	if !models.IsExpired(cookie.ExpiresAt) {
		if cookie.Remember {
			cookiewhere := TabUserCookie{
				ID: cookie.ID,
			}
			cookieupdata := TabUserCookie{
				UpdatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Duration(models.ConfigsUser.CookieTimeout) * time.Second),
			}
			models.DB.Where(&cookiewhere).Updates(&cookieupdata)

		}
		return true
	} else {
		//以过期
		return false
	}
	//return false
}

func ApiUserInit() {
	//用户模块初始化init
	fmt.Println("users init")

	// 自动创建表结构
	models.DB.AutoMigrate(&TabUser{})

	models.DB.AutoMigrate(&TabUserGroups{})

	models.DB.AutoMigrate(&TabUserGroupBinds{})

	models.DB.AutoMigrate(&TabUserInfo{})

	models.DB.AutoMigrate(&TabUserCookie{})

	models.DB.AutoMigrate(&TabUserLoginFailLog{})

	//创建admin用户
	var user TabUser
	user.Name = "admin"

	if models.DB.Where(&user).First(&user).Error == nil {

	} else {
		//fmt.Println("用户不存在")

		//对密码加盐
		user.Salt = models.RandStr32()
		user.Pass = "adminpassword"
		HashUserPass(&user)
		models.DB.Create(&user) // 传入指针
	}

	//创建admin group
	var usergroup TabUserGroups
	usergroup.Name = "admins"
	if models.DB.Where(&usergroup).First(&usergroup).Error == nil {

	} else {
		//fmt.Println("用户组不存在")
		models.DB.Create(&usergroup) // 传入指针
	}

	//创建用户与用户组绑定
	var usergroupbind TabUserGroupBinds
	usergroupbind.UserID = user.ID
	usergroupbind.GroupID = usergroup.ID

	if models.DB.Where(&usergroupbind).First(&usergroupbind).Error == nil {

	} else {
		models.DB.Create(&usergroupbind) // 传入指针
	}

	//更新系统管理员列表
	updateSysAdminsCash()
}

type From_user_add struct {
	Useremail string `json:"useremail"`
	Username  string `json:"username"`
	Userpass  string `json:"userpass"`
}

type From_user_login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type From_user_updateinfo struct {
	Username string `json:"username"`
	Remark   string `json:"remark"`
	Birthday string `json:"birthday"`
}

type From_user_changeemail struct {
	Newemail string `json:"newemail"`
}

type From_user_changepass struct {
	Oldpass string `json:"oldpass"`
	Newpass string `json:"newpass"`
}

func AuthenticationAuthorityFromCookie(c string) (*TabUser, error) {

	if c != "" {
		cookie := TabUserCookie{
			Value: c,
		}
		if models.DB.Where(&cookie).First(&cookie).Error == nil {
			//找到cookie，验证cookie有效性，以及更新cookie
			if CheckCookiesAndUpdate(&cookie) {
				//cookie有效
				//载入user
				user := TabUser{
					ID: cookie.UserID,
				}
				models.DB.Where(&user).First(&user)
				return &user, nil
			} else {
				return nil, errors.New("cookie  过期")
			}
		} else {
			return nil, errors.New("cookie Not Fund")
		}
	} else {
		return nil, errors.New("cookie 参数错误")
	}
}

func GetUserInfoFromUserID(userID uint) *TabUserInfo {
	//通过id获取用户info

	if userID <= 0 {
		return nil
	}

	//先查询用户是否存在
	var user TabUser
	user.ID = userID

	if models.DB.Where(&user).First(&user).Error == nil {
		var userinfo TabUserInfo
		userinfo.UserID = user.ID
		if models.DB.Where(&userinfo).First(&userinfo).Error == nil {
			return &userinfo
		} else {
			//无记录，创建一条
			userinfo.Username = user.Name
			userinfo.FirstName = user.Email
			userinfo.Birthdate = (time.Now())
			models.DB.Create(&userinfo)
			return &userinfo
		}
	}

	return nil

}

func AuthenticationAuthority(ctx *gin.Context) (bool, TabUser, map[string]interface{}) {

	data, cookieval := SeparateData(ctx)
	//fmt.Println("cookieis" + cookieval)
	var user TabUser
	if cookieval != "" {
		user_, error := AuthenticationAuthorityFromCookie(cookieval)
		if error == nil {
			user = *user_
			return true, user, data
		} else {
			return false, user, nil
		}

	} else {
		ReturnJson(ctx, "userCookieError", nil)
		return false, user, nil
	}

}

func ApiUser(r *gin.RouterGroup) {

	r.GET("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})
	r.POST("/test", func(ctx *gin.Context) {
		ReturnJson(ctx, "apiOK", nil)
	})

	//get获取用户info
	r.GET("/getuserinfo/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		var redata map[string]interface{} = make(map[string]interface{})
		if err == nil {
			userinfo := GetUserInfoFromUserID(uint(id))
			if userinfo != nil {
				redata["userinfo"] = *userinfo
			}

		}

		ReturnJson(ctx, "apiOK", redata)

	})

	//修改用户密码
	r.POST("/changePassword", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var jsonData From_user_changepass
			if err := mapstructure.Decode(data, &jsonData); err == nil {
				//验证旧密码
				//fmt.Println(user)
				//转换旧密码
				olduser := TabUser{
					Pass: jsonData.Oldpass,
					Salt: user.Salt,
				}
				HashUserPass(&olduser)
				if olduser.Pass == user.Pass {
					//旧密码正确，更新新密码
					var userupdate TabUser
					userupdate.Pass = jsonData.Newpass
					userupdate.Salt = models.RandStr32()
					HashUserPass(&userupdate)
					models.DB.Model(&user).Updates(&userupdate)
					ReturnJson(ctx, "apiOK", nil)
				} else {
					//旧密码错误
					ReturnJson(ctx, "userPassIncorrect", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}

		}
	})

	//更新用户邮箱
	r.POST("/changeEmail", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var jsonData From_user_changeemail
			if err := mapstructure.Decode(data, &jsonData); err == nil {
				//判断新邮箱格式
				if models.IsEmailValid(jsonData.Newemail) {
					var userupdate TabUser
					userupdate.Email = jsonData.Newemail
					models.DB.Model(&user).Updates(&userupdate)
					ReturnJson(ctx, "apiOK", nil)
				} else {
					ReturnJson(ctx, "userEmailFormatError", nil)

				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		}
	})

	//修改用户头像
	r.POST("/updateAvatar", func(ctx *gin.Context) {
		cookie := ctx.PostForm("cookie")
		user, err := AuthenticationAuthorityFromCookie(cookie)
		if err == nil {
			file, err := ctx.FormFile("file")
			if err == nil {
				if file.Filename != "" {
					//限制文件大小
					if file.Size > 512 {
						//头像裁剪过限制1M应该差不多
						if file.Size < 1048576 {

							//判断mime
							mimeType, err := models.GetFileMime(file)
							if err == nil {

								file_extname := models.ConfigsFile.AllowImageMime[mimeType]
								if file_extname != "" {

									//haxi文件

									file_hashi_name, err := models.SHA256HashFile(file)
									if err == nil {

										dst := path.Join(models.ConfigsFile.Pahts["avatar"], file_hashi_name+file_extname)

										var is_save_ok = false
										//判断文件是否存在避免重复保存
										if models.FileExists(dst) {
											//fmt.Println("文件存在")
											is_save_ok = true
											ReturnJson(ctx, "apiOK", nil)
										} else {
											//fmt.Println("文件no存在")
											ferr := ctx.SaveUploadedFile(file, dst)
											if ferr == nil {
												//文件保存成功
												//fmt.Print("save_ok")
												is_save_ok = true
												ReturnJson(ctx, "apiOK", nil)
											} else {
												//fmt.Print(ferr)
												ReturnJson(ctx, "postErr", nil)
											}

										}
										if is_save_ok {
											//修改数据库内容
											var user_info_fund TabUserInfo
											user_info_fund.UserID = user.ID

											var user_update_avatar TabUserInfo
											user_update_avatar.AvatarPath = file_hashi_name + file_extname

											//先查找是否有记录
											if models.DB.Where(&user_info_fund).First(&user_info_fund).Error == nil {
												//有记录，更新
												models.DB.Model(&user_info_fund).Updates(&user_update_avatar)
											} else {
												//无记录，创建
												user_update_avatar.UserID = user.ID
												models.DB.Create(&user_update_avatar)
											}

										}

									} else {
										ReturnJson(ctx, "postErr", nil)
									}

								} else {
									ReturnJson(ctx, "file_mime_err", nil)
								}

							} else {
								ReturnJson(ctx, "postErr", nil)
							}

						} else {
							ReturnJson(ctx, "file_size_err", nil)
						}
					} else {
						ReturnJson(ctx, "file_size_err", nil)
					}
				} else {
					ReturnJson(ctx, "file_name_err", nil)
				}
			} else {
				ReturnJson(ctx, "file_get_err", nil)
			}

		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

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
					var userinfo TabUserInfo
					userinfo.UserID = user.ID

					var userinfoupdate TabUserInfo
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
			var redata map[string]interface{} = make(map[string]interface{})

			info := GetUserInfoFromUserID(user.ID)
			redata["userInfo"] = *info

			user.Pass = ""
			user.Salt = ""
			redata["user"] = user

			// 只返回当前用户是否为系统管理员，不暴露完整列表
			isSysAdmin := false
			for _, adminID := range sysAdmins {
				if adminID == user.ID {
					isSysAdmin = true
					break
				}
			}
			redata["isSysAdmin"] = isSysAdmin

			ReturnJson(ctx, "apiOK", redata)

		}
	})

	//用户登陆
	r.POST("/login", func(ctx *gin.Context) {
		var loginuser From_user_login
		data, _ := SeparateData(ctx)
		if data != nil {
			if err := mapstructure.Decode(data, &loginuser); err == nil {
				if loginuser.Username != "" && loginuser.Password != "" {
					//传入的数据都ok，获取用户信息

					getuser := TabUser{
						Name: loginuser.Username,
					}

					if models.DB.Where(&getuser).First(&getuser).Error == nil {
						//倒入数据
						user := TabUser{
							Pass: loginuser.Password, //密码明文
							Salt: getuser.Salt,       //保存的盐制
						}
						//哈希密
						HashUserPass(&user)
						if user.Pass == getuser.Pass {
							//用户密码正确,生成cookie
							cookie := TabUserCookie{
								UserID:    getuser.ID,
								Name:      "login",
								Value:     models.RandStr32(),
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								ExpiresAt: time.Now().Add(time.Duration(models.ConfigsUser.CookieTimeout) * time.Second), //计算过期时间,
								Remember:  loginuser.Remember,
							}
							models.DB.Create(&cookie) // 传入指针

							// 登录成功，清除该用户的失败记录
							models.DB.Where("username = ?", loginuser.Username).Delete(&TabUserLoginFailLog{})

							redata := map[string]interface{}{
								"cookie": cookie,
							}
							//登录成功，记录日志
							//recordLoginFail(ctx, loginuser.Username, getuser.ID, "logined")
							ReturnJson(ctx, "apiOK", redata)
						} else {
							// 密码错误，记录失败日志
							recordLoginFail(ctx, loginuser.Username, getuser.ID, "password_error")
							ReturnJson(ctx, "userPassIncorrect", nil)
						}

					} else {
						//用户不存在，记录失败日志
						recordLoginFail(ctx, loginuser.Username, 0, "user_not_found")
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
				newUser := TabUser{
					Name:  jsonData.Username,
					Email: jsonData.Useremail,
					Pass:  jsonData.Userpass, // 实际应替换为哈希值
					Date:  time.Now(),
					// Date 字段无需赋值，数据库会自动填充默认值
				}
				if newUser.Name != "" && newUser.Pass != "" && newUser.Email != "" {

					//用户名是唯一的，先读取是否有这个用户名
					var user TabUser
					user.Name = newUser.Name

					if models.DB.Where(&user).First(&user).Error == nil {
						//fmt.Println("找到用户:", user.ID)
						ReturnJson(ctx, "userNameDup", nil)
					} else {
						//fmt.Println("用户不存在")

						//对密码加盐
						newUser.Salt = models.RandStr32()

						//对用户的密码进行哈希替换
						HashUserPass(&newUser)

						models.DB.Create(&newUser) // 传入指针

						//创建用户后写一个log

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
