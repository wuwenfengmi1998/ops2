package routers

import (
	"ops/models"
	"time"
)

//跨模块绑定区
//绑定数据统一处理区

type TabPurchaseFileBind struct {
	ID        uint       `gorm:"primarykey"`
	OrderID   uint       `gorm:"not null"`
	FileID    uint       `gorm:"not null"`
	CreatedAt *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseItemFileBind struct {
	ID        uint      `gorm:"primaryKey"`
	ItemID    uint      `gorm:"not null;index;comment:关联物品id"`
	FileID    uint      `gorm:"not null;comment:关联文件id"`
	CreatorID uint      `gorm:"not null;comment:上传人id"`
	CreatedAt time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseItemWorkOrderBind struct {
	ID          uint      `gorm:"primaryKey"`
	ItemID      uint      `gorm:"not null;index;comment:关联物品id"`
	WorkOrderID uint      `gorm:"not null;index;comment:关联工单id"`
	Remark      string    `gorm:"size:500;comment:备注"`
	CreatorID   uint      `gorm:"not null;comment:绑定人id"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWorkOrderFileBind struct {
	ID          uint       `gorm:"primarykey"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	FileID      uint       `gorm:"not null;comment:关联文件ID"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWorkOrderCommitFileBind struct {
	ID          uint       `gorm:"primarykey"`
	CommitID    uint       `gorm:"not null;index;comment:关联进度ID"`
	FileID      uint       `gorm:"not null;comment:关联文件ID"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabWorkOrderPurchaseOrderBind 工单与采购订单的关联表
type TabWorkOrderPurchaseOrderBind struct {
	ID              uint       `gorm:"primarykey"`
	WorkOrderID     uint       `gorm:"not null;index;comment:关联工单ID"`
	CommitID        uint       `gorm:"not null;index;comment:关联进度ID"`
	PurchaseOrderID uint       `gorm:"not null;comment:关联采购订单ID"`
	CreatedAt       *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabWorkOrderCustomerBind 工单与客户关联表
type TabWorkOrderCustomerBind struct {
	ID          uint       `gorm:"primarykey"`
	WorkOrderID uint       `gorm:"not null;index;comment:关联工单ID"`
	CustomerID  uint       `gorm:"not null;index;comment:关联客户ID"`
	CreatorID   uint       `gorm:"not null;comment:绑定人id"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime"`
}

type TabWarehouseContainerFileBind struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerID uint      `gorm:"not null;index;comment:关联容器id"`
	FileID      uint      `gorm:"not null;comment:关联文件id"`
	CreatorID   uint      `gorm:"not null;comment:上传人id"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabWarehouseItemCustomerBind 物品与客户关联表
type TabWarehouseItemCustomerBind struct {
	ID         uint       `gorm:"primarykey"`
	ItemID     uint       `gorm:"not null;index;comment:关联物品ID"`
	CustomerID uint       `gorm:"not null;index;comment:关联客户ID"`
	CreatorID  uint       `gorm:"not null;comment:绑定人id"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime"`
}

func BindsInit() {
	models.DB.AutoMigrate(
		&TabPurchaseFileBind{},
		&TabWarehouseItemFileBind{},
		&TabWarehouseItemWorkOrderBind{},
		&TabWarehouseContainerFileBind{},
		&TabWorkOrderFileBind{},
		&TabWorkOrderCommitFileBind{},
		&TabWorkOrderPurchaseOrderBind{},
		&TabWorkOrderCustomerBind{},
		&TabWarehouseItemCustomerBind{},
	)

}
