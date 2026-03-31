package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Config 全局配置
type Config struct {
	Web      WebConfig      `yaml:"web"`
	Database DatabaseConfig `yaml:"database"`
	User     UserConfig     `yaml:"user"`
	File     FileConfig     `yaml:"file"`
}

// WebConfig Web服务配置
type WebConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	TLS             bool   `yaml:"tls"`
	CertPrivatePath string `yaml:"certPrivatePath"`
	CertPublicPath  string `yaml:"certPublicPath"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type string `yaml:"type"` // sqlite, mysql, postgres
	Path string `yaml:"path"` // SQLite路径
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

// UserConfig 用户相关配置
type UserConfig struct {
	CookieTimeout int    `yaml:"cookieTimeout"`
	PassHashType  string `yaml:"passHashType"` // text, md5, md5salt
}

// FileConfig 文件上传配置
type FileConfig struct {
	MaxSize        uint64            `yaml:"maxSize"`
	Paths          map[string]string `yaml:"paths"`
	AllowImageMime map[string]string `yaml:"allowImageMime"`
	AllowVideoMime map[string]string `yaml:"allowVideoMime"`
	AllowMusicMime map[string]string `yaml:"allowMusicMime"`
	AllowPdfMime   map[string]string `yaml:"allowPdfMime"`
}

// Current 全局配置实例
var Current *Config

// Load 加载配置文件
func Load(configPath string) error {
	// 如果配置文件不存在，创建默认配置
	if !fileExists(configPath) {
		if err := createDefaultConfig(configPath); err != nil {
			return err
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 解析YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}

	Current = config
	return nil
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建默认配置文件
func createDefaultConfig(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 默认配置
	defaultConfig := &Config{
		Web: WebConfig{
			Host:   "127.0.0.1",
			Port:   "8080",
			TLS:    false,
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: "data/database.db",
		},
		User: UserConfig{
			CookieTimeout: 604800,
			PassHashType:  "md5",
		},
		File: FileConfig{
			MaxSize: 52428800, // 50MB
			Paths: map[string]string{
				"avatar": "data/static/avatar/",
				"image":  "data/upload/image/",
				"video":  "data/upload/video/",
				"music":  "data/upload/music/",
				"pdf":    "data/upload/pdf/",
				"other":  "data/upload/other/",
			},
			AllowImageMime: map[string]string{
				"image/jpeg": ".jpeg",
				"image/png":  ".png",
				"image/gif":  ".gif",
				"image/bmp":  ".bmp",
			},
			AllowVideoMime: map[string]string{
				"video/mp4":         ".mp4",
				"video/x-msvideo":   ".avi",
				"video/quicktime":   ".mov",
				"video/x-flv":       ".flv",
				"video/mpeg":        ".mpeg",
			},
			AllowMusicMime: map[string]string{
				"audio/mpeg": ".mpeg",
				"audio/aac":  ".aac",
				"audio/wav":  ".wav",
				"audio/flac": ".flac",
			},
			AllowPdfMime: map[string]string{
				"application/pdf": ".pdf",
			},
		},
	}

	// 序列化为YAML
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(path, data, 0644)
}