package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ScheduleQueryArgs struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	CalendarID uint   `json:"calendar_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
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
		Description: "按明确日期范围查询当前用户可见的 OPS 日历/日程。相对日期需先调用 time 获取 start_date/end_date。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"start_date":  map[string]interface{}{"type": "string", "description": "开始日期，格式 YYYY-MM-DD。"},
				"end_date":    map[string]interface{}{"type": "string", "description": "结束日期，格式 YYYY-MM-DD。"},
				"calendar_id": map[string]interface{}{"type": "integer", "description": "可选日历 ID；不传则查询全部可见日程。"},
				"limit":       map[string]interface{}{"type": "integer", "description": "可选返回上限，默认 100，最大 200。"},
			},
			"required": []string{"start_date", "end_date"},
		},
		// ScheduleQueryResult is the return type packed as JSON
	}
}

func executeOpsAIAssistantScheduleQuery(ctx context.Context, runtime FunctionToolRuntime, rawArgs []byte) ([]byte, error) {
	var args ScheduleQueryArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return nil, err
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
