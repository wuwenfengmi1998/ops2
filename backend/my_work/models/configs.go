package models

import "github.com/mitchellh/mapstructure"

var Configs map[string]interface{}

type ConfigsWeb_ struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Tls             bool   `mapstructure:"tls"`
	CertPrivatePath string `mapstructure:"certPrivatePath"`
	CertPublicPath  string `mapstructure:"certPublicPath"`
}

type ConfigsUser_ struct {
	CookieTimeout int    `mapstructure:"cookieTimeout"`
	PassHashType  string `mapstructure:"passHashType"`
}

type ConfigsFile_ struct {
	MaxSize        uint64            `mapstructure:"maxSize"`
	Pahts          map[string]string `mapstructure:"pahts"`
	AllowImageMime map[string]string `mapstructure:"allowImageMime"`
	AllowVideoMime map[string]string `mapstructure:"allowVideoMime"`
	AllowMusicMime map[string]string `mapstructure:"allowMusicMime"`
	AllowPdfMime   map[string]string `mapstructure:"allowPdfMime"`
}

type ConfigsAIChatOpenAI_ struct {
	Name                string `mapstructure:"name"`
	Active              bool   `mapstructure:"active"`
	ApiKey              string `mapstructure:"apiKey"`
	BaseUrl             string `mapstructure:"baseUrl"`
	Model               string `mapstructure:"model"`
	Timeout             int    `mapstructure:"timeout"`
	MaxTokens           int    `mapstructure:"maxTokens"`
	ContextWindowTokens int    `mapstructure:"contextWindowTokens"`
	SystemPrompt        string `mapstructure:"systemPrompt"`
}

type ConfigsAIChatTool_ struct {
	Name        string `mapstructure:"name"`
	Enabled     bool   `mapstructure:"enabled"`
	Description string `mapstructure:"description"`
}

type ConfigsAIChatToolRouter_ struct {
	Enabled    bool                 `mapstructure:"enabled"`
	OpenAIName string               `mapstructure:"openaiName"`
	Timeout    int                  `mapstructure:"timeout"`
	MaxTokens  int                  `mapstructure:"maxTokens"`
	Tools      []ConfigsAIChatTool_ `mapstructure:"tools"`
}

type ConfigsAIChat_ struct {
	Enabled    bool                     `mapstructure:"enabled"`
	OpenAI     []ConfigsAIChatOpenAI_   `mapstructure:"openai"`
	ToolRouter ConfigsAIChatToolRouter_ `mapstructure:"toolRouter"`
}

var ConfigsWed ConfigsWeb_
var ConfigsUser ConfigsUser_
var ConfigsFile ConfigsFile_
var ConfigsAIChat ConfigsAIChat_

func ConfigAllInit() error {

	//初始化数据库
	DatabaseInit()

	//读取web配置
	err := mapstructure.Decode(Configs["web"].(map[string]interface{}), &ConfigsWed)
	if err != nil {
		panic(err)
	}

	//初始化user config
	err = mapstructure.Decode(Configs["user"].(map[string]interface{}), &ConfigsUser)
	if err != nil {
		panic(err)
	}

	//初始化file config
	err = mapstructure.Decode(Configs["file"].(map[string]interface{}), &ConfigsFile)
	if err != nil {
		panic(err)
	}

	//初始化aichat config
	ConfigsAIChat = ConfigsAIChat_{
		Enabled: false,
		OpenAI: []ConfigsAIChatOpenAI_{
			{
				Name:         "default",
				Active:       true,
				BaseUrl:      "https://ark.cn-beijing.volces.com/api/v3",
				Timeout:      120,
				MaxTokens:    4096,
				SystemPrompt: "你是一个有帮助的 AI 助手。",
			},
		},
		ToolRouter: ConfigsAIChatToolRouter_{
			Enabled:    true,
			OpenAIName: "default",
			Timeout:    30,
			MaxTokens:  512,
			Tools: []ConfigsAIChatTool_{
				{Name: "time", Enabled: true, Description: "解析当前时间、相对日期和日期范围。"},
				{Name: "ops_ai_assistant_current_user", Enabled: true, Description: "返回当前登录用户信息；未登录时提示需要登录才能获取信息。"},
				{Name: "ops_ai_assistant_schedule_query", Enabled: true, Description: "按日期范围查询当前用户可见的 OPS 日历/日程。"},
			},
		},
	}
	if aiChatConfig, ok := Configs["aichat"].(map[string]interface{}); ok {
		err = mapstructure.Decode(aiChatConfig, &ConfigsAIChat)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
