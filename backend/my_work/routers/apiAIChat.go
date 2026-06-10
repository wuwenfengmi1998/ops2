package routers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	PromptTokens               int     `json:"prompt_tokens"`
	CompletionTokens           int     `json:"completion_tokens"`
	ToolPromptTokens           int     `json:"tool_prompt_tokens"`
	ToolCompletionTokens       int     `json:"tool_completion_tokens"`
	TotalTokens                int     `json:"total_tokens"`
	CompletionTokensPerSec     float64 `json:"completion_tokens_per_sec"`
	PeakCompletionTokensPerSec float64 `json:"peak_completion_tokens_per_sec"`
	Estimated                  bool    `json:"estimated"`
}

const maxImageDataSize = 4 * 1024 * 1024

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

// chatRequestFromFrontend is the expected POST body
type chatRequest struct {
	Messages       []chatMessage `json:"messages"`
	OpenAIName     string        `json:"openaiName,omitempty"`
	ConversationID uint          `json:"conversationId,omitempty"`
	ClientLocalID  string        `json:"clientLocalId,omitempty"`
	SaveToServer   bool          `json:"saveToServer,omitempty"`
}

type chatMessage struct {
	Role          string `json:"role"`
	Content       string `json:"content"`
	ImageURL      string `json:"image_url,omitempty"`
	ImageURLAlias string `json:"imageURL,omitempty"`
}

// openaiChatRequest is the request sent to the upstream OpenAI-compatible API
type openaiChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Tools       []openaiTool    `json:"tools,omitempty"`
	ToolChoice  any             `json:"tool_choice,omitempty"`
}

type openaiMessage struct {
	Role       string           `json:"role"`
	Content    any              `json:"content,omitempty"`
	Name       string           `json:"name,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
	ToolCalls  []openaiToolCall `json:"tool_calls,omitempty"`
}

type openaiTool struct {
	Type     string                   `json:"type"`
	Function openaiFunctionDefinition `json:"function"`
}

type openaiFunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type openaiToolCall struct {
	ID       string             `json:"id,omitempty"`
	Type     string             `json:"type,omitempty"`
	Function openaiFunctionCall `json:"function"`
}

type openaiFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openaiContentPart struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *openaiImageURL `json:"image_url,omitempty"`
}

type openaiImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type openaiResponseMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content"`
	ToolCalls []openaiToolCall `json:"tool_calls,omitempty"`
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
	Message openaiResponseMessage `json:"message"`
}

type toolSelection struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type toolRoutingDecision struct {
	Tools  []toolSelection `json:"tools"`
	Reason string          `json:"reason"`
}

type toolRoutingResult struct {
	Decision toolRoutingDecision
	Selected []string
	Messages []openaiMessage
	Response string
	Usage    *openaiUsage
}

type openaiChoice struct {
	Index  int         `json:"index"`
	Delta  openaiDelta `json:"delta"`
	Finish *string     `json:"finish_reason,omitempty"`
}

type openaiDelta struct {
	Role             string           `json:"role,omitempty"`
	Content          string           `json:"content,omitempty"`
	ReasoningContent string           `json:"reasoning_content,omitempty"`
	Reasoning        string           `json:"reasoning,omitempty"`
	Thinking         string           `json:"thinking,omitempty"`
	ToolCalls        []openaiToolCall `json:"tool_calls,omitempty"`
}

type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ApiAIChat(r *gin.RouterGroup) {
	r.GET("/openai", handleOpenAIProfiles)
	r.POST("/chat", handleChat)

	conversations := r.Group("/conversations")
	conversations.POST("/list", handleAIChatConversationList)
	conversations.POST("/get", handleAIChatConversationGet)
	conversations.POST("/update", handleAIChatConversationUpdate)
	conversations.POST("/delete", handleAIChatConversationDelete)

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
			"name":                profile.Name,
			"active":              profile.Active,
			"baseUrl":             profile.BaseUrl,
			"model":               profile.Model,
			"timeout":             profile.Timeout,
			"maxTokens":           profile.MaxTokens,
			"contextWindowTokens": profile.ContextWindowTokens,
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
	data, cookieValue := SeparateData(ctx)

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

	var currentUser *TabUser
	if cookieValue != "" {
		if user, err := AuthenticationAuthorityFromCookie(cookieValue); err == nil {
			currentUser = user
		}
	}

	// Check ai config
	cfg := getAIChatConfig()
	profile, ok := selectOpenAIProfile(cfg, req.OpenAIName)
	if !cfg.Enabled || !ok || profile.Model == "" || profile.ApiKey == "" {
		sendSSEError(ctx, "AI 聊天未配置，请在后台配置 API Key 和模型")
		return
	}
	chatMsgs := convertToChatMessages(req.Messages)

	// Set up SSE headers before routing/tools so progress can stream immediately.
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Writer.WriteHeader(http.StatusOK)
	flusher, _ := ctx.Writer.(http.Flusher)
	tracker := newTokenUsageTracker()

	traceEvents := []sseEvent{}
	emitTrace := func(tool, stage, status, message string, data map[string]interface{}) {
		event := sseEvent{Type: "trace", Tool: tool, Stage: stage, Status: status, Message: message, Data: data}
		traceEvents = append(traceEvents, event)
		sendSSE(ctx, flusher, event)
	}
	emitStats := func(stats tokenUsageStats) {
		sendSSE(ctx, flusher, sseEvent{Type: "stats", Stats: &stats})
	}

	conversation, persistErr := prepareAIChatPersistence(currentUser, req, profile.Name)
	if persistErr != nil {
		emitTrace("chat", "persist", "error", "聊天保存失败，将继续仅本次对话", map[string]interface{}{"error": persistErr.Error()})
		conversation = nil
	} else if conversation != nil {
		sendSSE(ctx, flusher, sseEvent{Type: "conversation", Data: map[string]interface{}{
			"id":            conversation.ID,
			"title":         conversation.Title,
			"clientLocalId": conversation.ClientLocalID,
		}})
	}

	toolConfigs := []agents.ToolConfig{}
	if cfg.ToolRouter.Enabled {
		toolConfigs = buildToolConfigs(cfg.ToolRouter.Tools)
	}

	// Build OpenAI-compatible request
	openaiMsgs, err := convertToOpenAIMessages(chatMsgs)
	if err != nil {
		sendSSE(ctx, flusher, sseEvent{Type: "error", Error: err.Error()})
		sendSSEDone(ctx, flusher)
		return
	}
	functionTools := buildFunctionTools(toolConfigs)
	if profile.SystemPrompt != "" {
		openaiMsgs = append([]openaiMessage{{Role: "system", Content: profile.SystemPrompt}}, openaiMsgs...)
	}
	if len(functionTools) > 0 {
		toolNames := make([]string, 0, len(functionTools))
		for _, tool := range functionTools {
			toolNames = append(toolNames, tool.Function.Name)
		}
		emitTrace("function_tools", "prepare", "success", "已启用 Function Calling 工具", map[string]interface{}{"tools": toolNames})
		openaiMsgs = append([]openaiMessage{{Role: "system", Content: "可用工具使用规则：当用户询问“我是谁”“当前登录用户是谁”“我的用户信息”等当前身份问题时，调用 ops_ai_assistant_current_user；工具返回 loggedIn=true 时按工具结果回答当前用户信息，返回 loggedIn=false 时说明不知道并提示需要登录才能获取信息。当用户询问本月、今天、本周、下周等相对日期的日程时，先调用 time 获取明确 start_date/end_date，再调用 ops_ai_assistant_schedule_query 查询日程。不要臆造工具结果中不存在的信息。"}}, openaiMsgs...)
		var toolExecuted bool
		openaiMsgs, toolExecuted, err = runOpenAIToolLoop(ctx.Request.Context(), profile, openaiMsgs, functionTools, currentUser, tracker, emitTrace)
		if err != nil {
			emitTrace("model", "tool_call", "error", "工具调用失败，将继续普通回答", map[string]interface{}{"error": err.Error()})
		} else if toolExecuted {
			emitTrace("model", "tool_call", "success", "工具调用完成，准备生成最终回答", nil)
		}
	}
	apiReq := openaiChatRequest{
		Model:       profile.Model,
		Messages:    openaiMsgs,
		Stream:      true,
		MaxTokens:   profile.MaxTokens,
		Temperature: 0.7,
	}

	trimmedMessages, trimStats := trimOpenAIMessagesToContextWindow(apiReq.Messages, profile.ContextWindowTokens)
	apiReq.Messages = trimmedMessages
	if trimStats.RemovedMessages > 0 {
		emitTrace("model", "context_window", "success", "上下文窗口已裁剪旧消息", map[string]interface{}{
			"limit":            trimStats.Limit,
			"before_tokens":    trimStats.BeforeTokens,
			"after_tokens":     trimStats.AfterTokens,
			"removed_messages": trimStats.RemovedMessages,
		})
	}

	modelPromptTokens := estimateOpenAIMessagesTokens(apiReq.Messages)
	completionTokens := 0
	modelUsageReceived := false
	assistantContent := strings.Builder{}
	reasoningContent := strings.Builder{}
	var finalStats *tokenUsageStats
	streamStarted := time.Now()
	windowStarted := streamStarted
	windowTokens := 0
	peakTokensPerSecond := 0.0
	emitTrace("model", "stream", "running", "正在请求模型回复", nil)

	err = streamOpenAI(ctx.Request.Context(), profile, apiReq, func(chunk openaiStreamChunk) {
		for _, choice := range chunk.Choices {
			reasoningText := choice.Delta.ReasoningContent
			if reasoningText == "" {
				reasoningText = choice.Delta.Reasoning
			}
			if reasoningText == "" {
				reasoningText = choice.Delta.Thinking
			}
			if reasoningText != "" {
				reasoningContent.WriteString(reasoningText)
				sendSSE(ctx, flusher, sseEvent{Type: "reasoning", Text: reasoningText})
			}

			if choice.Delta.Content != "" {
				assistantContent.WriteString(choice.Delta.Content)
				deltaTokens := estimateTokenCount(choice.Delta.Content)
				completionTokens += deltaTokens
				windowTokens += deltaTokens
				elapsedWindow := time.Since(windowStarted).Seconds()
				if elapsedWindow >= 1 {
					peakTokensPerSecond = maxFloat(peakTokensPerSecond, float64(windowTokens)/elapsedWindow)
					windowStarted = time.Now()
					windowTokens = 0
				} else if peakTokensPerSecond == 0 && elapsedWindow > 0.25 {
					peakTokensPerSecond = maxFloat(peakTokensPerSecond, float64(windowTokens)/elapsedWindow)
				}
				stats := tracker.setModelEstimate(modelPromptTokens, completionTokens).snapshot(tokensPerSecond(completionTokens, streamStarted), peakTokensPerSecond)
				finalStats = &stats
				sendSSE(ctx, flusher, sseEvent{Type: "delta", Text: choice.Delta.Content, Stats: &stats})
			}
		}
		if chunk.Usage != nil {
			modelUsageReceived = true
			stats := tracker.setModelUsage(chunk.Usage).snapshot(tokensPerSecond(tracker.completionTokens, streamStarted), peakTokensPerSecond)
			finalStats = &stats
			emitStats(stats)
		}
	})
	if err != nil {
		errorText := "请求失败: " + err.Error()
		if conversation != nil && assistantContent.Len() > 0 {
			_ = persistAIChatAssistantMessage(conversation, assistantContent.String(), reasoningContent.String(), traceEvents, finalStats)
		}
		sendSSE(ctx, flusher, sseEvent{Type: "error", Error: errorText})
		sendSSEDone(ctx, flusher)
		return
	}

	if windowTokens > 0 {
		elapsedWindow := time.Since(windowStarted).Seconds()
		if elapsedWindow > 0 {
			peakTokensPerSecond = maxFloat(peakTokensPerSecond, float64(windowTokens)/elapsedWindow)
		}
	}
	emitTrace("model", "stream", "success", "模型回复完成", nil)
	if modelUsageReceived {
		stats := tracker.snapshot(tokensPerSecond(tracker.completionTokens, streamStarted), peakTokensPerSecond)
		finalStats = &stats
		emitStats(stats)
	} else {
		stats := tracker.setModelEstimate(modelPromptTokens, completionTokens).snapshot(tokensPerSecond(completionTokens, streamStarted), peakTokensPerSecond)
		finalStats = &stats
		emitStats(stats)
	}
	if err := persistAIChatAssistantMessage(conversation, assistantContent.String(), reasoningContent.String(), traceEvents, finalStats); err != nil {
		sendSSE(ctx, flusher, sseEvent{Type: "trace", Tool: "chat", Stage: "persist", Status: "error", Message: "助手回复保存失败", Data: map[string]interface{}{"error": err.Error()}})
	}
	sendSSEDone(ctx, flusher)
	flusher.Flush()
}

func callOpenAIChat(ctx context.Context, cfg models.ConfigsAIChatOpenAI_, req openaiChatRequest) (*openaiChatResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := strings.TrimRight(cfg.BaseUrl, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.ApiKey)

	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("连接上游服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("上游返回 %d: %s", resp.StatusCode, string(body))
	}

	var result openaiChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &result, nil
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
		imageURL := m.ImageURL
		if imageURL == "" {
			imageURL = m.ImageURLAlias
		}
		result = append(result, agents.ChatMessage{Role: m.Role, Content: m.Content, ImageURL: imageURL})
	}
	return result
}

func convertToOpenAIMessages(msgs []agents.ChatMessage) ([]openaiMessage, error) {
	result := make([]openaiMessage, 0, len(msgs))
	for _, m := range msgs {
		content, err := buildOpenAIContent(m)
		if err != nil {
			return nil, err
		}
		result = append(result, openaiMessage{Role: m.Role, Content: content})
	}
	return result, nil
}

func buildOpenAIContent(m agents.ChatMessage) (any, error) {
	if strings.TrimSpace(m.ImageURL) == "" {
		return m.Content, nil
	}

	imageURL, err := normalizeImageURL(m.ImageURL)
	if err != nil {
		return nil, err
	}

	parts := []openaiContentPart{
		{
			Type: "image_url",
			ImageURL: &openaiImageURL{
				URL:    imageURL,
				Detail: "auto",
			},
		},
	}
	if m.Content != "" {
		parts = append(parts, openaiContentPart{Type: "text", Text: m.Content})
	}
	return parts, nil
}

func normalizeImageURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("图片地址为空")
	}

	if strings.HasPrefix(strings.ToLower(value), "data:") {
		return normalizeImageDataURI(value)
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", errors.New("图片地址无效，仅支持 http/https URL 或 base64 data URI")
	}
	return value, nil
}

func normalizeImageDataURI(raw string) (string, error) {
	commaIndex := strings.Index(raw, ",")
	if commaIndex == -1 {
		return "", errors.New("图片 data URI 格式无效")
	}

	metadata := strings.TrimSpace(raw[len("data:"):commaIndex])
	payload := strings.TrimSpace(raw[commaIndex+1:])
	if metadata == "" || payload == "" {
		return "", errors.New("图片 data URI 格式无效")
	}

	metadataParts := strings.Split(metadata, ";")
	mimeType := strings.ToLower(strings.TrimSpace(metadataParts[0]))
	if !allowedImageTypes[mimeType] {
		return "", errors.New("图片格式不支持，仅支持 jpeg/png/webp/gif")
	}

	hasBase64 := false
	for _, part := range metadataParts[1:] {
		if strings.EqualFold(strings.TrimSpace(part), "base64") {
			hasBase64 = true
			break
		}
	}
	if !hasBase64 {
		return "", errors.New("图片 data URI 必须使用 base64 编码")
	}

	if len(payload) > maxImageDataSize*4/3+16 {
		return "", errors.New("图片过大，请选择小于 4MB 的图片")
	}
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", errors.New("图片 base64 数据无效")
	}
	if len(decoded) > maxImageDataSize {
		return "", errors.New("图片过大，请选择小于 4MB 的图片")
	}

	return "data:" + mimeType + ";base64," + payload, nil
}

type contextWindowTrimStats struct {
	Enabled         bool
	Limit           int
	BeforeTokens    int
	AfterTokens     int
	RemovedMessages int
}

func trimOpenAIMessagesToContextWindow(messages []openaiMessage, maxTokens int) ([]openaiMessage, contextWindowTrimStats) {
	stats := contextWindowTrimStats{Enabled: maxTokens > 0, Limit: maxTokens}
	if maxTokens <= 0 || len(messages) == 0 {
		stats.BeforeTokens = estimateOpenAIMessagesTokens(messages)
		stats.AfterTokens = stats.BeforeTokens
		return messages, stats
	}

	result := append([]openaiMessage(nil), messages...)
	stats.BeforeTokens = estimateOpenAIMessagesTokens(result)
	stats.AfterTokens = stats.BeforeTokens
	if stats.BeforeTokens <= maxTokens {
		return result, stats
	}

	for stats.AfterTokens > maxTokens {
		startIndex := 0
		if len(result) > 0 && result[0].Role == "system" {
			startIndex = 1
		}
		latestUserIndex := latestUserMessageIndex(result)
		removeIndex := -1
		for i := startIndex; i < len(result); i++ {
			if i == latestUserIndex {
				continue
			}
			if result[i].Role == "system" {
				continue
			}
			removeIndex = i
			break
		}
		if removeIndex == -1 {
			break
		}

		removeCount := 1
		if result[removeIndex].Role == "user" {
			nextIndex := removeIndex + 1
			if nextIndex < len(result) && nextIndex != latestUserIndex && result[nextIndex].Role == "assistant" {
				removeCount = 2
			}
		}
		result = append(result[:removeIndex], result[removeIndex+removeCount:]...)
		stats.RemovedMessages += removeCount
		stats.AfterTokens = estimateOpenAIMessagesTokens(result)
	}

	return result, stats
}

func latestUserMessageIndex(messages []openaiMessage) int {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return i
		}
	}
	return -1
}

type tokenUsageTracker struct {
	promptTokens         int
	completionTokens     int
	toolPromptTokens     int
	toolCompletionTokens int
	estimated            bool
}

func newTokenUsageTracker() *tokenUsageTracker {
	return &tokenUsageTracker{estimated: true}
}

func (t *tokenUsageTracker) addToolUsage(usage *openaiUsage, estimatedPromptTokens, estimatedCompletionTokens int) {
	if usage != nil {
		t.toolPromptTokens += usage.PromptTokens
		t.toolCompletionTokens += usage.CompletionTokens
		return
	}
	t.toolPromptTokens += estimatedPromptTokens
	t.toolCompletionTokens += estimatedCompletionTokens
	t.estimated = true
}

func (t *tokenUsageTracker) setModelEstimate(promptTokens, completionTokens int) *tokenUsageTracker {
	t.promptTokens = promptTokens
	t.completionTokens = completionTokens
	t.estimated = true
	return t
}

func (t *tokenUsageTracker) setModelUsage(usage *openaiUsage) *tokenUsageTracker {
	if usage == nil {
		return t
	}
	t.promptTokens = usage.PromptTokens
	t.completionTokens = usage.CompletionTokens
	return t
}

func (t *tokenUsageTracker) snapshot(completionTokensPerSec, peakCompletionTokensPerSec float64) tokenUsageStats {
	totalTokens := t.promptTokens + t.completionTokens + t.toolPromptTokens + t.toolCompletionTokens
	return tokenUsageStats{
		PromptTokens:               t.promptTokens,
		CompletionTokens:           t.completionTokens,
		ToolPromptTokens:           t.toolPromptTokens,
		ToolCompletionTokens:       t.toolCompletionTokens,
		TotalTokens:                totalTokens,
		CompletionTokensPerSec:     completionTokensPerSec,
		PeakCompletionTokensPerSec: peakCompletionTokensPerSec,
		Estimated:                  t.estimated,
	}
}

func estimateOpenAIMessagesTokens(messages []openaiMessage) int {
	total := 0
	for _, message := range messages {
		total += estimateTokenCount(message.Role) + 4
		total += estimateOpenAIContentTokens(message.Content)
	}
	return total
}

func estimateOpenAIContentTokens(content any) int {
	switch value := content.(type) {
	case string:
		return estimateTokenCount(value)
	case []openaiContentPart:
		total := 0
		for _, part := range value {
			switch part.Type {
			case "text":
				total += estimateTokenCount(part.Text)
			case "image_url":
				total += 85
			}
		}
		return total
	case []interface{}:
		data, err := json.Marshal(value)
		if err != nil {
			return 0
		}
		return estimateTokenCount(string(data))
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return 0
		}
		return estimateTokenCount(string(data))
	}
}

func estimateTokenCount(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	tokens := 0
	asciiRunes := 0
	flushASCII := func() {
		if asciiRunes > 0 {
			tokens += (asciiRunes + 3) / 4
			asciiRunes = 0
		}
	}

	for _, r := range text {
		if r <= 127 {
			if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
				flushASCII()
				continue
			}
			asciiRunes++
			continue
		}
		flushASCII()
		tokens++
	}
	flushASCII()
	if tokens == 0 {
		return 1
	}
	return tokens
}

func tokensPerSecond(tokens int, start time.Time) float64 {
	elapsed := time.Since(start).Seconds()
	if tokens <= 0 || elapsed <= 0 {
		return 0
	}
	return float64(tokens) / elapsed
}

func maxFloat(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
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

func buildFunctionTools(configs []agents.ToolConfig) []openaiTool {
	schemas := agents.FunctionToolSchemas(configs)
	tools := make([]openaiTool, 0, len(schemas))
	for _, schema := range schemas {
		tools = append(tools, openaiTool{
			Type: "function",
			Function: openaiFunctionDefinition{
				Name:        schema.Name,
				Description: schema.Description,
				Parameters:  schema.Parameters,
			},
		})
	}
	return tools
}

func parseJSONTraceValue(raw string) interface{} {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return value
	}
	return parsed
}

func runOpenAIToolLoop(ctx context.Context, profile models.ConfigsAIChatOpenAI_, messages []openaiMessage, tools []openaiTool, currentUser *TabUser, tracker *tokenUsageTracker, trace agents.TraceFunc) ([]openaiMessage, bool, error) {
	toolExecuted := false
	for round := 0; round < 5; round++ {
		if trace != nil {
			trace("model", "tool_call", "running", "正在请求模型决定是否调用工具", map[string]interface{}{"round": round + 1})
		}
		req := openaiChatRequest{
			Model:       profile.Model,
			Messages:    messages,
			Stream:      false,
			MaxTokens:   profile.MaxTokens,
			Temperature: 0.1,
			Tools:       tools,
			ToolChoice:  "auto",
		}
		resp, err := callOpenAIChat(ctx, profile, req)
		if err != nil {
			return messages, toolExecuted, err
		}
		responseText := ""
		if len(resp.Choices) == 0 {
			return messages, toolExecuted, nil
		}
		message := resp.Choices[0].Message
		responseText = message.Content
		tracker.addToolUsage(resp.Usage, estimateOpenAIMessagesTokens(messages), estimateTokenCount(responseText))
		if len(message.ToolCalls) == 0 {
			return messages, toolExecuted, nil
		}

		toolExecuted = true
		messages = append(messages, openaiMessage{Role: "assistant", Content: message.Content, ToolCalls: message.ToolCalls})
		for _, toolCall := range message.ToolCalls {
			toolName := strings.TrimSpace(toolCall.Function.Name)
			parsedArgs := parseJSONTraceValue(toolCall.Function.Arguments)
			if trace != nil {
				trace(toolName, "call", "running", "模型调用工具："+toolName, map[string]interface{}{
					"tool":      toolName,
					"arguments": parsedArgs,
				})
			}
			resultJSON, err := executeAIFunctionTool(ctx, toolName, []byte(toolCall.Function.Arguments), currentUser)
			status := "success"
			if err != nil {
				status = "error"
				resultJSON, _ = json.Marshal(map[string]interface{}{"ok": false, "error": err.Error()})
			}
			if trace != nil {
				data := map[string]interface{}{
					"tool":   toolName,
					"result": parseJSONTraceValue(string(resultJSON)),
				}
				if len(resultJSON) > 1200 {
					data["result"] = string(resultJSON[:1200]) + "..."
					data["truncated"] = true
				}
				trace(toolName, "execute", status, "工具执行完成："+toolName, data)
			}
			messages = append(messages, openaiMessage{Role: "tool", ToolCallID: toolCall.ID, Name: toolName, Content: string(resultJSON)})
		}
	}
	return messages, toolExecuted, fmt.Errorf("工具调用超过最大轮数")
}

func executeAIFunctionTool(ctx context.Context, name string, rawArgs []byte, currentUser *TabUser) ([]byte, error) {
	runtime := agents.FunctionToolRuntime{}
	if currentUser != nil {
		runtime.UserID = currentUser.ID
		runtime.UserName = currentUser.Name
		runtime.UserEmail = currentUser.Email
		runtime.UserType = currentUser.Type
		if userInfo := GetUserInfoFromUserID(currentUser.ID); userInfo != nil {
			runtime.UserInfo = &agents.CurrentUserInfo{
				ID:         userInfo.ID,
				UserID:     userInfo.UserID,
				FirstName:  userInfo.FirstName,
				Username:   userInfo.Username,
				Birthdate:  userInfo.Birthdate.Format("2006-01-02"),
				Gender:     userInfo.Gender,
				AvatarPath: userInfo.AvatarPath,
				Region:     userInfo.Region,
				Language:   userInfo.Language,
			}
		}
	}
	return agents.ExecuteFunctionTool(ctx, runtime, name, rawArgs)
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

func routeTools(ctx context.Context, profile models.ConfigsAIChatOpenAI_, router models.ConfigsAIChatToolRouter_, messages []agents.ChatMessage) (*toolRoutingResult, error) {
	lastUserContent := strings.TrimSpace(agents.LastUserContent(messages))
	if lastUserContent == "" {
		return nil, nil
	}

	openaiMsgs := []openaiMessage{{Role: "user", Content: lastUserContent}}

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

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("工具路由返回 %d: %s", resp.StatusCode, string(body))
	}

	var result openaiChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return &toolRoutingResult{Messages: openaiMsgs, Usage: result.Usage}, nil
	}

	response := result.Choices[0].Message.Content
	decision := extractToolRoutingDecision(response)
	selected := make([]string, 0, len(decision.Tools))
	for _, t := range decision.Tools {
		name := strings.TrimSpace(t.Name)
		if name != "" {
			selected = append(selected, name)
		}
	}
	return &toolRoutingResult{Decision: decision, Selected: selected, Messages: openaiMsgs, Response: response, Usage: result.Usage}, nil
}

func extractToolRoutingDecision(response string) toolRoutingDecision {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || end <= start {
		return toolRoutingDecision{}
	}
	var parsed toolRoutingDecision
	if err := json.Unmarshal([]byte(response[start:end+1]), &parsed); err != nil {
		return toolRoutingDecision{}
	}
	return parsed
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
