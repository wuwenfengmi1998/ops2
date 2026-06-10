package routers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ops/agents"
	"ops/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SSE frame types sent to frontend
type sseEvent struct {
	Type    string                 `json:"type"`
	Text    string                 `json:"text,omitempty"`
	Tool    string                 `json:"tool,omitempty"`
	Stage   string                 `json:"stage,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Stats   *tokenUsageStats       `json:"stats,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

type tokenUsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// chatRequestFromFrontend is the expected POST body
type chatRequest struct {
	Messages   []chatMessage `json:"messages"`
	OpenAIName string        `json:"openaiName,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openaiChatRequest is the request sent to the upstream OpenAI-compatible API
type openaiChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openaiStreamChunk is one SSE data line from the upstream
type openaiStreamChunk struct {
	ID      string         `json:"id,omitempty"`
	Object  string         `json:"object,omitempty"`
	Created int64          `json:"created,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []openaiChoice `json:"choices"`
	Usage   *openaiUsage   `json:"usage,omitempty"`
}

type openaiChatResponse struct {
	Choices []openaiResponseChoice `json:"choices"`
	Usage   *openaiUsage           `json:"usage,omitempty"`
}

type openaiResponseChoice struct {
	Message openaiMessage `json:"message"`
}

type toolRouteResponse struct {
	Tools []struct {
		Name   string `json:"name"`
		Reason string `json:"reason"`
	} `json:"tools"`
	Reason string `json:"reason"`
}

type openaiChoice struct {
	Index  int         `json:"index"`
	Delta  openaiDelta `json:"delta"`
	Finish *string     `json:"finish_reason,omitempty"`
}

type openaiDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ApiAIChat(r *gin.RouterGroup) {
	r.GET("/openai", handleOpenAIProfiles)
	r.POST("/chat", handleChat)

	admin := r.Group("/admin")
	admin.POST("/config", handleAIChatAdminGetConfig)
	admin.POST("/config/update", handleAIChatAdminUpdateConfig)
	admin.POST("/refresh", handleAIChatAdminRefreshCache)
}

func handleOpenAIProfiles(ctx *gin.Context) {
	cfg := getAIChatConfig()
	active := ""
	profiles := make([]map[string]interface{}, 0, len(cfg.OpenAI))
	for _, profile := range cfg.OpenAI {
		if profile.Active {
			active = profile.Name
		}
		profiles = append(profiles, map[string]interface{}{
			"name":      profile.Name,
			"active":    profile.Active,
			"baseUrl":   profile.BaseUrl,
			"model":     profile.Model,
			"timeout":   profile.Timeout,
			"maxTokens": profile.MaxTokens,
		})
	}
	ReturnJson(ctx, "apiOK", gin.H{
		"enabled":  cfg.Enabled,
		"active":   active,
		"profiles": profiles,
		"toolRouter": gin.H{
			"enabled":    cfg.ToolRouter.Enabled,
			"openaiName": cfg.ToolRouter.OpenAIName,
			"timeout":    cfg.ToolRouter.Timeout,
			"maxTokens":  cfg.ToolRouter.MaxTokens,
		},
	})
}

func handleChat(ctx *gin.Context) {
	data, _ := SeparateData(ctx)

	if data == nil {
		sendSSEError(ctx, "请求数据为空")
		return
	}

	var req chatRequest
	if err := decodeJSON(data, &req); err != nil {
		sendSSEError(ctx, "解析消息失败: "+err.Error())
		return
	}

	if len(req.Messages) == 0 {
		sendSSEError(ctx, "消息不能为空")
		return
	}

	// Check ai config
	cfg := getAIChatConfig()
	profile, ok := selectOpenAIProfile(cfg, req.OpenAIName)
	if !cfg.Enabled || !ok || profile.Model == "" || profile.ApiKey == "" {
		sendSSEError(ctx, "AI 聊天未配置，请在后台配置 API Key 和模型")
		return
	}
	toolRouterProfile, hasToolRouterProfile := selectOpenAIProfile(cfg, cfg.ToolRouter.OpenAIName)

	// Convert to agent messages and enrich with tools
	chatMsgs := convertToChatMessages(req.Messages)
	toolConfigs := []agents.ToolConfig{}
	if cfg.ToolRouter.Enabled {
		toolConfigs = buildToolConfigs(cfg.ToolRouter.Tools)
		if hasToolRouterProfile && toolRouterProfile.Model != "" && toolRouterProfile.ApiKey != "" {
			selected, err := routeTools(ctx.Request.Context(), toolRouterProfile, cfg.ToolRouter, chatMsgs)
			if err == nil && selected != nil {
				toolConfigs = filterToolConfigs(toolConfigs, selected)
			}
		}
	}

	// Set up SSE headers
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.WriteHeader(http.StatusOK)
	flusher, _ := ctx.Writer.(http.Flusher)

	// Enrich messages with tools (pre-process)
	chatMsgs = agents.EnrichMessages(ctx.Request.Context(), chatMsgs, toolConfigs, func(tool, stage, status, message string, data map[string]interface{}) {
		sendSSE(ctx, flusher, sseEvent{
			Type:    "trace",
			Tool:    tool,
			Stage:   stage,
			Status:  status,
			Message: message,
			Data:    data,
		})
	})

	// Build OpenAI-compatible request
	openaiMsgs := convertToOpenAIMessages(chatMsgs)
	apiReq := openaiChatRequest{
		Model:       profile.Model,
		Messages:    openaiMsgs,
		Stream:      true,
		MaxTokens:   profile.MaxTokens,
		Temperature: 0.7,
	}

	// Add system prompt if configured
	if profile.SystemPrompt != "" {
		apiReq.Messages = append([]openaiMessage{{Role: "system", Content: profile.SystemPrompt}}, apiReq.Messages...)
	}

	err := streamOpenAI(ctx.Request.Context(), profile, apiReq, func(chunk openaiStreamChunk) {
		for _, choice := range chunk.Choices {
			if choice.Delta.Content != "" {
				sendSSE(ctx, flusher, sseEvent{
					Type: "delta",
					Text: choice.Delta.Content,
				})
			}
		}
		if chunk.Usage != nil {
			sendSSE(ctx, flusher, sseEvent{
				Type: "stats",
				Stats: &tokenUsageStats{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:      chunk.Usage.TotalTokens,
				},
			})
		}
	})
	if err != nil {
		sendSSE(ctx, flusher, sseEvent{Type: "error", Error: "请求失败: " + err.Error()})
		sendSSEDone(ctx, flusher)
		return
	}

	sendSSEDone(ctx, flusher)
	flusher.Flush()
}

func streamOpenAI(ctx context.Context, cfg models.ConfigsAIChatOpenAI_, req openaiChatRequest, onData func(openaiStreamChunk)) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	url := strings.TrimRight(cfg.BaseUrl, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("连接上游服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上游返回 %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 256*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		payload := strings.TrimPrefix(line, "data: ")
		payload = strings.TrimSpace(payload)

		if payload == "[DONE]" {
			continue
		}

		var chunk openaiStreamChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			continue
		}

		onData(chunk)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流失败: %w", err)
	}

	return nil
}

// Initialize with system prompt if present
func sendSSE(ctx *gin.Context, flusher http.Flusher, event sseEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(ctx.Writer, "data: %s\n\n", string(data))
	flusher.Flush()
}

func sendSSEDone(ctx *gin.Context, flusher http.Flusher) {
	_, _ = fmt.Fprint(ctx.Writer, "data: [DONE]\n\n")
	flusher.Flush()
}

func sendSSEError(ctx *gin.Context, message string) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.WriteHeader(http.StatusOK)
	flusher, _ := ctx.Writer.(http.Flusher)

	sendSSE(ctx, flusher, sseEvent{
		Type:  "error",
		Error: message,
	})

	sendSSEDone(ctx, flusher)
	flusher.Flush()
}

func convertToChatMessages(msgs []chatMessage) []agents.ChatMessage {
	result := make([]agents.ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, agents.ChatMessage{Role: m.Role, Content: m.Content})
	}
	return result
}

func convertToOpenAIMessages(msgs []agents.ChatMessage) []openaiMessage {
	result := make([]openaiMessage, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, openaiMessage{Role: m.Role, Content: m.Content})
	}
	return result
}

func buildToolConfigs(configs []models.ConfigsAIChatTool_) []agents.ToolConfig {
	result := make([]agents.ToolConfig, 0, len(configs))
	for _, c := range configs {
		result = append(result, agents.ToolConfig{
			Name:        c.Name,
			Enabled:     c.Enabled,
			Description: c.Description,
		})
	}
	return result
}

func selectOpenAIProfile(cfg models.ConfigsAIChat_, name string) (models.ConfigsAIChatOpenAI_, bool) {
	if name != "" {
		for _, p := range cfg.OpenAI {
			if p.Name == name {
				return p, true
			}
		}
		return models.ConfigsAIChatOpenAI_{}, false
	}
	for _, p := range cfg.OpenAI {
		if p.Active {
			return p, true
		}
	}
	if len(cfg.OpenAI) > 0 {
		return cfg.OpenAI[0], true
	}
	return models.ConfigsAIChatOpenAI_{}, false
}

func routeTools(ctx context.Context, profile models.ConfigsAIChatOpenAI_, router models.ConfigsAIChatToolRouter_, messages []agents.ChatMessage) ([]string, error) {
	openaiMsgs := []openaiMessage{}
	lastUserContent := agents.LastUserContent(messages)
	if lastUserContent != "" {
		openaiMsgs = append(openaiMsgs, openaiMessage{Role: "user", Content: lastUserContent})
	}

	toolNames := make([]string, 0, len(router.Tools))
	for _, t := range router.Tools {
		if t.Enabled {
			toolNames = append(toolNames, t.Name+" - "+t.Description)
		}
	}
	if len(toolNames) == 0 {
		return nil, nil
	}

	sysPrompt := "请根据用户的最新一条消息，判断需要启用哪些工具来完成用户需求。\n可选工具：\n" + strings.Join(toolNames, "\n") + "\n\n回复格式要求：\n```json\n{\"tools\":[{\"name\":\"工具名称\",\"reason\":\"选择原因\"}],\"reason\":\"整体判断理由\"}\n```\n仅输出 JSON 代码块。如果没有需要启用的工具，返回 {\"tools\":[]}。"
	openaiMsgs = append([]openaiMessage{{Role: "system", Content: sysPrompt}}, openaiMsgs...)

	req := openaiChatRequest{
		Model:       profile.Model,
		Messages:    openaiMsgs,
		Stream:      false,
		MaxTokens:   router.MaxTokens,
		Temperature: 0.1,
	}

	bodyBytes, _ := json.Marshal(req)
	url := strings.TrimRight(profile.BaseUrl, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+profile.ApiKey)
	client := &http.Client{Timeout: time.Duration(router.Timeout) * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result openaiChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, nil
	}

	response := result.Choices[0].Message.Content
	toolRouteResponse := extractToolsFromResponse(response)
	return toolRouteResponse, nil
}

func extractToolsFromResponse(response string) []string {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end <= start {
		return nil
	}
	var parsed toolRouteResponse
	if err := json.Unmarshal([]byte(response[start:end+1]), &parsed); err != nil {
		return nil
	}
	tools := make([]string, 0, len(parsed.Tools))
	for _, t := range parsed.Tools {
		tools = append(tools, t.Name)
	}
	return tools
}

func filterToolConfigs(configs []agents.ToolConfig, selected []string) []agents.ToolConfig {
	if len(selected) == 0 {
		return []agents.ToolConfig{}
	}
	selectedSet := make(map[string]bool, len(selected))
	for _, s := range selected {
		selectedSet[s] = true
	}
	filtered := make([]agents.ToolConfig, 0, len(configs))
	for _, c := range configs {
		if selectedSet[c.Name] {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
