package models

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type TabFileInfo_ struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	Name   string    `gorm:"not null;size:256;index"` // 前端报告的文件名
	Path   string    `gorm:"not null;size:300"`       //
	Sha256 string    `gorm:"not null;size:256;index"` //
	Mime   string    `gorm:"size:64;index"`
	Type   string    `gorm:"size:64;index"`
	Const  uint      `gorm:"default:1;index"`
	Per    uint      `gorm:"default:1"`
	UserID uint      `gorm:"not null;index"`
	Date   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUser_ struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"` // 自增主键
	Name  string    `gorm:"size:100;uniqueIndex"`     // 唯一约束索引
	Email string    `gorm:"size:255;index"`           // 字符串长度限制100 索引
	Pass  string    `gorm:"size:128"`                 // 建议存储哈希后的密码
	Type  string    `gorm:"size:64;default:user"`     //
	Salt  string    `gorm:"size:64;"`
	Date  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserGroups_ struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`                // 自增主键
	Name  string    `gorm:"size:100;uniqueIndex"`                    // 唯一约束索引
	Email string    `gorm:"size:255;index"`                          // 字符串长度限制100 索引
	Type  string    `gorm:"size:64;default:usergroup"`               //
	Date  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserGroupBinds_ struct {
	ID      uint      `gorm:"primaryKey;autoIncrement"` // 自增主键
	UserID  uint      `gorm:"index"`
	GroupID uint      `gorm:"index"`
	Date    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 默认当前时间
}

type TabUserInfo_ struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null;uniqueIndex"`
	FirstName  string    `gorm:"size:50;null"`
	Username   string    `gorm:"size:30;null"`
	Birthdate  time.Time `gorm:"type:datetime;null"`
	Gender     string    `gorm:"type:char(1);check:gender IN ('M', 'F', 'U');default:'U'"`
	AvatarPath string    `gorm:"size:255"`
	Region     string    `gorm:"size:50"`
	Language   string    `gorm:"size:10;default:'zh-CN'"`
	CreatedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:created_at"`
}

// var def_user_info = User_info{
// 	ID:0,
// 	UserID:0,
// }

type TabCookie_ struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"size:16;not null"`
	Name      string    `gorm:"size:255;not null;index"`
	Value     string    `gorm:"size:255;not null;index"`
	ExpiresAt time.Time `gorm:"type:datetime;index"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;index;not null;default:CURRENT_TIMESTAMP"`
	Remember  bool      `gorm:"default:false"`
}

type APIRequestLog_ struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	IPAddress  string    `gorm:"column:ip_address;size:45;not null" json:"ip_address"`
	Path       string    `gorm:"column:path;size:500;not null" json:"path"`
	Method     string    `gorm:"column:method;size:10;not null" json:"method"`
	StatusCode int       `gorm:"column:status_code;index" json:"status_code"`
	Message    string    `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func DatabaseInit() error {
	var err error
	fmt.Println("database_init")
	DatabaseConfigs := Configs["database"].(map[string]interface{})

	if DatabaseConfigs["type"].(string) == "sqlite" {
		//sqlite init
		fmt.Println("sqlite")
		DB, err = gorm.Open(sqlite.Open(DatabaseConfigs["path"].(string)), &gorm.Config{})
	} else if DatabaseConfigs["type"].(string) == "mysql" {
		//mysql init
		fmt.Println("mysql")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DatabaseConfigs["user"].(string), DatabaseConfigs["pass"].(string), DatabaseConfigs["host"].(string), DatabaseConfigs["port"].(string), DatabaseConfigs["name"].(string))
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if DatabaseConfigs["type"].(string) == "pg" {
		//postgresql init
		fmt.Println("postgresql")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", DatabaseConfigs["host"].(string), DatabaseConfigs["user"].(string), DatabaseConfigs["pass"].(string), DatabaseConfigs["name"].(string), DatabaseConfigs["port"].(string))
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		fmt.Println(err)
		panic("数据库连接失败")
	}

	// 自动创建表结构
	DB.AutoMigrate(&TabUser_{})

	DB.AutoMigrate(&TabUserGroups_{})

	DB.AutoMigrate(&TabUserGroupBinds_{})

	DB.AutoMigrate(&TabUserInfo_{})

	DB.AutoMigrate(&TabCookie_{})

	DB.AutoMigrate(&TabFileInfo_{})

	DB.AutoMigrate(&APIRequestLog_{})

	return nil
}
