package routers

//2026-4-2开始每个功能的数据表在各自的api路由下初始化

import (
	"encoding/json"
	"ops/models"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TabSchedule struct {
	ID        uint   `gorm:"primarykey"`
	UserID    uint   `gorm:"not null;comment:创建人ID"`
	Title     string `gorm:"size:200;not null;comment:日程标题"`
	StartDate string `gorm:"size:10;not null;index;comment:开始日期 YYYY-MM-DD"`
	EndDate   string `gorm:"size:10;not null;index;comment:结束日期 YYYY-MM-DD"`
	BgColor   string `gorm:"size:50;default:#3788d9;comment:背景颜色"`
	Remark    string `gorm:"type:text;comment:备注"`

	CreatedAt *time.Time     `gorm:"type:datetime;autoCreateTime;comment:创建时间"`
	UpdatedAt *time.Time     `gorm:"type:datetime;autoUpdateTime;comment:最后修改时间"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TabScheduleLog struct {
	ID         uint   `gorm:"primarykey"`
	ScheduleID uint   `gorm:"not null;index;comment:关联日程ID"`
	UserID     uint   `gorm:"not null;comment:操作人ID"`
	ActionType string `gorm:"size:50;not null;comment:操作类型: create-创建 update-修改 delete-删除 query-查询"`
	OldContent string `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent string `gorm:"type:text;comment:修改后内容(JSON)"`
	IP         string `gorm:"size:50;comment:操作IP"`
	Remark     string `gorm:"size:500;comment:备注/操作描述"`

	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime;comment:操作时间"`
}
type fromAddEvent struct {
	ID uint 	 `json:"id"`
	Title string `json:"title" binding:"required"` // 日程标题
	Start string `json:"start" binding:"required"` // 开始日期
	End   string `json:"end" binding:"required"`   // 结束日期
	Color string `json:"color" binding:"required"` // 背景颜色
}

type fromGetEvents struct {
	Start string `json:"start" binding:"required"` // 开始日期
	End   string `json:"end" binding:"required"`   // 结束日期
}

var (
	userGroup models.TabUserGroups_
	scheduleAdmins []uint
)

//更新管理员成员缓存
func updateAdminsCash(){
//先清空切片
scheduleAdmins=nil

//获取管理员用户组id
//id 1是系统管理员 直接appen
	scheduleAdmins=append(scheduleAdmins, 1)
	//读取所有绑定了这个用户组的用户id
	var usergroupbind []models.TabUserGroupBinds_
	usergroupbindfind:= models.TabUserGroupBinds_{
		GroupID: userGroup.ID,
	}
	models.DB.Where(&usergroupbindfind).Find(&usergroupbind)
	for _ , item:= range usergroupbind{
		if !slices.Contains(scheduleAdmins,item.UserID){
			scheduleAdmins=append(scheduleAdmins, item.UserID)
		}
	}

}

func ApiScheduleInit() {
	//先初始化数据表
	models.DB.AutoMigrate(&TabSchedule{})
	models.DB.AutoMigrate(&TabScheduleLog{})
	
	//先检查用户组有没有这个key
	userGroup.Name = "schedule_admin"
	if models.DB.Where(&userGroup).First(&userGroup).Error == nil {
		updateAdminsCash()
	} else {
		userGroup.Type = "usergroup"
		models.DB.Create(&userGroup)
	}
}

func ApiSchedule(r *gin.RouterGroup) {
	r.POST("/getevents", func(ctx *gin.Context) {
		data, cookie := SeparateData(ctx)
		user, er := AuthenticationAuthorityFromCookie(cookie)

		var from fromGetEvents
		if err := mapstructure.Decode(data, &from); err == nil {
			//fmt.Println(from)
			//从数据库获取相关数据
			var list []TabSchedule
			models.DB.Where("start_date <= ? AND end_date >= ?", from.End, from.Start).Where("deleted_at IS NULL").Find(&list)
			var relist []map[string]interface{}
			for _, item := range list {
				data, _ := json.Marshal(item)
				var temp map[string]interface{}
				json.Unmarshal(data, &temp)
				// 加自定义字段
				if er == nil {
					//已登录 进一步判断编辑权限
					temp["edit"] = false

					if slices.Contains(scheduleAdmins,user.ID){
						temp["edit"] = true
					}
					if item.UserID == user.ID {
						temp["edit"] = true
					}

					// user_group_find := models.TabUserGroupBinds_{}
					// if models.DB.Where("user_id = ? AND group_id = ?", user.ID, userGroup.ID).First(&user_group_find).Error == nil { //是应用管理员
					// 	temp["edit"] = true
					// }
				} else {
					temp["edit"] = false
				}

				relist = append(relist, temp)

			}

			//ReturnJson(ctx, "ApiOK", list)
			ReturnJson(ctx, "apiOK", gin.H{"list": relist})

		} else {
			ReturnJson(ctx, "jsonErr", nil)
		}

	})

	r.POST("/deleevent", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromAddEvent
			if err := mapstructure.Decode(data, &from); err == nil {
				//先从数据库拉取原始event数据
				 oldEvent:=TabSchedule{}
				 if models.DB.Where("id = ?", from.ID).First(&oldEvent).Error==nil{
					//需要先判断修改权限
					var isCanEdit=false
					if slices.Contains(scheduleAdmins,user.ID){ //用户id是管理员
						isCanEdit = true
					}
					if oldEvent.UserID==user.ID{//event是用户创建的
						isCanEdit = true
					}
					if isCanEdit{
						tosql := TabSchedule{}
						tosql.DeletedAt.Scan(time.Now())
						//fmt.Println(tosql)
						findEvent:=TabSchedule{
							ID: oldEvent.ID,
						}
						if models.DB.Where(&findEvent).Updates(&tosql).Error==nil{
							//应该修改完了  写日志
							//把最新数据再读出来
							models.DB.Where(&findEvent).First(&findEvent)
							newContent, _ := json.Marshal(findEvent) //转 JSON
							oldContent, _ := json.Marshal(oldEvent) //转 JSON
							tosqllog := TabScheduleLog{
								UserID:     user.ID,
								ScheduleID: oldEvent.ID,
								ActionType: "delete",
								NewContent: string(newContent),
								OldContent: string(oldContent),
								IP:         ctx.ClientIP(),
							}
							models.DB.Create(&tosqllog)
							ReturnJson(ctx, "apiOK", nil)
						}else{
							ReturnJson(ctx, "apiErr", nil)
						}
					}else{
						ReturnJson(ctx, "schedule_permission_denied", nil)
					}
				}else{
					ReturnJson(ctx, "schedule_event_not_find", nil)
				}
			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

	})

	r.POST("/editevent", func(ctx *gin.Context) {
	isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			
			var from fromAddEvent
			if err := mapstructure.Decode(data, &from); err == nil {
				//先从数据库拉取原始event数据
				oldEvent:=TabSchedule{}
				if models.DB.Where("id = ?", from.ID).First(&oldEvent).Error==nil{
					//需要先判断修改权限
					var isCanEdit=false
					if slices.Contains(scheduleAdmins,user.ID){ //用户id是管理员
						isCanEdit = true
					}
					if oldEvent.UserID==user.ID{//event是用户创建的
						isCanEdit = true
					}

					if isCanEdit{
						tosql := TabSchedule{
							// UserID:    user.ID,  //如果是管理员修改的话会覆盖掉创建者的id
							Title:     from.Title,
							StartDate: from.Start,
							EndDate:   from.End,
							BgColor:   from.Color,
						}
						//fmt.Println(tosql)

						findEvent:=TabSchedule{
							ID: oldEvent.ID,
						}

						if models.DB.Where(&findEvent).Updates(&tosql).Error==nil{
							//应该修改完了  写日志
							//把最新数据再读出来
							models.DB.Where(&findEvent).First(&findEvent)
							newContent, _ := json.Marshal(findEvent) //转 JSON
							oldContent, _ := json.Marshal(oldEvent) //转 JSON
							tosqllog := TabScheduleLog{
								UserID:     user.ID,
								ScheduleID: oldEvent.ID,
								ActionType: "update",
								NewContent: string(newContent),
								OldContent: string(oldContent),
								IP:         ctx.ClientIP(),
							}
							models.DB.Create(&tosqllog)
							ReturnJson(ctx, "apiOK", nil)


						}else{
							ReturnJson(ctx, "apiErr", nil)
						}
						
					}else{
						ReturnJson(ctx, "schedule_permission_denied", nil)
					}

				}else{
					ReturnJson(ctx, "schedule_event_not_find", nil)
				}

				
				

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})
	r.POST("/addevent", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {

			var from fromAddEvent
			if err := mapstructure.Decode(data, &from); err == nil {

				tosql := TabSchedule{
					UserID:    user.ID,
					Title:     from.Title,
					StartDate: from.Start,
					EndDate:   from.End,
					BgColor:   from.Color,
				}
				if models.DB.Create(&tosql).Error == nil {
					//记录日志
					newContent, _ := json.Marshal(tosql) // 👈 转 JSON
					tosqllog := TabScheduleLog{
						UserID:     user.ID,
						ScheduleID: tosql.ID,
						ActionType: "create",
						NewContent: string(newContent), // 👈 直接赋值
						OldContent: "",
						IP:         ctx.ClientIP(),
					}
					models.DB.Create(&tosqllog)
					ReturnJson(ctx, "apiOK", nil)
				} else {
					ReturnJson(ctx, "apiErr", nil)
				}

			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}

	})

}
