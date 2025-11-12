package models

import "time"

type APIRequestLog struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	IPAddress  string    `gorm:"column:ip_address;size:45;not null" json:"ip_address"`
	Path       string    `gorm:"column:path;size:500;not null" json:"path"`
	Method     string    `gorm:"column:method;size:10;not null" json:"method"`
	StatusCode int       `gorm:"column:status_code;index" json:"status_code"`
	Message    string    `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func LogAdd() {

}
