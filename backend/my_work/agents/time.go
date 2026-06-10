package agents

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type TimeTool struct{}

type TimeRangeArgs struct {
	Range     string `json:"range"`
	Timezone  string `json:"timezone,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type TimeRangeResult struct {
	Ok           bool   `json:"ok"`
	Now          string `json:"now"`
	NowWeekday   string `json:"now_weekday"`
	StartDate    string `json:"start_date"`
	StartWeekday string `json:"start_weekday"`
	EndDate      string `json:"end_date"`
	EndWeekday   string `json:"end_weekday"`
	Label        string `json:"label"`
	Timezone     string `json:"timezone"`
}

func init() {
	Register(TimeTool{})
}

func (TimeTool) Name() string {
	return "time"
}

func (TimeTool) Enabled(config ToolConfig) bool {
	return config.Enabled
}

func (TimeTool) ShouldUse(messages []ChatMessage) bool {
	content := strings.ToLower(LastUserContent(messages))
	keywords := []string{
		"时间", "日期", "今天", "昨天", "明天", "本周", "这周", "上周", "下周", "本月", "这个月", "上月", "下月", "今年", "去年", "明年",
		"time", "date", "today", "yesterday", "tomorrow", "week", "month", "year", "now",
	}
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	return false
}

func (TimeTool) Enrich(ctx context.Context, messages []ChatMessage, config ToolConfig, trace TraceFunc) ([]ChatMessage, error) {
	select {
	case <-ctx.Done():
		return messages, ctx.Err()
	default:
	}

	now := time.Now()
	content := buildTimeContext(now)
	if trace != nil {
		trace("time", "execute", "success", "已获取当前时间上下文", map[string]interface{}{
			"now":   now.Format("2006-01-02 15:04:05"),
			"today": now.Format("2006-01-02"),
		})
	}

	enriched := append([]ChatMessage{}, messages...)
	enriched = append(enriched, SystemMessage(content))
	return enriched, nil
}

func ResolveTimeRange(args TimeRangeArgs, now time.Time) (TimeRangeResult, error) {
	loc := now.Location()
	if args.Timezone != "" {
		if loaded, err := time.LoadLocation(args.Timezone); err == nil {
			loc = loaded
			now = now.In(loc)
		}
	}

	todayStart := dateStart(now)
	weekStart := todayStart.AddDate(0, 0, -int((int(todayStart.Weekday())+6)%7))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)

	rangeName := strings.TrimSpace(strings.ToLower(args.Range))
	start := todayStart
	end := todayStart
	label := "今天"

	switch rangeName {
	case "", "today":
		label = "今天"
		start = todayStart
		end = todayStart
	case "yesterday":
		label = "昨天"
		start = todayStart.AddDate(0, 0, -1)
		end = start
	case "tomorrow":
		label = "明天"
		start = todayStart.AddDate(0, 0, 1)
		end = start
	case "this_week":
		label = "本周"
		start = weekStart
		end = weekStart.AddDate(0, 0, 6)
	case "last_week":
		label = "上周"
		start = weekStart.AddDate(0, 0, -7)
		end = weekStart.AddDate(0, 0, -1)
	case "next_week":
		label = "下周"
		start = weekStart.AddDate(0, 0, 7)
		end = weekStart.AddDate(0, 0, 13)
	case "this_month":
		label = "本月"
		start = monthStart
		end = monthStart.AddDate(0, 1, -1)
	case "last_month":
		label = "上月"
		start = monthStart.AddDate(0, -1, 0)
		end = monthStart.AddDate(0, 0, -1)
	case "next_month":
		label = "下月"
		start = monthStart.AddDate(0, 1, 0)
		end = monthStart.AddDate(0, 2, -1)
	case "this_year":
		label = "今年"
		start = yearStart
		end = yearStart.AddDate(1, 0, -1)
	case "custom":
		if args.StartDate == "" || args.EndDate == "" {
			return TimeRangeResult{Ok: false, Now: now.Format("2006-01-02 15:04:05"), Timezone: loc.String()}, fmt.Errorf("custom range requires start_date and end_date")
		}
		parsedStart, err := time.ParseInLocation("2006-01-02", args.StartDate, loc)
		if err != nil {
			return TimeRangeResult{Ok: false, Now: now.Format("2006-01-02 15:04:05"), Timezone: loc.String()}, fmt.Errorf("invalid start_date: %w", err)
		}
		parsedEnd, err := time.ParseInLocation("2006-01-02", args.EndDate, loc)
		if err != nil {
			return TimeRangeResult{Ok: false, Now: now.Format("2006-01-02 15:04:05"), Timezone: loc.String()}, fmt.Errorf("invalid end_date: %w", err)
		}
		if parsedEnd.Before(parsedStart) {
			return TimeRangeResult{Ok: false, Now: now.Format("2006-01-02 15:04:05"), Timezone: loc.String()}, fmt.Errorf("end_date must be after start_date")
		}
		label = "自定义"
		start = parsedStart
		end = parsedEnd
	default:
		return TimeRangeResult{Ok: false, Now: now.Format("2006-01-02 15:04:05"), Timezone: loc.String()}, fmt.Errorf("unsupported range: %s", args.Range)
	}

	return TimeRangeResult{
		Ok:           true,
		Now:          now.Format("2006-01-02 15:04:05"),
		NowWeekday:   weekdayName(now.Weekday()),
		StartDate:    start.Format("2006-01-02"),
		StartWeekday: weekdayName(start.Weekday()),
		EndDate:      end.Format("2006-01-02"),
		EndWeekday:   weekdayName(end.Weekday()),
		Label:        label,
		Timezone:     loc.String(),
	}, nil
}

func buildTimeContext(now time.Time) string {
	todayStart := dateStart(now)
	weekStart := todayStart.AddDate(0, 0, -int((int(todayStart.Weekday())+6)%7))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	return fmt.Sprintf(`以下是当前时间上下文，请在回答涉及相对日期/时间的问题时使用：
- 当前时间：%s
- 今天：%s，范围 %s 至 %s
- 昨天：%s
- 明天：%s
- 本周：%s 至 %s
- 本月：%s 至 %s
- 今年：%s 至 %s`,
		formatDateTimeWithWeek(now, "2006-01-02 15:04:05 MST"),
		formatDateWithWeek(todayStart), formatDateTimeWithWeek(todayStart, "2006-01-02 15:04:05"), formatDateTimeWithWeek(todayStart.AddDate(0, 0, 1).Add(-time.Second), "2006-01-02 15:04:05"),
		formatDateWithWeek(todayStart.AddDate(0, 0, -1)),
		formatDateWithWeek(todayStart.AddDate(0, 0, 1)),
		formatDateWithWeek(weekStart), formatDateWithWeek(weekStart.AddDate(0, 0, 7).Add(-time.Second)),
		formatDateWithWeek(monthStart), formatDateWithWeek(monthStart.AddDate(0, 1, 0).Add(-time.Second)),
		formatDateWithWeek(yearStart), formatDateWithWeek(yearStart.AddDate(1, 0, 0).Add(-time.Second)),
	)
}

func formatDateWithWeek(t time.Time) string {
	return formatDateTimeWithWeek(t, "2006-01-02")
}

func formatDateTimeWithWeek(t time.Time, layout string) string {
	return fmt.Sprintf("%s（%s）", t.Format(layout), weekdayName(t.Weekday()))
}

func weekdayName(weekday time.Weekday) string {
	names := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	return names[weekday]
}

func dateStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
