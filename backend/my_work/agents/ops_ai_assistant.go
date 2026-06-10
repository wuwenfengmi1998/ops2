package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ScheduleQueryArgs struct {
	Action     string `json:"action,omitempty"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	CalendarID uint   `json:"calendar_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type CurrentUserArgs struct {
	Action string `json:"action,omitempty"`
}

type CurrentUserInfo struct {
	ID         uint   `json:"id,omitempty"`
	UserID     uint   `json:"user_id,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	Username   string `json:"username,omitempty"`
	Birthdate  string `json:"birthdate,omitempty"`
	Gender     string `json:"gender,omitempty"`
	AvatarPath string `json:"avatar_path,omitempty"`
	Region     string `json:"region,omitempty"`
	Language   string `json:"language,omitempty"`
}

type ScheduleCalendar struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ScheduleEvent struct {
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

type ScheduleQuery struct {
	StartDate  string
	EndDate    string
	CalendarID uint
	UserID     uint
	Limit      int
}

type ScheduleProvider interface {
	QuerySchedules(ctx context.Context, query ScheduleQuery) ([]ScheduleEvent, error)
}

var registeredScheduleProvider ScheduleProvider = nil

func RegisterScheduleProvider(provider ScheduleProvider) {
	registeredScheduleProvider = provider
}

func opsAIAssistantScheduleQuerySchema() FunctionToolSchema {
	return FunctionToolSchema{
		Name:        "ops_ai_assistant_schedule_query",
		Description: "只读工具：按明确日期范围查询当前用户可见的 OPS 日历/日程。禁止新增、修改、删除数据，禁止批量添加、批量导入或执行任何写操作。相对日期需先调用 time 获取 start_date/end_date。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"query"},
					"description": "固定为 query。本工具只允许查询，禁止 create/update/delete/batch_add/import 等写操作。",
				},
				"start_date":  map[string]interface{}{"type": "string", "description": "开始日期，格式 YYYY-MM-DD。"},
				"end_date":    map[string]interface{}{"type": "string", "description": "结束日期，格式 YYYY-MM-DD。"},
				"calendar_id": map[string]interface{}{"type": "integer", "description": "可选日历 ID；不传则查询全部可见日程。"},
				"limit":       map[string]interface{}{"type": "integer", "description": "可选返回上限，默认 100，最大 200。"},
			},
			"required": []string{"action", "start_date", "end_date"},
		},
		// ScheduleQueryResult is the return type packed as JSON
	}
}

func opsAIAssistantCurrentUserSchema() FunctionToolSchema {
	return FunctionToolSchema{
		Name:        "ops_ai_assistant_current_user",
		Description: "只读工具：当用户询问“我是谁”“当前登录用户是谁”“我的用户信息”等当前身份问题时调用。先判断是否已登录；已登录则返回当前登录用户信息，未登录则返回 unknown 并提示需要登录才能获取信息。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"get"},
					"description": "固定为 get。本工具只读取当前登录用户信息。",
				},
			},
			"required": []string{"action"},
		},
	}
}

func executeOpsAIAssistantCurrentUser(ctx context.Context, runtime FunctionToolRuntime, rawArgs []byte) ([]byte, error) {
	var args CurrentUserArgs
	if len(rawArgs) > 0 {
		if err := json.Unmarshal(rawArgs, &args); err != nil {
			return nil, err
		}
	}
	if args.Action != "" && args.Action != "get" {
		return json.Marshal(map[string]interface{}{
			"ok":    false,
			"error": "ops_ai_assistant_current_user 是只读工具，仅允许 get 读取当前用户信息",
		})
	}
	if runtime.UserID <= 0 {
		return json.Marshal(map[string]interface{}{
			"ok":       true,
			"loggedIn": false,
			"unknown":  true,
			"message":  "不知道你是谁，需要登录才能获取当前用户信息。",
		})
	}

	return json.Marshal(map[string]interface{}{
		"ok":       true,
		"loggedIn": true,
		"user": map[string]interface{}{
			"id":    runtime.UserID,
			"name":  runtime.UserName,
			"email": runtime.UserEmail,
			"type":  runtime.UserType,
			"info":  runtime.UserInfo,
		},
	})
}

func executeOpsAIAssistantScheduleQuery(ctx context.Context, runtime FunctionToolRuntime, rawArgs []byte) ([]byte, error) {
	var args ScheduleQueryArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return nil, err
	}
	if args.Action != "" && args.Action != "query" {
		return json.Marshal(map[string]interface{}{
			"ok":    false,
			"error": "ops_ai_assistant 是只读工具，仅允许 query 查询操作，禁止新增、修改、删除、批量添加或导入数据",
		})
	}
	startDate, err := time.Parse("2006-01-02", args.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", args.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}
	limit := args.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}

	if registeredScheduleProvider == nil {
		return json.Marshal(map[string]interface{}{
			"ok":    false,
			"error": "日程查询服务未注册",
		})
	}

	events, err := registeredScheduleProvider.QuerySchedules(ctx, ScheduleQuery{
		StartDate:  args.StartDate,
		EndDate:    args.EndDate,
		CalendarID: args.CalendarID,
		UserID:     runtime.UserID,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}
	if events == nil {
		events = []ScheduleEvent{}
	}

	return json.Marshal(map[string]interface{}{
		"ok":         true,
		"start_date": args.StartDate,
		"end_date":   args.EndDate,
		"count":      len(events),
		"limit":      limit,
		"events":     events,
	})
}
