package repository

import (
	"ops/models"

	"gorm.io/gorm"
)

type PurchaseRepository interface {
	GetOrders(userID uint, search string, page, entries int) ([]models.TabPurchaseOrder, int64, error)
	GetOrderByID(orderID uint) (*models.TabPurchaseOrder, error)
	CreateOrder(order *models.TabPurchaseOrder) error
	CreateCost(cost *models.TabPurchaseCosts) error
	GetOrderCosts(orderID uint) ([]models.TabPurchaseCosts, error)
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) GetOrders(userID uint, search string, page, entries int) ([]models.TabPurchaseOrder, int64, error) {
	var orders []models.TabPurchaseOrder
	var total int64

	query := r.db.Model(&models.TabPurchaseOrder{}).Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("title LIKE ? OR part_name LIKE ? OR remark LIKE ? OR tracking_number LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := entries * (page - 1)
	if err := query.Order("created_at DESC").Offset(offset).Limit(entries).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *purchaseRepository) GetOrderByID(orderID uint) (*models.TabPurchaseOrder, error) {
	var order models.TabPurchaseOrder
	if err := r.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *purchaseRepository) CreateOrder(order *models.TabPurchaseOrder) error {
	return r.db.Create(order).Error
}

func (r *purchaseRepository) CreateCost(cost *models.TabPurchaseCosts) error {
	return r.db.Create(cost).Error
}

func (r *purchaseRepository) GetOrderCosts(orderID uint) ([]models.TabPurchaseCosts, error) {
	var costs []models.TabPurchaseCosts
	if err := r.db.Where("order_id = ?", orderID).Find(&costs).Error; err != nil {
		return nil, err
	}
	return costs, nil
}