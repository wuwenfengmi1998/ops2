package handler

import (
	"ops/internal/service"
	"ops/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PurchaseHandler struct {
	service service.PurchaseService
}

func NewPurchaseHandler(db *gorm.DB) *PurchaseHandler {
	return &PurchaseHandler{
		service: service.NewPurchaseService(db),
	}
}

// GetOrders 获取采购订单列表
// @Summary 获取采购订单列表
// @Description 获取用户采购订单列表，支持搜索和分页
// @Tags 采购管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param search query string false "搜索关键词"
// @Param page query int true "页码" default(1)
// @Param entries query int true "每页数量" default(20)
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 400 {object} response.StandardResponse "参数错误"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 500 {object} response.StandardResponse "服务器错误"
// @Router /api/v1/purchase/orders [get]
func (h *PurchaseHandler) GetOrders(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取查询参数
	search := c.Query("search")
	page := GetIntParam(c, "page", 1)
	entries := GetIntParam(c, "entries", 20)

	// 调用Service
	result, success := h.service.GetOrders(c, userID.(uint), search, page, entries)
	if !success {
		response.BadRequest(c, "参数错误")
		return
	}

	response.Success(c, result)
}

// CreateOrder 创建采购订单
// @Summary 创建采购订单
// @Description 创建新的采购订单
// @Tags 采购管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param request body service.CreateOrderRequest true "订单信息"
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 400 {object} response.StandardResponse "参数错误"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 500 {object} response.StandardResponse "服务器错误"
// @Router /api/v1/purchase/orders [post]
func (h *PurchaseHandler) CreateOrder(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 解析请求体
	var request service.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field()+" "+err.Tag())
		}
		if len(validationErrors) > 0 {
			response.BadRequest(c, "参数错误: "+validationErrors[0])
		} else {
			response.BadRequest(c, "请求格式错误")
		}
		return
	}

	// 调用Service
	success := h.service.CreateOrder(c, userID.(uint), request)
	if !success {
		response.BadRequest(c, "创建订单失败，请检查数据")
		return
	}

	response.Success(c, gin.H{"message": "订单创建成功"})
}

// GetOrderDetails 获取订单详情
// @Summary 获取订单详情
// @Description 获取采购订单的详细信息
// @Tags 采购管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param id path int true "订单ID"
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 404 {object} response.StandardResponse "订单不存在"
// @Failure 500 {object} response.StandardResponse "服务器错误"
// @Router /api/v1/purchase/orders/{id} [get]
func (h *PurchaseHandler) GetOrderDetails(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取订单ID
	orderID := GetUintParam(c, "id")
	if orderID == 0 {
		response.BadRequest(c, "订单ID无效")
		return
	}

	// 调用Service
	order, costs, err := h.service.GetOrderDetails(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "-5", "订单不存在")
		} else {
			response.InternalError(c, err)
		}
		return
	}

	// 检查订单所属用户
	if order.UserID != userID.(uint) {
		response.Unauthorized(c)
		return
	}

	response.Success(c, gin.H{
		"order": order,
		"costs": costs,
	})
}