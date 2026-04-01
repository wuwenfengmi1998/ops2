package database

// AutoMigrate 自动迁移所有表
func AutoMigrate() error {
	models := []interface{}{
		&TabUser{},
		&TabUserGroups{},
		&TabUserGroupBinds{},
		&TabUserInfo{},
		&TabCookie{},
		&TabFileInfo{},
		&APIRequestLog{},
		&TabPurchaseOrder{},
		&TabPurchaseCosts{},
	}

	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}

// TabUser 用户表
type TabUser struct {
	ID   uint   `gorm:"primarykey;autoIncrement"`
	Name string `gorm:"type:varchar(64);uniqueIndex"`
}

// TabUserGroups 用户组表
type TabUserGroups struct {
	ID   uint   `gorm:"primarykey;autoIncrement"`
	Name string `gorm:"type:varchar(64);uniqueIndex"`
}

// TabUserGroupBinds 用户-组绑定关系表
type TabUserGroupBinds struct {
	UserID  uint `gorm:"index"`
	GroupID uint `gorm:"index"`
}

// TabUserInfo 用户详情表
type TabUserInfo struct {
	UserID       uint   `gorm:"primaryKey"`
	AvatarPath   string `gorm:"type:text"`
	Birthdate    string `gorm:"type:varchar(16)"`
	Gender       int
	Introduction string `gorm:"type:text"`
}

// TabCookie Session Cookie表
type TabCookie struct {
	Value     string `gorm:"primaryKey;type:varchar(64)"`
	UserID    uint   `gorm:"index"`
	ExpiresAt int64
	CreateAt  int64
	Remember  bool
}

// TabFileInfo 文件信息表
type TabFileInfo struct {
	ID         uint   `gorm:"primarykey;autoIncrement"`
	Path       string `gorm:"type:text"`
	Hash       string `gorm:"index"`
	Size       int64
	CreateTime int64
	ExtName    string `gorm:"type:varchar(16)"`
	MimeType   string `gorm:"type:varchar(128)"`
	StoreType  int    // 1=image 2=video 3=music 4=pdf 5=other
}

// APIRequestLog API请求日志表
type APIRequestLog struct {
	ID       uint   `gorm:"primarykey;autoIncrement"`
	Time     int64  `gorm:"index"`
	IP       string `gorm:"type:varchar(64)"`
	Path     string `gorm:"type:varchar(255)"`
	Method   string `gorm:"type:varchar(16)"`
	Status   int
	UserID   uint
	UserType int
	DataSize int
}

// TabPurchaseOrder 采购订单表
type TabPurchaseOrder struct {
	ID           uint   `gorm:"primarykey;autoIncrement"`
	Title        string `gorm:"type:varchar(255)"`
	CreateTime   int64  `gorm:"index"`
	CompleteTime int64
	Status       int // 状态：0=进行中 1=已完成 2=已取消
	CourierNum   string `gorm:"type:text"`         // 快递单号
	Photos       string `gorm:"type:text"`         // 照片JSON数组
	Creater      uint   `gorm:"index"`             // 创建者ID
	Remark       string `gorm:"type:text"`         // 备注
}

// TabPurchaseCosts 采购费用明细表
type TabPurchaseCosts struct {
	ID           uint   `gorm:"primarykey;autoIncrement"`
	OrderID      uint   `gorm:"index"`
	Name         string `gorm:"type:varchar(255)"`
	PricePerUnit string `gorm:"type:varchar(32)"`
	Quantity     string `gorm:"type:varchar(32)"`
	Unit         string `gorm:"type:varchar(32)"`
}