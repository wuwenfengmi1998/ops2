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

type PurchaseQueryArgs struct {
	Action    string `json:"action"`
	OrderID   uint   `json:"order_id,omitempty"`
	Search    string `json:"search,omitempty"`
	Status    string `json:"status,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type PurchaseQuery struct {
	Action    string
	OrderID   uint
	Search    string
	Status    string
	StartDate string
	EndDate   string
	Page      int
	Limit     int
	UserID    uint
}

type PurchaseCost struct {
	ID           uint   `json:"id"`
	OrderID      uint   `json:"order_id"`
	UserID       uint   `json:"user_id"`
	Price        int    `json:"price"`
	Quantity     int    `json:"quantity"`
	CurrencyType int    `json:"currency_type"`
	CurrencyName string `json:"currency_name"`
	CostType     int    `json:"cost_type"`
	CostTypeName string `json:"cost_type_name"`
}

type PurchaseCommit struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	UserID    uint   `json:"user_id"`
	Action    string `json:"action"`
	Status    string `json:"status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type PurchaseOrder struct {
	ID              uint             `json:"id"`
	UserID          uint             `json:"user_id"`
	Title           string           `json:"title"`
	Remark          string           `json:"remark,omitempty"`
	Link            string           `json:"link,omitempty"`
	DetailURL       string           `json:"detail_url"`
	Styles          string           `json:"styles,omitempty"`
	OrderStatus     string           `json:"order_status"`
	OrderStatusName string           `json:"order_status_name"`
	CreatedAt       string           `json:"created_at,omitempty"`
	UpdatedAt       string           `json:"updated_at,omitempty"`
	CanModify       bool             `json:"can_modify,omitempty"`
	Costs           []PurchaseCost   `json:"costs,omitempty"`
	Commits         []PurchaseCommit `json:"commits,omitempty"`
}

type PurchaseQueryResult struct {
	Ok       bool                   `json:"ok"`
	Action   string                 `json:"action"`
	LoggedIn bool                   `json:"loggedIn"`
	Count    int                    `json:"count,omitempty"`
	Total    int64                  `json:"total,omitempty"`
	Page     int                    `json:"page,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Orders   []PurchaseOrder        `json:"orders,omitempty"`
	Order    *PurchaseOrder         `json:"order,omitempty"`
	Counts   map[string]int64       `json:"counts,omitempty"`
	Filters  map[string]interface{} `json:"filters,omitempty"`
	Message  string                 `json:"message,omitempty"`
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

type PurchaseProvider interface {
	QueryPurchases(ctx context.Context, query PurchaseQuery) (*PurchaseQueryResult, error)
}

var registeredScheduleProvider ScheduleProvider = nil
var registeredPurchaseProvider PurchaseProvider = nil

func RegisterScheduleProvider(provider ScheduleProvider) {
	registeredScheduleProvider = provider
}

func RegisterPurchaseProvider(provider PurchaseProvider) {
	registeredPurchaseProvider = provider
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

func opsAIAssistantPurchaseQuerySchema() FunctionToolSchema {
	return FunctionToolSchema{
		Name:        "ops_ai_assistant_purchase_query",
		Description: "只读工具：查询采购模块订单。用户询问采购订单、采购状态、待处理/已下单/已到达/已收件/丢件/退件数量或列表、指定采购订单详情时调用。禁止新增、修改、删除订单或状态。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"list", "get", "count"},
					"description": "list 查询订单列表；get 查询单个订单详情；count 统计各状态数量。",
				},
				"order_id":   map[string]interface{}{"type": "integer", "description": "action=get 时的采购订单 ID。"},
				"search":     map[string]interface{}{"type": "string", "description": "按订单 ID、标题或备注搜索。"},
				"status":     map[string]interface{}{"type": "string", "enum": []string{"", "pending", "ordered", "arrived", "received", "lost", "returned"}, "description": "订单状态过滤：pending 待处理，ordered 已下单，arrived 已到达，received 已收件，lost 丢件，returned 退件。"},
				"start_date": map[string]interface{}{"type": "string", "description": "可选创建日期开始，格式 YYYY-MM-DD。"},
				"end_date":   map[string]interface{}{"type": "string", "description": "可选创建日期结束，格式 YYYY-MM-DD。"},
				"page":       map[string]interface{}{"type": "integer", "description": "分页页码，默认 1。"},
				"limit":      map[string]interface{}{"type": "integer", "description": "返回上限，默认 20，最大 100。"},
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

func executeOpsAIAssistantPurchaseQuery(ctx context.Context, runtime FunctionToolRuntime, rawArgs []byte) ([]byte, error) {
	var args PurchaseQueryArgs
	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return nil, err
	}
	if args.Action != "list" && args.Action != "get" && args.Action != "count" {
		return json.Marshal(map[string]interface{}{
			"ok":    false,
			"error": "ops_ai_assistant_purchase_query 是只读工具，仅允许 list/get/count 查询操作",
		})
	}
	if runtime.UserID <= 0 {
		return json.Marshal(map[string]interface{}{
			"ok":       true,
			"action":   args.Action,
			"loggedIn": false,
			"message":  "需要登录才能查询采购模块信息。",
		})
	}
	if args.Action == "get" && args.OrderID <= 0 {
		return nil, fmt.Errorf("order_id is required when action=get")
	}
	if args.StartDate != "" {
		if _, err := time.Parse("2006-01-02", args.StartDate); err != nil {
			return nil, fmt.Errorf("invalid start_date: %w", err)
		}
	}
	if args.EndDate != "" {
		if _, err := time.Parse("2006-01-02", args.EndDate); err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
	}
	if args.StartDate != "" && args.EndDate != "" {
		startDate, _ := time.Parse("2006-01-02", args.StartDate)
		endDate, _ := time.Parse("2006-01-02", args.EndDate)
		if endDate.Before(startDate) {
			return nil, fmt.Errorf("end_date must be after start_date")
		}
	}
	if args.Page <= 0 {
		args.Page = 1
	}
	if args.Limit <= 0 {
		args.Limit = 20
	}
	if args.Limit > 100 {
		args.Limit = 100
	}
	if registeredPurchaseProvider == nil {
		return json.Marshal(map[string]interface{}{
			"ok":    false,
			"error": "采购查询服务未注册",
		})
	}

	result, err := registeredPurchaseProvider.QueryPurchases(ctx, PurchaseQuery{
		Action:    args.Action,
		OrderID:   args.OrderID,
		Search:    args.Search,
		Status:    args.Status,
		StartDate: args.StartDate,
		EndDate:   args.EndDate,
		Page:      args.Page,
		Limit:     args.Limit,
		UserID:    runtime.UserID,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
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
