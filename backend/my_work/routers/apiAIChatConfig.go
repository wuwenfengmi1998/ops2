package routers

import (
	"errors"
	"ops/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TabAIChatSetting struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	Enabled   bool       `gorm:"default:false;index"`
	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"type:datetime;autoUpdateTime"`
}

type TabAIChatOpenAIProfile struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement"`
	Name                string     `gorm:"size:100;not null;uniqueIndex"`
	Active              bool       `gorm:"default:false;index"`
	ApiKey              string     `gorm:"type:text"`
	BaseUrl             string     `gorm:"size:500"`
	Model               string     `gorm:"size:200"`
	Timeout             int        `gorm:"default:120"`
	MaxTokens           int        `gorm:"default:4096"`
	ContextWindowTokens int        `gorm:"default:0"`
	SystemPrompt        string     `gorm:"type:text"`
	SortOrder           int        `gorm:"default:0;index"`
	CreatedAt           *time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt           *time.Time `gorm:"type:datetime;autoUpdateTime"`
}

type TabAIChatToolRouter struct {
	ID         uint       `gorm:"primaryKey;autoIncrement"`
	Enabled    bool       `gorm:"default:true"`
	OpenAIName string     `gorm:"size:100;index"`
	Timeout    int        `gorm:"default:30"`
	MaxTokens  int        `gorm:"default:512"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt  *time.Time `gorm:"type:datetime;autoUpdateTime"`
}

type TabAIChatTool struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	Name        string     `gorm:"size:100;not null;uniqueIndex"`
	Enabled     bool       `gorm:"default:true;index"`
	Description string     `gorm:"type:text"`
	SortOrder   int        `gorm:"default:0;index"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"type:datetime;autoUpdateTime"`
}

type TabAIChatConversation struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	UserID        uint           `gorm:"not null;index:idx_ai_chat_conversation_user_updated,priority:1"`
	Title         string         `gorm:"size:200"`
	OpenAIName    string         `gorm:"size:100;index"`
	ClientLocalID string         `gorm:"size:100;index"`
	MessageCount  int            `gorm:"default:0"`
	LastMessageAt *time.Time     `gorm:"type:datetime;index"`
	CreatedAt     *time.Time     `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"type:datetime;autoUpdateTime;index:idx_ai_chat_conversation_user_updated,priority:2"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type TabAIChatMessage struct {
	ID             uint       `gorm:"primaryKey;autoIncrement"`
	ConversationID uint       `gorm:"not null;index:idx_ai_chat_message_conversation_seq,priority:1"`
	UserID         uint       `gorm:"not null;index"`
	Role           string     `gorm:"size:20;index"`
	Content        string     `gorm:"type:text"`
	ImageURL       string     `gorm:"type:text"`
	Reasoning      string     `gorm:"type:text"`
	TracesJSON     string     `gorm:"type:text"`
	StatsJSON      string     `gorm:"type:text"`
	Seq            int        `gorm:"index:idx_ai_chat_message_conversation_seq,priority:2"`
	CreatedAt      *time.Time `gorm:"type:datetime;autoCreateTime"`
}

var aiChatConfigMu sync.RWMutex

func ApiAIChatInit() {
	err := models.DB.AutoMigrate(
		&TabAIChatSetting{},
		&TabAIChatOpenAIProfile{},
		&TabAIChatToolRouter{},
		&TabAIChatTool{},
		&TabAIChatConversation{},
		&TabAIChatMessage{},
	)
	if err != nil {
		panic(err)
	}

	if err := seedAIChatConfigFromYAMLIfEmpty(); err != nil {
		panic(err)
	}
	if err := RefreshAIChatConfigCache(); err != nil {
		panic(err)
	}
}

func seedAIChatConfigFromYAMLIfEmpty() error {
	var settingCount int64
	var profileCount int64
	models.DB.Model(&TabAIChatSetting{}).Count(&settingCount)
	models.DB.Model(&TabAIChatOpenAIProfile{}).Count(&profileCount)
	if settingCount > 0 || profileCount > 0 {
		return nil
	}

	cfg := models.ConfigsAIChat
	return models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&TabAIChatSetting{ID: 1, Enabled: cfg.Enabled}).Error; err != nil {
			return err
		}

		profiles := cfg.OpenAI
		if len(profiles) == 0 {
			profiles = []models.ConfigsAIChatOpenAI_{{
				Name:                "default",
				Active:              true,
				BaseUrl:             "https://ark.cn-beijing.volces.com/api/v3",
				Timeout:             120,
				MaxTokens:           4096,
				ContextWindowTokens: 0,
				SystemPrompt:        "你是一个有帮助的 AI 助手。",
			}}
		}
		for i, profile := range profiles {
			if profile.Name == "" {
				profile.Name = "default"
			}
			if profile.Timeout <= 0 {
				profile.Timeout = 120
			}
			if profile.MaxTokens <= 0 {
				profile.MaxTokens = 4096
			}
			if err := tx.Create(&TabAIChatOpenAIProfile{
				Name:                profile.Name,
				Active:              profile.Active,
				ApiKey:              profile.ApiKey,
				BaseUrl:             profile.BaseUrl,
				Model:               profile.Model,
				Timeout:             profile.Timeout,
				MaxTokens:           profile.MaxTokens,
				ContextWindowTokens: nonNegativeInt(profile.ContextWindowTokens),
				SystemPrompt:        profile.SystemPrompt,
				SortOrder:           i,
			}).Error; err != nil {
				return err
			}
		}

		toolRouter := cfg.ToolRouter
		if toolRouter.Timeout <= 0 {
			toolRouter.Timeout = 30
		}
		if toolRouter.MaxTokens <= 0 {
			toolRouter.MaxTokens = 512
		}
		if toolRouter.OpenAIName == "" && len(profiles) > 0 {
			toolRouter.OpenAIName = profiles[0].Name
		}
		if err := tx.Create(&TabAIChatToolRouter{
			ID:         1,
			Enabled:    toolRouter.Enabled,
			OpenAIName: toolRouter.OpenAIName,
			Timeout:    toolRouter.Timeout,
			MaxTokens:  toolRouter.MaxTokens,
		}).Error; err != nil {
			return err
		}

		tools := toolRouter.Tools
		if len(tools) == 0 {
			tools = []models.ConfigsAIChatTool_{{Name: "time", Enabled: true, Description: "提供当前日期、时间和相对日期换算。"}}
		}
		for i, tool := range tools {
			if tool.Name == "" {
				continue
			}
			if err := tx.Create(&TabAIChatTool{
				Name:        tool.Name,
				Enabled:     tool.Enabled,
				Description: tool.Description,
				SortOrder:   i,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func RefreshAIChatConfigCache() error {
	var setting TabAIChatSetting
	if err := models.DB.First(&setting, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	var profiles []TabAIChatOpenAIProfile
	if err := models.DB.Order("sort_order asc, id asc").Find(&profiles).Error; err != nil {
		return err
	}

	var toolRouter TabAIChatToolRouter
	if err := models.DB.First(&toolRouter, 1).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var tools []TabAIChatTool
	if err := models.DB.Order("sort_order asc, id asc").Find(&tools).Error; err != nil {
		return err
	}

	cfg := models.ConfigsAIChat_{
		Enabled: setting.Enabled,
		OpenAI:  make([]models.ConfigsAIChatOpenAI_, 0, len(profiles)),
		ToolRouter: models.ConfigsAIChatToolRouter_{
			Enabled:    toolRouter.Enabled,
			OpenAIName: toolRouter.OpenAIName,
			Timeout:    defaultInt(toolRouter.Timeout, 30),
			MaxTokens:  defaultInt(toolRouter.MaxTokens, 512),
			Tools:      make([]models.ConfigsAIChatTool_, 0, len(tools)),
		},
	}

	for _, profile := range profiles {
		cfg.OpenAI = append(cfg.OpenAI, models.ConfigsAIChatOpenAI_{
			Name:                profile.Name,
			Active:              profile.Active,
			ApiKey:              profile.ApiKey,
			BaseUrl:             profile.BaseUrl,
			Model:               profile.Model,
			Timeout:             defaultInt(profile.Timeout, 120),
			MaxTokens:           defaultInt(profile.MaxTokens, 4096),
			ContextWindowTokens: nonNegativeInt(profile.ContextWindowTokens),
			SystemPrompt:        profile.SystemPrompt,
		})
	}
	for _, tool := range tools {
		cfg.ToolRouter.Tools = append(cfg.ToolRouter.Tools, models.ConfigsAIChatTool_{
			Name:        tool.Name,
			Enabled:     tool.Enabled,
			Description: tool.Description,
		})
	}

	aiChatConfigMu.Lock()
	models.ConfigsAIChat = cfg
	aiChatConfigMu.Unlock()
	return nil
}

func getAIChatConfig() models.ConfigsAIChat_ {
	aiChatConfigMu.RLock()
	defer aiChatConfigMu.RUnlock()
	return models.ConfigsAIChat
}

func defaultInt(value int, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func nonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func handleAIChatAdminGetConfig(ctx *gin.Context) {
	if ok, _ := requireSysAdmin(ctx); !ok {
		return
	}
	cfg := getAIChatConfig()
	ReturnJson(ctx, "apiOK", gin.H{
		"enabled": cfg.Enabled,
		"openai":  maskAIChatProfiles(cfg.OpenAI),
		"toolRouter": gin.H{
			"enabled":    cfg.ToolRouter.Enabled,
			"openaiName": cfg.ToolRouter.OpenAIName,
			"timeout":    cfg.ToolRouter.Timeout,
			"maxTokens":  cfg.ToolRouter.MaxTokens,
			"tools":      maskAIChatTools(cfg.ToolRouter.Tools),
		},
	})
}

func handleAIChatAdminRefreshCache(ctx *gin.Context) {
	if ok, _ := requireSysAdmin(ctx); !ok {
		return
	}
	if err := RefreshAIChatConfigCache(); err != nil {
		ReturnJson(ctx, "apiErr", gin.H{"error": err.Error()})
		return
	}
	ReturnJson(ctx, "apiOK", nil)
}

func handleAIChatAdminUpdateConfig(ctx *gin.Context) {
	ok, data := requireSysAdmin(ctx)
	if !ok {
		return
	}
	if data == nil {
		ReturnJson(ctx, "apiErr", nil)
		return
	}

	var req models.ConfigsAIChat_
	if err := decodeJSON(data, &req); err != nil {
		ReturnJson(ctx, "apiErr", gin.H{"error": err.Error()})
		return
	}
	if len(req.OpenAI) == 0 {
		ReturnJson(ctx, "apiErr", gin.H{"error": "openai profiles cannot be empty"})
		return
	}

	if err := saveAIChatConfig(req); err != nil {
		ReturnJson(ctx, "apiErr", gin.H{"error": err.Error()})
		return
	}
	if err := RefreshAIChatConfigCache(); err != nil {
		ReturnJson(ctx, "apiErr", gin.H{"error": err.Error()})
		return
	}
	ReturnJson(ctx, "apiOK", nil)
}

func saveAIChatConfig(req models.ConfigsAIChat_) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&TabAIChatSetting{ID: 1, Enabled: req.Enabled}).Error; err != nil {
			return err
		}

		var existingProfiles []TabAIChatOpenAIProfile
		if err := tx.Find(&existingProfiles).Error; err != nil {
			return err
		}
		existingByName := make(map[string]TabAIChatOpenAIProfile, len(existingProfiles))
		for _, profile := range existingProfiles {
			existingByName[profile.Name] = profile
		}

		incomingNames := make(map[string]bool, len(req.OpenAI))
		activeSet := false
		for i, profile := range req.OpenAI {
			if profile.Name == "" {
				continue
			}
			incomingNames[profile.Name] = true
			if profile.Timeout <= 0 {
				profile.Timeout = 120
			}
			if profile.MaxTokens <= 0 {
				profile.MaxTokens = 4096
			}
			profile.ContextWindowTokens = nonNegativeInt(profile.ContextWindowTokens)
			if profile.Active {
				if activeSet {
					profile.Active = false
				} else {
					activeSet = true
				}
			}

			tab := TabAIChatOpenAIProfile{
				Name:                profile.Name,
				Active:              profile.Active,
				ApiKey:              profile.ApiKey,
				BaseUrl:             profile.BaseUrl,
				Model:               profile.Model,
				Timeout:             profile.Timeout,
				MaxTokens:           profile.MaxTokens,
				ContextWindowTokens: profile.ContextWindowTokens,
				SystemPrompt:        profile.SystemPrompt,
				SortOrder:           i,
			}
			if old, ok := existingByName[profile.Name]; ok {
				tab.ID = old.ID
				if tab.ApiKey == "" {
					tab.ApiKey = old.ApiKey
				}
			}
			if err := tx.Save(&tab).Error; err != nil {
				return err
			}
		}
		for _, old := range existingProfiles {
			if !incomingNames[old.Name] {
				if err := tx.Delete(&old).Error; err != nil {
					return err
				}
			}
		}
		if !activeSet && len(req.OpenAI) > 0 {
			if err := tx.Model(&TabAIChatOpenAIProfile{}).Where("name = ?", req.OpenAI[0].Name).Update("active", true).Error; err != nil {
				return err
			}
		}

		toolRouter := req.ToolRouter
		if toolRouter.Timeout <= 0 {
			toolRouter.Timeout = 30
		}
		if toolRouter.MaxTokens <= 0 {
			toolRouter.MaxTokens = 512
		}
		if err := tx.Save(&TabAIChatToolRouter{ID: 1, Enabled: toolRouter.Enabled, OpenAIName: toolRouter.OpenAIName, Timeout: toolRouter.Timeout, MaxTokens: toolRouter.MaxTokens}).Error; err != nil {
			return err
		}

		var existingTools []TabAIChatTool
		if err := tx.Find(&existingTools).Error; err != nil {
			return err
		}
		existingToolByName := make(map[string]TabAIChatTool, len(existingTools))
		for _, tool := range existingTools {
			existingToolByName[tool.Name] = tool
		}
		incomingToolNames := make(map[string]bool, len(toolRouter.Tools))
		for i, tool := range toolRouter.Tools {
			if tool.Name == "" {
				continue
			}
			incomingToolNames[tool.Name] = true
			tab := TabAIChatTool{Name: tool.Name, Enabled: tool.Enabled, Description: tool.Description, SortOrder: i}
			if old, ok := existingToolByName[tool.Name]; ok {
				tab.ID = old.ID
			}
			if err := tx.Save(&tab).Error; err != nil {
				return err
			}
		}
		for _, old := range existingTools {
			if !incomingToolNames[old.Name] {
				if err := tx.Delete(&old).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func maskAIChatProfiles(profiles []models.ConfigsAIChatOpenAI_) []gin.H {
	items := make([]gin.H, 0, len(profiles))
	for _, profile := range profiles {
		items = append(items, gin.H{
			"name":                profile.Name,
			"active":              profile.Active,
			"apiKeySet":           profile.ApiKey != "",
			"baseUrl":             profile.BaseUrl,
			"model":               profile.Model,
			"timeout":             profile.Timeout,
			"maxTokens":           profile.MaxTokens,
			"contextWindowTokens": profile.ContextWindowTokens,
			"systemPrompt":        profile.SystemPrompt,
		})
	}
	return items
}

func maskAIChatTools(tools []models.ConfigsAIChatTool_) []gin.H {
	items := make([]gin.H, 0, len(tools))
	for _, tool := range tools {
		items = append(items, gin.H{
			"name":        tool.Name,
			"enabled":     tool.Enabled,
			"description": tool.Description,
		})
	}
	return items
}

func requireSysAdmin(ctx *gin.Context) (bool, map[string]interface{}) {
	isAuth, user, data := AuthenticationAuthority(ctx)
	if !isAuth {
		ReturnJson(ctx, "userCookieError", nil)
		return false, nil
	}
	if !SysAdminCheck(user.ID) {
		ReturnJson(ctx, "permission_denied", nil)
		return false, nil
	}
	return true, data
}
