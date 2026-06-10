package routers

import (
	"encoding/json"
	"errors"
	"ops/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type aiChatConversationListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type aiChatConversationIDReq struct {
	ID uint `json:"id"`
}

type aiChatConversationUpdateReq struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func handleAIChatConversationList(ctx *gin.Context) {
	isAuth, user, data := AuthenticationAuthority(ctx)
	if !isAuth {
		ReturnJson(ctx, "userCookieError", nil)
		return
	}

	var req aiChatConversationListReq
	if data != nil {
		_ = decodeJSON(data, &req)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 50
	}

	query := models.DB.Model(&TabAIChatConversation{}).Where("user_id = ?", user.ID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}

	var conversations []TabAIChatConversation
	if err := query.Order("last_message_at desc, updated_at desc, id desc").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&conversations).Error; err != nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}

	items := make([]map[string]interface{}, 0, len(conversations))
	for _, conversation := range conversations {
		items = append(items, aiChatConversationResponse(conversation))
	}
	ReturnJson(ctx, "apiOK", gin.H{
		"conversations": items,
		"total":         total,
		"page":          req.Page,
		"page_size":     req.PageSize,
	})
}

func handleAIChatConversationGet(ctx *gin.Context) {
	isAuth, user, data := AuthenticationAuthority(ctx)
	if !isAuth {
		ReturnJson(ctx, "userCookieError", nil)
		return
	}

	var req aiChatConversationIDReq
	if err := decodeJSON(data, &req); err != nil || req.ID == 0 {
		ReturnJson(ctx, "jsonErr", nil)
		return
	}

	var conversation TabAIChatConversation
	if err := models.DB.Where("id = ? AND user_id = ?", req.ID, user.ID).First(&conversation).Error; err != nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}

	var messages []TabAIChatMessage
	if err := models.DB.Where("conversation_id = ? AND user_id = ?", conversation.ID, user.ID).
		Order("seq asc, id asc").
		Find(&messages).Error; err != nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}

	items := make([]map[string]interface{}, 0, len(messages))
	for _, message := range messages {
		items = append(items, aiChatMessageResponse(message))
	}
	ReturnJson(ctx, "apiOK", gin.H{
		"conversation": aiChatConversationResponse(conversation),
		"messages":     items,
	})
}

func handleAIChatConversationUpdate(ctx *gin.Context) {
	isAuth, user, data := AuthenticationAuthority(ctx)
	if !isAuth {
		ReturnJson(ctx, "userCookieError", nil)
		return
	}

	var req aiChatConversationUpdateReq
	if err := decodeJSON(data, &req); err != nil || req.ID == 0 {
		ReturnJson(ctx, "jsonErr", nil)
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "新对话"
	}
	if len([]rune(title)) > 200 {
		title = string([]rune(title)[:200])
	}

	result := models.DB.Model(&TabAIChatConversation{}).
		Where("id = ? AND user_id = ?", req.ID, user.ID).
		Updates(map[string]interface{}{"title": title})
	if result.Error != nil || result.RowsAffected == 0 {
		ReturnJson(ctx, "apiErr", nil)
		return
	}
	ReturnJson(ctx, "apiOK", gin.H{"id": req.ID, "title": title})
}

func handleAIChatConversationDelete(ctx *gin.Context) {
	isAuth, user, data := AuthenticationAuthority(ctx)
	if !isAuth {
		ReturnJson(ctx, "userCookieError", nil)
		return
	}

	var req aiChatConversationIDReq
	if err := decodeJSON(data, &req); err != nil || req.ID == 0 {
		ReturnJson(ctx, "jsonErr", nil)
		return
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND user_id = ?", req.ID, user.ID).Delete(&TabAIChatConversation{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Where("conversation_id = ? AND user_id = ?", req.ID, user.ID).Delete(&TabAIChatMessage{}).Error
	})
	if err != nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}
	ReturnJson(ctx, "apiOK", gin.H{"id": req.ID})
}

func aiChatConversationResponse(conversation TabAIChatConversation) map[string]interface{} {
	return map[string]interface{}{
		"id":            conversation.ID,
		"title":         conversation.Title,
		"openaiName":    conversation.OpenAIName,
		"clientLocalId": conversation.ClientLocalID,
		"messageCount":  conversation.MessageCount,
		"lastMessageAt": conversation.LastMessageAt,
		"createdAt":     conversation.CreatedAt,
		"updatedAt":     conversation.UpdatedAt,
	}
}

func aiChatMessageResponse(message TabAIChatMessage) map[string]interface{} {
	item := map[string]interface{}{
		"id":             message.ID,
		"conversationId": message.ConversationID,
		"role":           message.Role,
		"content":        message.Content,
		"image_url":      message.ImageURL,
		"reasoning":      message.Reasoning,
		"seq":            message.Seq,
		"createdAt":      message.CreatedAt,
	}
	var traces []sseEvent
	if message.TracesJSON != "" && json.Unmarshal([]byte(message.TracesJSON), &traces) == nil {
		item["traces"] = traces
	}
	var stats tokenUsageStats
	if message.StatsJSON != "" && json.Unmarshal([]byte(message.StatsJSON), &stats) == nil {
		item["stats"] = stats
	}
	return item
}

func prepareAIChatPersistence(user *TabUser, req chatRequest, openAIName string) (*TabAIChatConversation, error) {
	if user == nil || !req.SaveToServer {
		return nil, nil
	}

	lastUserMessage, ok := lastAIChatUserMessage(req.Messages)
	if !ok {
		return nil, nil
	}

	var conversation TabAIChatConversation
	if req.ConversationID > 0 {
		if err := models.DB.Where("id = ? AND user_id = ?", req.ConversationID, user.ID).First(&conversation).Error; err != nil {
			return nil, err
		}
	} else if strings.TrimSpace(req.ClientLocalID) != "" {
		err := models.DB.Where("user_id = ? AND client_local_id = ?", user.ID, strings.TrimSpace(req.ClientLocalID)).First(&conversation).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			conversation = TabAIChatConversation{}
		}
	}

	now := time.Now()
	if conversation.ID == 0 {
		conversation = TabAIChatConversation{
			UserID:        user.ID,
			Title:         makeAIChatTitle(lastUserMessage),
			OpenAIName:    openAIName,
			ClientLocalID: strings.TrimSpace(req.ClientLocalID),
			LastMessageAt: &now,
		}
		if err := models.DB.Create(&conversation).Error; err != nil {
			return nil, err
		}
	} else if conversation.OpenAIName == "" && openAIName != "" {
		_ = models.DB.Model(&conversation).Update("open_ai_name", openAIName).Error
		conversation.OpenAIName = openAIName
	}

	if err := createAIChatMessage(conversation.ID, user.ID, lastUserMessage, "", nil, nil); err != nil {
		return nil, err
	}
	if err := refreshAIChatConversationSummary(conversation.ID, user.ID); err != nil {
		return nil, err
	}
	if err := models.DB.Where("id = ? AND user_id = ?", conversation.ID, user.ID).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func persistAIChatAssistantMessage(conversation *TabAIChatConversation, content string, reasoning string, traces []sseEvent, stats *tokenUsageStats) error {
	if conversation == nil {
		return nil
	}
	message := chatMessage{Role: "assistant", Content: content}
	if err := createAIChatMessage(conversation.ID, conversation.UserID, message, reasoning, traces, stats); err != nil {
		return err
	}
	return refreshAIChatConversationSummary(conversation.ID, conversation.UserID)
}

func createAIChatMessage(conversationID uint, userID uint, message chatMessage, reasoning string, traces []sseEvent, stats *tokenUsageStats) error {
	var maxSeq int
	models.DB.Model(&TabAIChatMessage{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Select("COALESCE(MAX(seq), 0)").
		Scan(&maxSeq)

	imageURL := message.ImageURL
	if imageURL == "" {
		imageURL = message.ImageURLAlias
	}
	tracesJSON := ""
	if len(traces) > 0 {
		if data, err := json.Marshal(traces); err == nil {
			tracesJSON = string(data)
		}
	}
	statsJSON := ""
	if stats != nil {
		if data, err := json.Marshal(stats); err == nil {
			statsJSON = string(data)
		}
	}

	row := TabAIChatMessage{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           message.Role,
		Content:        message.Content,
		ImageURL:       imageURL,
		Reasoning:      reasoning,
		TracesJSON:     tracesJSON,
		StatsJSON:      statsJSON,
		Seq:            maxSeq + 1,
	}
	return models.DB.Create(&row).Error
}

func refreshAIChatConversationSummary(conversationID uint, userID uint) error {
	var count int64
	if err := models.DB.Model(&TabAIChatMessage{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Count(&count).Error; err != nil {
		return err
	}
	var lastMessage TabAIChatMessage
	lastMessageAt := time.Now()
	if err := models.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).Order("seq desc, id desc").First(&lastMessage).Error; err == nil && lastMessage.CreatedAt != nil {
		lastMessageAt = *lastMessage.CreatedAt
	}
	return models.DB.Model(&TabAIChatConversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{
			"message_count":   int(count),
			"last_message_at": lastMessageAt,
		}).Error
}

func lastAIChatUserMessage(messages []chatMessage) (chatMessage, bool) {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i], true
		}
	}
	return chatMessage{}, false
}

func makeAIChatTitle(message chatMessage) string {
	title := strings.TrimSpace(message.Content)
	if title == "" {
		title = "新对话"
	}
	runes := []rune(title)
	if len(runes) > 40 {
		title = string(runes[:40])
	}
	return title
}
