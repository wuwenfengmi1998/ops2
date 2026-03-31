package service

import (
	"encoding/json"
	"ops/internal/repository"
	"ops/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PurchaseService interface {
	GetOrders(c *gin.Context, userID uint, search string, page, entries int) (gin.H, bool)
	CreateOrder(c *gin.Context, userID uint, request CreateOrderRequest) bool
	GetOrderDetails(orderID uint) (*models.TabPurchaseOrder, []models.TabPurchaseCosts, error)
}

type purchaseService struct {
	repo repository.PurchaseRepository
}

func NewPurchaseService(db *gorm.DB) PurchaseService {
	return &purchaseService{
		repo: repository.NewPurchaseRepository(db),
	}
}

// 请求结构体
type CostItem struct {
	Cost         int    `json:"cost" binding:"required,min=1"`
	CostT        int    `json:"costt" binding:"required,min=0"`
	CurrencyType string `json:"currencytype" binding:"required"`
	Int          int    `json:"int" binding:"required,min=1"`
	Type         string `json:"type" binding:"required"`
}

type CreateOrderRequest struct {
	Costs          []CostItem `json:"costs" binding:"required,min=1,dive"`
	Link           string     `json:"link"`
	OrderStatus    string     `json:"order_status" binding:"required"`
	PartName       string     `json:"partname"`
	Photos         []string   `json:"photos"`
	Remark         string     `json:"remark"`
	Styles         string     `json:"styles"`
	Title          string     `json:"title" binding:"required"`
	TrackingNumber string     `json:"tracking_number"`
	UpdateTime     string     `json:"update_time"`
}

func (s *purchaseService) GetOrders(c *gin.Context, userID uint, search string, page, entries int) (gin.H, bool) {
	// 验证分页参数
	if entries <= 0 || entries > 300 {
		return nil, false
	}
	if page <= 0 {
		return nil, false
	}

	orders, total, err := s.repo.GetOrders(userID, search, page, entries)
	if err != nil {
		return nil, false
	}

	// 构建响应
	result := gin.H{
		"all_count":  total,
		"all_orders": orders,
	}

	return result, true
}

func (s *purchaseService) CreateOrder(c *gin.Context, userID uint, request CreateOrderRequest) bool {
	// 验证数据
	if request.Title == "" {
		return false
	}

	// 验证价格和数量
	for _, cost := range request.Costs {
		if cost.Cost <= 0 {
			return false
		}
		if cost.Int <= 0 {
			return false
		}
	}

	// 验证图片哈希（简单检查是否包含特殊字符）
	for _, photo := range request.Photos {
		if models.IsContainsSpecialChar(photo) {
			return false
		}
	}

	// 解析更新时间
	var updateTime *time.Time
	if request.UpdateTime != "" {
		parsedTime, err := models.StringToTimePtr(request.UpdateTime)
		if err != nil {
			return false
		}
		updateTime = parsedTime
	}

	// 转换照片数组为JSON
	var photosJSON datatypes.JSON
	if len(request.Photos) > 0 {
		photosBytes, err := json.Marshal(request.Photos)
		if err != nil {
			return false
		}
		photosJSON = datatypes.JSON(photosBytes)
	}

	// 创建订单
	order := &models.TabPurchaseOrder{
		UserID:         userID,
		Title:          request.Title,
		Remark:         request.Remark,
		Photos:         photosJSON,
		Link:           request.Link,
		PartName:       request.PartName,
		Styles:         request.Styles,
		UpdateTime:     updateTime,
		TrackingNumber: request.TrackingNumber,
		OrderStatus:    request.OrderStatus,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return false
	}

	// 创建费用明细
	for _, costItem := range request.Costs {
		cost := &models.TabPurchaseCosts{
			UserID:   userID,
			OrderID:  order.ID,
			Price:    costItem.Cost,
			Quantity: costItem.Int,
		}

		if err := s.repo.CreateCost(cost); err != nil {
			return false
		}
	}

	return true
}

func (s *purchaseService) GetOrderDetails(orderID uint) (*models.TabPurchaseOrder, []models.TabPurchaseCosts, error) {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, nil, err
	}

	costs, err := s.repo.GetOrderCosts(orderID)
	if err != nil {
		return nil, nil, err
	}

	return order, costs, nil
}