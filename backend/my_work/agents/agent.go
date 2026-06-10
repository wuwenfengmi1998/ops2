package agents

import (
	"context"
	"fmt"
	"strings"
)

type ToolConfig struct {
	Name        string
	Enabled     bool
	Description string
}

type TraceFunc func(tool string, stage string, status string, message string, data map[string]interface{})

type Tool interface {
	Name() string
	Enabled(config ToolConfig) bool
	ShouldUse(messages []ChatMessage) bool
	Enrich(ctx context.Context, messages []ChatMessage, config ToolConfig, trace TraceFunc) ([]ChatMessage, error)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var registry []Tool

func Register(tool Tool) {
	registry = append(registry, tool)
}

func EnabledTools(configs []ToolConfig) []Tool {
	tools := make([]Tool, 0)
	for _, tool := range registry {
		config := FindToolConfig(configs, tool.Name())
		if tool.Enabled(config) {
			tools = append(tools, tool)
		}
	}
	return tools
}

func EnrichMessages(ctx context.Context, messages []ChatMessage, configs []ToolConfig, trace TraceFunc) []ChatMessage {
	enriched := append([]ChatMessage{}, messages...)
	for _, tool := range EnabledTools(configs) {
		if !tool.ShouldUse(enriched) {
			continue
		}
		config := FindToolConfig(configs, tool.Name())
		var err error
		enriched, err = tool.Enrich(ctx, enriched, config, trace)
		if err != nil && trace != nil {
			trace(tool.Name(), "execute", "error", err.Error(), nil)
		}
	}
	return enriched
}

func FindToolConfig(configs []ToolConfig, name string) ToolConfig {
	for _, config := range configs {
		if strings.EqualFold(config.Name, name) {
			return config
		}
	}
	return ToolConfig{Name: name, Enabled: false}
}

func LastUserContent(messages []ChatMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i].Content
		}
	}
	return ""
}

func SystemMessage(content string) ChatMessage {
	return ChatMessage{Role: "system", Content: content}
}

func SafeString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
