package agents

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type TimeTool struct{}

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
		now.Format("2006-01-02 15:04:05 MST"),
		todayStart.Format("2006-01-02"), todayStart.Format("2006-01-02 15:04:05"), todayStart.AddDate(0, 0, 1).Add(-time.Second).Format("2006-01-02 15:04:05"),
		todayStart.AddDate(0, 0, -1).Format("2006-01-02"),
		todayStart.AddDate(0, 0, 1).Format("2006-01-02"),
		weekStart.Format("2006-01-02"), weekStart.AddDate(0, 0, 7).Add(-time.Second).Format("2006-01-02"),
		monthStart.Format("2006-01-02"), monthStart.AddDate(0, 1, 0).Add(-time.Second).Format("2006-01-02"),
		yearStart.Format("2006-01-02"), yearStart.AddDate(1, 0, 0).Add(-time.Second).Format("2006-01-02"),
	)
}

func dateStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
