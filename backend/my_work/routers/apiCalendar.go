package routers

import (
	"context"
	"encoding/json"
	"ops/agents"
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
	ID           uint       `gorm:"primarykey"`
	CalendarID   uint       `gorm:"not null;index;comment:关联日历ID"`
	UserID       uint       `gorm:"not null;comment:创建人ID"`
	Title        string     `gorm:"size:200;not null;comment:事件标题"`
	StartDate    *time.Time `gorm:"size:10;not null;index;comment:开始日期 YYYY-MM-DD"`
	EndDate      *time.Time `gorm:"size:10;not null;index;comment:结束日期 YYYY-MM-DD"`
	IsAllDay     bool       `gorm:"default:true;comment:是否全日事件"`
	ScheduleType string     `gorm:"size:50;default:work;comment:日程类型: work-工作 duty-值班 exam-考试 standby-待命 personal_holiday-调休 personal_holiday-公假"`
	BgColor      string     `gorm:"size:50;default:#3788d9;comment:背景颜色"`
	IsPublic     bool       `gorm:"default:false;comment:是否为公共日程"`
	Remark       string     `gorm:"type:text;comment:备注"`

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

type fromRestoreCalendar struct {
	ID uint `json:"id" binding:"required"`
}

var (
	calendarUserGroup TabUserGroups
	calendarAdmins    []uint
)

type CalendarScheduleQuery struct {
	CalendarID uint
	StartDate  time.Time
	EndDate    time.Time
	User       *TabUser
	Limit      int
}

type CalendarScheduleEvent struct {
	ID           uint   `json:"id"`
	CalendarID   uint   `json:"calendar_id"`
	UserID       uint   `json:"user_id"`
	Title        string `json:"title"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	ScheduleType string `json:"schedule_type"`
	IsPublic     bool   `json:"is_public"`
	Remark       string `json:"remark,omitempty"`
	AccessNote   string `json:"access_note,omitempty"`
	CanEdit      bool   `json:"canEdit"`
}

type calendarScheduleProvider struct{}

func (calendarScheduleProvider) QuerySchedules(ctx context.Context, query agents.ScheduleQuery) ([]agents.ScheduleEvent, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		return nil, err
	}

	var user *TabUser
	if query.UserID > 0 {
		user = &TabUser{ID: query.UserID}
	}
	events, err := QueryCalendarSchedulesForAI(CalendarScheduleQuery{
		CalendarID: query.CalendarID,
		StartDate:  startDate,
		EndDate:    endDate,
		User:       user,
		Limit:      query.Limit,
	})
	if err != nil {
		return nil, err
	}

	result := make([]agents.ScheduleEvent, 0, len(events))
	for _, event := range events {
		result = append(result, agents.ScheduleEvent{
			ID:           event.ID,
			CalendarID:   event.CalendarID,
			UserID:       event.UserID,
			Title:        event.Title,
			StartDate:    event.StartDate,
			EndDate:      event.EndDate,
			ScheduleType: event.ScheduleType,
			IsPublic:     event.IsPublic,
			Remark:       event.Remark,
			AccessNote:   event.AccessNote,
			CanEdit:      event.CanEdit,
		})
	}
	return result, nil
}

func QueryCalendarSchedulesForAI(query CalendarScheduleQuery) ([]CalendarScheduleEvent, error) {
	startDate := dateOnly(query.StartDate)
	endDate := dateOnly(query.EndDate)
	if endDate.Before(startDate) {
		return nil, nil
	}

	currentUserID := uint(0)
	if query.User != nil {
		currentUserID = query.User.ID
	}

	db := models.DB.Where("start_date <= ? AND end_date >= ? AND deleted_at IS NULL", &endDate, &startDate)
	if query.CalendarID > 0 {
		db = db.Where("calendar_id = ? OR is_public = ?", query.CalendarID, true)
	} else if currentUserID > 0 {
		db = db.Where("is_public = ? OR calendar_id IN (?)", true, models.DB.Model(&TabCalendar{}).Select("id").Where("deleted_at IS NULL AND (is_public = ? OR user_id = ?)", true, currentUserID))
	} else {
		db = db.Where("is_public = ? OR calendar_id IN (?)", true, models.DB.Model(&TabCalendar{}).Select("id").Where("deleted_at IS NULL AND is_public = ?", true))
	}

	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}

	var events []TabCalendarEvent
	if err := db.Order("start_date asc, end_date asc, id asc").Find(&events).Error; err != nil {
		return nil, err
	}

	calendarIDs := make([]uint, 0)
	calendarIDSet := map[uint]bool{}
	if query.CalendarID > 0 {
		calendarIDs = append(calendarIDs, query.CalendarID)
		calendarIDSet[query.CalendarID] = true
	}
	for _, event := range events {
		if !calendarIDSet[event.CalendarID] {
			calendarIDs = append(calendarIDs, event.CalendarID)
			calendarIDSet[event.CalendarID] = true
		}
	}

	calendarCreators := map[uint]uint{}
	if len(calendarIDs) > 0 {
		var calendars []TabCalendar
		models.DB.Where("id IN ?", calendarIDs).Find(&calendars)
		for _, calendar := range calendars {
			calendarCreators[calendar.ID] = calendar.UserID
		}
	}

	result := make([]CalendarScheduleEvent, 0, len(events))
	for _, event := range events {
		canEdit := false
		accessNote := ""
		if currentUserID > 0 {
			calendarCreatorID := calendarCreators[event.CalendarID]
			canEdit = event.UserID == currentUserID || calendarCreatorID == currentUserID || slices.Contains(calendarAdmins, currentUserID)
			if !event.IsPublic && event.ScheduleType == "work" {
				accessNote = "非公开工作日程，仅因当前用户具备相关权限可见"
			} else if !event.IsPublic {
				accessNote = "非公开日程，仅因当前用户具备相关权限可见"
			}
		}
		result = append(result, CalendarScheduleEvent{
			ID:           event.ID,
			CalendarID:   event.CalendarID,
			UserID:       event.UserID,
			Title:        event.Title,
			StartDate:    formatDatePtr(event.StartDate),
			EndDate:      formatDatePtr(event.EndDate),
			ScheduleType: event.ScheduleType,
			IsPublic:     event.IsPublic,
			Remark:       event.Remark,
			AccessNote:   accessNote,
			CanEdit:      canEdit,
		})
	}
	return result, nil
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func formatDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// CalendarUpdateAdminsCash
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

	CalendarUpdateAdminsCash()
	agents.RegisterScheduleProvider(calendarScheduleProvider{})
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

	// 获取所有日历（包括已删除的，管理员专用）
	r.POST("/calendar/list_all", func(ctx *gin.Context) {
		isAuth, user, _ := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		// 限制只有日历管理员可访问
		if !slices.Contains(calendarAdmins, user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		// 使用 Unscoped 查询所有日历（包括软删除的）
		var calendars []TabCalendar
		models.DB.Unscoped().Order("created_at DESC").Find(&calendars)

		// 一次性查询所有日历的事件数量（仅统计未删除的事件）
		type calendarEventCount struct {
			CalendarID uint `gorm:"column:calendar_id"`
			Cnt        int  `gorm:"column:cnt"`
		}
		var rows []calendarEventCount
		models.DB.Model(&TabCalendarEvent{}).
			Select("calendar_id, COUNT(*) as cnt").
			Where("deleted_at IS NULL").
			Group("calendar_id").
			Scan(&rows)
		eventCountMap := make(map[uint]int)
		for _, row := range rows {
			eventCountMap[row.CalendarID] = row.Cnt
		}

		type CalendarWithEdit struct {
			TabCalendar
			CanEdit    bool `json:"canEdit"`
			EventCount int  `json:"event_count"`
		}
		var result []CalendarWithEdit
		for _, cal := range calendars {
			result = append(result, CalendarWithEdit{
				TabCalendar: cal,
				CanEdit:     true,
				EventCount:  eventCountMap[cal.ID],
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{"list": result})
	})

	// 恢复已删除的日历
	r.POST("/calendar/restore", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}

		// 限制只有日历管理员可操作
		if !slices.Contains(calendarAdmins, user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var from fromRestoreCalendar
		if err := mapstructure.Decode(data, &from); err != nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		// 使用 Unscoped 查询（包括软删除的）
		var calendar TabCalendar
		if models.DB.Unscoped().Where("id = ?", from.ID).First(&calendar).Error != nil {
			ReturnJson(ctx, "calendar_not_find", nil)
			return
		}

		// 恢复软删除（将 deleted_at 设为 NULL）
		if models.DB.Unscoped().Model(&calendar).Update("deleted_at", nil).Error != nil {
			ReturnJson(ctx, "apiErr", nil)
			return
		}

		// 记录日志
		newContent, _ := json.Marshal(calendar)
		log := TabCalendarLog{
			CalendarID: calendar.ID,
			UserID:     user.ID,
			ActionType: "restore",
			NewContent: string(newContent),
			IP:         ctx.ClientIP(),
		}
		models.DB.Create(&log)
		ReturnJson(ctx, "apiOK", nil)
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

	// 获取日历事件（公开接口，无需登录）
	r.POST("/calendar/events", func(ctx *gin.Context) {
		data, cookieval := SeparateData(ctx)

		if data == nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		calendarIDRaw, ok := data["calendar_id"].(float64)
		if !ok || calendarIDRaw == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		calendarID := uint(calendarIDRaw)

		startStr, _ := data["start"].(string)
		endStr, _ := data["end"].(string)
		startDate, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		endDate, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}

		var currentUser *TabUser
		if cookieval != "" {
			if user, err := AuthenticationAuthorityFromCookie(cookieval); err == nil {
				currentUser = user
			}
		}

		events, err := QueryCalendarSchedulesForAI(CalendarScheduleQuery{
			CalendarID: calendarID,
			StartDate:  startDate,
			EndDate:    endDate,
			User:       currentUser,
		})
		if err != nil {
			ReturnJson(ctx, "apiErr", nil)
			return
		}

		relist := make([]map[string]interface{}, 0, len(events))
		for _, event := range events {
			relist = append(relist, map[string]interface{}{
				"ID":           event.ID,
				"CalendarID":   event.CalendarID,
				"UserID":       event.UserID,
				"Title":        event.Title,
				"StartDate":    event.StartDate,
				"EndDate":      event.EndDate,
				"ScheduleType": event.ScheduleType,
				"IsPublic":     event.IsPublic,
				"Remark":       event.Remark,
				"canEdit":      event.CanEdit,
			})
		}

		ReturnJson(ctx, "apiOK", gin.H{"list": relist})
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
			remark, _ := data["remark"].(string)
			isPublic, _ := data["is_public"].(bool)
			scheduleType, _ := data["schedule_type"].(string)
			if scheduleType == "" {
				scheduleType = "work"
			}

			startDate, _ := time.Parse("2006-01-02 15:04:05", startStr)
			endDate, _ := time.Parse("2006-01-02 15:04:05", endStr)

			event := TabCalendarEvent{
				CalendarID:   calendarID,
				UserID:       user.ID,
				Title:        title,
				StartDate:    &startDate,
				EndDate:      &endDate,
				ScheduleType: scheduleType,
				IsPublic:     isPublic,
				Remark:       remark,
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
				// 检查权限（事件创建人、日历创建人或管理员可修改）
				var calendarCreatorID uint
				var calendar TabCalendar
				if models.DB.Where("id = ?", oldEvent.CalendarID).First(&calendar).Error == nil {
					calendarCreatorID = calendar.UserID
				}
				if !canModifyCalendar(user.ID, oldEvent.UserID) && calendarCreatorID != user.ID {
					ReturnJson(ctx, "permission_denied", nil)
					return
				}

				// 解析字段
				startStr, _ := data["start"].(string)
				endStr, _ := data["end"].(string)
				title, _ := data["title"].(string)
				remark, _ := data["remark"].(string)
				isPublic, _ := data["is_public"].(bool)
				scheduleType, _ := data["schedule_type"].(string)
				if scheduleType == "" {
					scheduleType = "work"
				}

				startDate, _ := time.Parse("2006-01-02 15:04:05", startStr)
				endDate, _ := time.Parse("2006-01-02 15:04:05", endStr)

				newEvent := TabCalendarEvent{
					Title:        title,
					StartDate:    &startDate,
					EndDate:      &endDate,
					ScheduleType: scheduleType,
					IsPublic:     isPublic,
					Remark:       remark,
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
					// 检查权限（事件创建人、日历创建人或管理员可删除）
					var calendarCreatorID uint
					var calendar TabCalendar
					if models.DB.Where("id = ?", oldEvent.CalendarID).First(&calendar).Error == nil {
						calendarCreatorID = calendar.UserID
					}
					if !canModifyCalendar(user.ID, oldEvent.UserID) && calendarCreatorID != user.ID {
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
