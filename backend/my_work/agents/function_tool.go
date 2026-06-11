package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type FunctionToolSchema struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

type FunctionToolRuntime struct {
	UserID    uint
	UserName  string
	UserEmail string
	UserType  string
	UserInfo  *CurrentUserInfo
}

func FunctionToolSchemas(configs []ToolConfig) []FunctionToolSchema {
	tools := make([]FunctionToolSchema, 0)
	for _, config := range configs {
		if !config.Enabled {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(config.Name)) {
		case "time":
			tools = append(tools, timeFunctionToolSchema())
		case "ops_ai_assistant_schedule_query", "ops_ai_assistant":
			tools = append(tools, opsAIAssistantScheduleQuerySchema())
		case "ops_ai_assistant_calendar_list":
			tools = append(tools, opsAIAssistantCalendarListSchema())
		case "ops_ai_assistant_schedule_create":
			tools = append(tools, opsAIAssistantScheduleCreateSchema())
		case "ops_ai_assistant_schedule_update":
			tools = append(tools, opsAIAssistantScheduleUpdateSchema())
		case "ops_ai_assistant_schedule_delete":
			tools = append(tools, opsAIAssistantScheduleDeleteSchema())
		case "ops_ai_assistant_current_user":
			tools = append(tools, opsAIAssistantCurrentUserSchema())
		case "ops_ai_assistant_purchase_query":
			tools = append(tools, opsAIAssistantPurchaseQuerySchema())
		}
	}
	return tools
}

func ExecuteFunctionTool(ctx context.Context, runtime FunctionToolRuntime, name string, rawArgs []byte) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	switch strings.ToLower(strings.TrimSpace(name)) {
	case "time":
		var args TimeRangeArgs
		if len(rawArgs) > 0 {
			if err := json.Unmarshal(rawArgs, &args); err != nil {
				return nil, err
			}
		}
		result, err := ResolveTimeRange(args, time.Now())
		if err != nil {
			return nil, err
		}
		return json.Marshal(result)
	case "ops_ai_assistant_schedule_query", "ops_ai_assistant":
		return executeOpsAIAssistantScheduleQuery(ctx, runtime, rawArgs)
	case "ops_ai_assistant_calendar_list":
		return executeOpsAIAssistantCalendarList(ctx, runtime, rawArgs)
	case "ops_ai_assistant_schedule_create":
		return executeOpsAIAssistantScheduleCreate(ctx, runtime, rawArgs)
	case "ops_ai_assistant_schedule_update":
		return executeOpsAIAssistantScheduleUpdate(ctx, runtime, rawArgs)
	case "ops_ai_assistant_schedule_delete":
		return executeOpsAIAssistantScheduleDelete(ctx, runtime, rawArgs)
	case "ops_ai_assistant_current_user":
		return executeOpsAIAssistantCurrentUser(ctx, runtime, rawArgs)
	case "ops_ai_assistant_purchase_query":
		return executeOpsAIAssistantPurchaseQuery(ctx, runtime, rawArgs)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func timeFunctionToolSchema() FunctionToolSchema {
	return FunctionToolSchema{
		Name:        "time",
		Description: "解析当前时间、相对日期和日期范围。遇到本月、今天、本周等相对日期时先调用本工具获得明确日期。",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"range": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"today", "yesterday", "tomorrow", "this_week", "last_week", "next_week", "this_month", "last_month", "next_month", "this_year", "custom"},
					"description": "要解析的日期范围。",
				},
				"timezone":   map[string]interface{}{"type": "string", "description": "可选时区，例如 Asia/Shanghai。"},
				"start_date": map[string]interface{}{"type": "string", "description": "custom 范围开始日期，格式 YYYY-MM-DD。"},
				"end_date":   map[string]interface{}{"type": "string", "description": "custom 范围结束日期，格式 YYYY-MM-DD。"},
			},
			"required": []string{"range"},
		},
	}
}
