package routers

//2026-4-2开始每个功能的数据表在各自的api路由下初始化

import (
	"encoding/json"
	"fmt"
	"ops/models"
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
)

func ApiScheduleInit() {
	//先初始化数据表
	models.DB.AutoMigrate(&TabSchedule{})
	models.DB.AutoMigrate(&TabScheduleLog{})
	//获取管理员用户组id
	//先检查用户组有没有这个key
	userGroup.Name = "schedule_admin"

	if models.DB.Where(&userGroup).First(&userGroup).Error == nil {

	} else {
		userGroup.Type = "usergroup"
		models.DB.Create(&userGroup)
	}
}

func ApiSchedule(r *gin.RouterGroup) {
	r.POST("/getevents", func(ctx *gin.Context) {
		data, _ := SeparateData(ctx)
		//fmt.Println(cookieval, data)
		var from fromGetEvents
		if err := mapstructure.Decode(data, &from); err == nil {
			//fmt.Println(from)
			//从数据库获取相关数据
			var list []TabSchedule
			models.DB.Where("start_date <= ? AND end_date >= ?", from.End, from.Start).Where("deleted_at IS NULL").Find(&list)
			fmt.Println(list)
			//ReturnJson(ctx, "ApiOK", list)
			ReturnJson(ctx, "apiOK", gin.H{"list": list})

		} else {
			ReturnJson(ctx, "jsonErr", nil)
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
