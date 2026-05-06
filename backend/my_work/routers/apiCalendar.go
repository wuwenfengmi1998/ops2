package routers

import (
	"encoding/json"
	"ops/models"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// TabCalendar 日历表
type TabCalendar struct {
	ID          uint   `gorm:"primarykey"`
	UserID      uint   `gorm:"not null;comment:创建人ID"`
	Name        string `gorm:"size:100;not null;comment:日历名称"`
	Description string `gorm:"size:500;comment:日历描述"`
	Color       string `gorm:"size:50;default:#3788d9;comment:日历颜色"`
	IsPublic    bool   `gorm:"default:false;comment:是否公开"`

	CreatedAt *time.Time     `gorm:"type:datetime;autoCreateTime;comment:创建时间"`
	UpdatedAt *time.Time     `gorm:"type:datetime;autoUpdateTime;comment:最后修改时间"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TabCalendarEvent 日历事件表
type TabCalendarEvent struct {
	ID         uint `gorm:"primarykey"`
	CalendarID uint `gorm:"not null;index;comment:关联日历ID"`
	UserID     uint `gorm:"not null;comment:创建人ID"`
	//UsersID    []uint     `gorm:"type:json; null;comment:其他关联用户ID"`
	Title     string     `gorm:"size:200;not null;comment:事件标题"`
	StartDate *time.Time `gorm:"size:10;not null;index;comment:开始日期 YYYY-MM-DD"`
	EndDate   *time.Time `gorm:"size:10;not null;index;comment:结束日期 YYYY-MM-DD"`
	IsAllDay  bool       `gorm:"default:true;comment:是否全日事件"`
	BgColor   string     `gorm:"size:50;default:#3788d9;comment:背景颜色"`
	IsPublic  bool       `gorm:"default:false;comment:是否为公共日程"`
	Remark    string     `gorm:"type:text;comment:备注"`

	CreatedAt *time.Time     `gorm:"type:datetime;autoCreateTime;comment:创建时间"`
	UpdatedAt *time.Time     `gorm:"type:datetime;autoUpdateTime;comment:最后修改时间"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TabCalendarLog 日历操作日志表
type TabCalendarLog struct {
	ID         uint   `gorm:"primarykey"`
	CalendarID uint   `gorm:"not null;index;comment:关联日历ID"`
	EventID    uint   `gorm:"not null;index;comment:关联事件ID（可选）"`
	UserID     uint   `gorm:"not null;comment:操作人ID"`
	ActionType string `gorm:"size:50;not null;comment:操作类型: create-创建 update-修改 delete-删除"`
	OldContent string `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent string `gorm:"type:text;comment:修改后内容(JSON)"`
	IP         string `gorm:"size:50;comment:操作IP"`
	Remark     string `gorm:"size:500;comment:备注/操作描述"`

	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime;comment:操作时间"`
}

// 请求结构体
type fromCreateCalendar struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Is_public   bool   `json:"is_public"`
}

type fromUpdateCalendar struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Is_public   bool   `json:"is_public"`
}

type fromDeleteCalendar struct {
	ID uint `json:"id" binding:"required"`
}

type fromGetCalendarEvents struct {
	CalendarID uint       `json:"calendar_id" binding:"required"`
	Start      *time.Time `json:"start" binding:"required"`
	End        *time.Time `json:"end" binding:"required"`
}

type fromAddCalendarEvent struct {
	CalendarID   uint   `json:"calendar_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Start        string `json:"start" binding:"required"`
	End          string `json:"end" binding:"required"`
	Color        string `json:"color"`
	ScheduleType string `json:"schedule_type"`
	Is_public    bool   `json:"is_public"`
	Remark       string `json:"remark"`
}

type fromUpdateCalendarEvent struct {
	ID           uint   `json:"id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Start        string `json:"start" binding:"required"`
	End          string `json:"end" binding:"required"`
	Color        string `json:"color"`
	ScheduleType string `json:"schedule_type"`
	Is_public    bool   `json:"is_public"`
	Remark       string `json:"remark"`
}

type fromDeleteCalendarEvent struct {
	ID uint `json:"id" binding:"required"`
}

var (
	calendarUserGroup TabUserGroups
	calendarAdmins    []uint
)

// CalendarUpdateAdminsCash 更新客户管理员缓存
func CalendarUpdateAdminsCash() {
	calendarAdmins = nil
	calendarAdmins = append(calendarAdmins, 1) // id=1 系统管理员默认拥有所有权限
	var binds []TabUserGroupBinds
	models.DB.Where("group_id = ?", calendarUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(calendarAdmins, item.UserID) {
			calendarAdmins = append(calendarAdmins, item.UserID)
		}
	}
}

// canModifyCustomer 判断是否有权限修改/删除客户（创建者或管理员）
func canModifyCalendar(userID, creatorUserID uint) bool {
	if slices.Contains(calendarAdmins, userID) {
		return true
	}
	return userID == creatorUserID
}

func ApiCalendarInit() {
	// 初始化数据表
	models.DB.AutoMigrate(&TabCalendar{})
	models.DB.AutoMigrate(&TabCalendarEvent{})
	models.DB.AutoMigrate(&TabCalendarLog{})

	// 自动创建 calendar_admin 用户组
	models.DB.Where("name = ?", "calendar_admin").FirstOrCreate(&calendarUserGroup, TabUserGroups{
		Name: "calendar_admin",
		Type: "usergroup",
	})
}

func ApiCalendar(r *gin.RouterGroup) {
	// 创建日历
	r.POST("/calendar/create", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromCreateCalendar
			if err := mapstructure.Decode(data, &from); err == nil {
				calendar := TabCalendar{
					UserID:      user.ID,
					Name:        from.Name,
					Description: from.Description,
					Color:       from.Color,
					IsPublic:    from.Is_public,
				}
				if calendar.Color == "" {
					calendar.Color = "#3788d9"
				}
				if models.DB.Create(&calendar).Error == nil {
					// 记录日志
					newContent, _ := json.Marshal(calendar)
					log := TabCalendarLog{
						CalendarID: calendar.ID,
						UserID:     user.ID,
						ActionType: "create",
						NewContent: string(newContent),
						IP:         ctx.ClientIP(),
					}
					models.DB.Create(&log)
					ReturnJson(ctx, "apiOK", gin.H{"id": calendar.ID})
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

	// 获取日历列表（不需要登录）
	r.POST("/calendar/list", func(ctx *gin.Context) {
		isAuth, user, _ := AuthenticationAuthority(ctx)

		var calendars []TabCalendar
		models.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&calendars)

		type CalendarWithEdit struct {
			TabCalendar
			CanEdit bool `json:"canEdit"`
		}
		var result []CalendarWithEdit
		for _, cal := range calendars {
			// 私有日历：只有创建者可见
			if !cal.IsPublic {
				if !isAuth || cal.UserID != user.ID {
					continue
				}
				result = append(result, CalendarWithEdit{
					TabCalendar: cal,
					CanEdit:     true,
				})
				continue
			}
			// 公开日历
			canEdit := false
			if isAuth {
				canEdit = canModifyCalendar(user.ID, cal.UserID)
			}
			result = append(result, CalendarWithEdit{
				TabCalendar: cal,
				CanEdit:     canEdit,
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{"list": result})

	})

	// 更新日历
	r.POST("/calendar/update", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromUpdateCalendar
			if err := mapstructure.Decode(data, &from); err == nil {
				oldCalendar := TabCalendar{}
				if models.DB.Where("id = ?", from.ID).First(&oldCalendar).Error == nil {
					// 检查权限（只有创建人可以修改）
					if oldCalendar.UserID != user.ID {
						ReturnJson(ctx, "permission_denied", nil)
						return
					}

					newCalendar := TabCalendar{
						Name:        from.Name,
						Description: from.Description,
						Color:       from.Color,
						IsPublic:    from.Is_public,
					}
					if newCalendar.Color == "" {
						newCalendar.Color = "#3788d9"
					}

					if models.DB.Model(&oldCalendar).Updates(&newCalendar).Error == nil {
						// 记录日志
						newContent, _ := json.Marshal(newCalendar)
						oldContent, _ := json.Marshal(oldCalendar)
						log := TabCalendarLog{
							CalendarID: oldCalendar.ID,
							UserID:     user.ID,
							ActionType: "update",
							OldContent: string(oldContent),
							NewContent: string(newContent),
							IP:         ctx.ClientIP(),
						}
						models.DB.Create(&log)
						ReturnJson(ctx, "apiOK", nil)
					} else {
						ReturnJson(ctx, "apiErr", nil)
					}
				} else {
					ReturnJson(ctx, "calendar_not_find", nil)
				}
			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	// 删除日历
	r.POST("/calendar/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromDeleteCalendar
			if err := mapstructure.Decode(data, &from); err == nil {
				oldCalendar := TabCalendar{}
				if models.DB.Where("id = ?", from.ID).First(&oldCalendar).Error == nil {
					// 检查权限（只有创建人可以删除）
					if oldCalendar.UserID != user.ID {
						ReturnJson(ctx, "permission_denied", nil)
						return
					}

					// 软删除日历
					if models.DB.Delete(&oldCalendar).Error == nil {
						// 记录日志
						oldContent, _ := json.Marshal(oldCalendar)
						log := TabCalendarLog{
							CalendarID: oldCalendar.ID,
							UserID:     user.ID,
							ActionType: "delete",
							OldContent: string(oldContent),
							IP:         ctx.ClientIP(),
						}
						models.DB.Create(&log)
						ReturnJson(ctx, "apiOK", nil)
					} else {
						ReturnJson(ctx, "apiErr", nil)
					}
				} else {
					ReturnJson(ctx, "calendar_not_find", nil)
				}
			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	// 获取日历事件
	r.POST("/calendar/events", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromGetCalendarEvents
			if err := mapstructure.Decode(data, &from); err == nil {
				var events []TabCalendarEvent
				models.DB.Where("calendar_id = ? AND start_date <= ? AND end_date >= ? AND deleted_at IS NULL",
					from.CalendarID, from.End, from.Start).Find(&events)

				// 为事件添加编辑权限标识
				var relist []map[string]interface{}
				for _, event := range events {
					data, _ := json.Marshal(event)
					var temp map[string]interface{}
					json.Unmarshal(data, &temp)
					// 这里可以根据需要添加 edit 字段
					relist = append(relist, temp)
				}

				ReturnJson(ctx, "apiOK", gin.H{"list": relist})
			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	// 添加日历事件
	r.POST("/calendar/addevent", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			// 先检查必需字段
			calendarIDRaw, ok := data["calendar_id"].(float64)
			if !ok || calendarIDRaw == 0 {
				ReturnJson(ctx, "jsonErr", nil)
				return
			}
			calendarID := uint(calendarIDRaw)

			// 检查日历是否存在
			var calendar TabCalendar
			if models.DB.Where("id = ? AND deleted_at IS NULL", calendarID).First(&calendar).Error != nil {
				ReturnJson(ctx, "calendar_not_find", nil)
				return
			}

			// 解析日期
			startStr, _ := data["start"].(string)
			endStr, _ := data["end"].(string)
			title, _ := data["title"].(string)
			color, _ := data["color"].(string)
			remark, _ := data["remark"].(string)
			isPublic, _ := data["is_public"].(bool)

			startDate, _ := time.Parse("2006-01-02 15:04:05", startStr)
			endDate, _ := time.Parse("2006-01-02 15:04:05", endStr)

			event := TabCalendarEvent{
				CalendarID: calendarID,
				UserID:     user.ID,
				Title:      title,
				StartDate:  &startDate,
				EndDate:    &endDate,
				BgColor:    color,
				IsPublic:   isPublic,
				Remark:     remark,
			}
			if event.BgColor == "" {
				event.BgColor = calendar.Color
			}

			if models.DB.Create(&event).Error == nil {
				// 记录日志
				newContent, _ := json.Marshal(event)
				log := TabCalendarLog{
					CalendarID: event.CalendarID,
					EventID:    event.ID,
					UserID:     user.ID,
					ActionType: "create_event",
					NewContent: string(newContent),
					IP:         ctx.ClientIP(),
				}
				models.DB.Create(&log)
				ReturnJson(ctx, "apiOK", gin.H{"id": event.ID})
			} else {
				ReturnJson(ctx, "apiErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	// 更新日历事件
	r.POST("/calendar/updateevent", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			// 先检查必需字段
			idRaw, ok := data["id"].(float64)
			if !ok || idRaw == 0 {
				ReturnJson(ctx, "jsonErr", nil)
				return
			}
			eventID := uint(idRaw)

			oldEvent := TabCalendarEvent{}
			if models.DB.Where("id = ?", eventID).First(&oldEvent).Error == nil {
				// 检查权限（只有创建人可以修改）
				if oldEvent.UserID != user.ID {
					ReturnJson(ctx, "permission_denied", nil)
					return
				}

				// 解析字段
				startStr, _ := data["start"].(string)
				endStr, _ := data["end"].(string)
				title, _ := data["title"].(string)
				color, _ := data["color"].(string)
				remark, _ := data["remark"].(string)
				isPublic, _ := data["is_public"].(bool)

				startDate, _ := time.Parse("2006-01-02 15:04:05", startStr)
				endDate, _ := time.Parse("2006-01-02 15:04:05", endStr)

				newEvent := TabCalendarEvent{
					Title:     title,
					StartDate: &startDate,
					EndDate:   &endDate,
					BgColor:   color,
					IsPublic:  isPublic,
					Remark:    remark,
				}
				if newEvent.BgColor == "" {
					// 获取日历颜色
					var calendar TabCalendar
					models.DB.Where("id = ?", oldEvent.CalendarID).First(&calendar)
					newEvent.BgColor = calendar.Color
				}

				if models.DB.Model(&oldEvent).Updates(&newEvent).Error == nil {
					// 记录日志
					newContent, _ := json.Marshal(newEvent)
					oldContent, _ := json.Marshal(oldEvent)
					log := TabCalendarLog{
						CalendarID: oldEvent.CalendarID,
						EventID:    oldEvent.ID,
						UserID:     user.ID,
						ActionType: "update_event",
						OldContent: string(oldContent),
						NewContent: string(newContent),
						IP:         ctx.ClientIP(),
					}
					models.DB.Create(&log)
					ReturnJson(ctx, "apiOK", nil)
				} else {
					ReturnJson(ctx, "apiErr", nil)
				}
			} else {
				ReturnJson(ctx, "event_not_find", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})

	// 删除日历事件
	r.POST("/calendar/deleteevent", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if isAuth {
			var from fromDeleteCalendarEvent
			if err := mapstructure.Decode(data, &from); err == nil {
				oldEvent := TabCalendarEvent{}
				if models.DB.Where("id = ?", from.ID).First(&oldEvent).Error == nil {
					// 检查权限（只有创建人可以删除）
					if oldEvent.UserID != user.ID {
						ReturnJson(ctx, "permission_denied", nil)
						return
					}

					if models.DB.Delete(&oldEvent).Error == nil {
						// 记录日志
						oldContent, _ := json.Marshal(oldEvent)
						log := TabCalendarLog{
							CalendarID: oldEvent.CalendarID,
							EventID:    oldEvent.ID,
							UserID:     user.ID,
							ActionType: "delete_event",
							OldContent: string(oldContent),
							IP:         ctx.ClientIP(),
						}
						models.DB.Create(&log)
						ReturnJson(ctx, "apiOK", nil)
					} else {
						ReturnJson(ctx, "apiErr", nil)
					}
				} else {
					ReturnJson(ctx, "event_not_find", nil)
				}
			} else {
				ReturnJson(ctx, "jsonErr", nil)
			}
		} else {
			ReturnJson(ctx, "userCookieError", nil)
		}
	})
}
