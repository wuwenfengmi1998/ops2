package models

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// type APIRequestLog_ struct {
// 	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
// 	IPAddress  string    `gorm:"column:ip_address;size:45;not null" json:"ip_address"`
// 	Path       string    `gorm:"column:path;size:500;not null" json:"path"`
// 	Method     string    `gorm:"column:method;size:10;not null" json:"method"`
// 	StatusCode int       `gorm:"column:status_code;index" json:"status_code"`
// 	Message    string    `gorm:"column:error_message;type:text" json:"error_message"`
// 	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
// }

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

	//DB.AutoMigrate(&APIRequestLog_{})

	return nil
}
