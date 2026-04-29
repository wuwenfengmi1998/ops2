package routers

import (
	"slices"
	"time"

	"ops/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	customerUserGroup TabUserGroups
	customerAdmins    []uint
)

// TabCustomer 客户主表
type TabCustomer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	FirstName string         `gorm:"size:100;comment:名" json:"first_name"`
	LastName  string         `gorm:"size:100;comment:姓" json:"last_name"`
	Title     string         `gorm:"size:20;comment:称呼:Unit/Mr/Ms" json:"title"`
	CreatedBy uint           `gorm:"not null;comment:创建人ID" json:"created_by"`
	CreatedAt *time.Time     `gorm:"type:datetime;autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"type:datetime;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TabCustomerPhone 客户电话绑定表
type TabCustomerPhone struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CustomerID uint       `gorm:"not null;index;comment:关联客户ID" json:"customer_id"`
	Prefix     string     `gorm:"size:10;comment:地区前缀:86/852/853" json:"prefix"`
	Phone      string     `gorm:"size:50;comment:电话号码" json:"phone"`
	Label      string     `gorm:"size:20;comment:标签:mobile/work/home/other" json:"label"`
	IsPrimary  bool       `gorm:"default:false;comment:是否主号码" json:"is_primary"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime" json:"created_at"`
}

// TabCustomerEmail 客户邮箱绑定表
type TabCustomerEmail struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CustomerID uint       `gorm:"not null;index;comment:关联客户ID" json:"customer_id"`
	Email      string     `gorm:"size:200;comment:邮箱地址" json:"email"`
	Label      string     `gorm:"size:20;comment:标签:work/personal/other" json:"label"`
	IsPrimary  bool       `gorm:"default:false;comment:是否主邮箱" json:"is_primary"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime" json:"created_at"`
}

// TabCustomerCompany 客户单位绑定表
type TabCustomerCompany struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CustomerID  uint       `gorm:"not null;index;comment:关联客户ID" json:"customer_id"`
	CompanyName string     `gorm:"size:200;comment:单位名称" json:"company_name"`
	Department  string     `gorm:"size:100;comment:部门" json:"department"`
	Position    string     `gorm:"size:100;comment:职位" json:"position"`
	IsPrimary   bool       `gorm:"default:false;comment:是否主单位" json:"is_primary"`
	CreatedAt   *time.Time `gorm:"type:datetime;autoCreateTime" json:"created_at"`
}

// TabCustomerLog 客户操作日志
type TabCustomerLog struct {
	ID         uint       `gorm:"primarykey"`
	CustomerID uint       `gorm:"not null;index;comment:关联客户ID"`
	UserID     uint       `gorm:"not null;comment:操作人ID"`
	ActionType string     `gorm:"size:50;not null;comment:操作类型:create/update/delete/query"`
	OldContent string     `gorm:"type:text;comment:修改前内容(JSON)"`
	NewContent string     `gorm:"type:text;comment:修改后内容(JSON)"`
	IP         string     `gorm:"size:50;comment:操作IP"`
	Remark     string     `gorm:"size:500;comment:备注"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime;comment:操作时间"`
}

type From_customer_add struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string `json:"title"`
	Phones    []struct {
		Prefix    string `json:"prefix"`
		Phone     string `json:"phone"`
		Label     string `json:"label"`
		IsPrimary bool   `json:"is_primary"`
	} `json:"phones"`
	Emails []struct {
		Email     string `json:"email"`
		Label     string `json:"label"`
		IsPrimary bool   `json:"is_primary"`
	} `json:"emails"`
	Companies []struct {
		CompanyName string `json:"company_name"`
		Department  string `json:"department"`
		Position    string `json:"position"`
		IsPrimary   bool   `json:"is_primary"`
	} `json:"companies"`
}

type From_customer_update struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string `json:"title"`
	Phones    []struct {
		Prefix    string `json:"prefix"`
		Phone     string `json:"phone"`
		Label     string `json:"label"`
		IsPrimary bool   `json:"is_primary"`
	} `json:"phones"`
	Emails []struct {
		Email     string `json:"email"`
		Label     string `json:"label"`
		IsPrimary bool   `json:"is_primary"`
	} `json:"emails"`
	Companies []struct {
		CompanyName string `json:"company_name"`
		Department  string `json:"department"`
		Position    string `json:"position"`
		IsPrimary   bool   `json:"is_primary"`
	} `json:"companies"`
}

// CustomerUpdateAdminsCash 更新客户管理员缓存
func CustomerUpdateAdminsCash() {
	customerAdmins = nil
	customerAdmins = append(customerAdmins, 1) // id=1 系统管理员默认拥有所有权限
	var binds []TabUserGroupBinds
	models.DB.Where("group_id = ?", customerUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(customerAdmins, item.UserID) {
			customerAdmins = append(customerAdmins, item.UserID)
		}
	}
}

// customerAdminCheck 检查用户是否为客户管理员
func customerAdminCheck(userID uint) bool {
	return slices.Contains(customerAdmins, userID)
}

// canModifyCustomer 判断是否有权限修改/删除客户（创建者或管理员）
func canModifyCustomer(userID, creatorUserID uint) bool {
	if slices.Contains(customerAdmins, userID) {
		return true
	}
	return userID == creatorUserID
}

// ApiCustomerInit 初始化客户模块
func ApiCustomerInit() {
	// 自动创建 customer_admin 用户组
	models.DB.Where("name = ?", "customer_admin").FirstOrCreate(&customerUserGroup, TabUserGroups{
		Name: "customer_admin",
		Type: "usergroup",
	})

	// 自动迁移客户相关表
	models.DB.AutoMigrate(
		&TabCustomer{},
		&TabCustomerPhone{},
		&TabCustomerEmail{},
		&TabCustomerCompany{},
		&TabCustomerLog{},
	)

	CustomerUpdateAdminsCash()
}

func ApiCustomer(r *gin.RouterGroup) {
	// POST /add - 新增客户（需 customer_admin 权限）
	r.POST("/add", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}
		if !customerAdminCheck(user.ID) {
			ReturnJson(ctx, "permission_denied", nil)
			return
		}

		var params From_customer_add
		if err := decodeJSON(data, &params); err != nil {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 创建客户
		customer := TabCustomer{
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Title:     params.Title,
			CreatedBy: user.ID,
		}
		if err := models.DB.Create(&customer).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}

		// 写入电话
		for _, p := range params.Phones {
			models.DB.Create(&TabCustomerPhone{
				CustomerID: customer.ID,
				Prefix:     p.Prefix,
				Phone:      p.Phone,
				Label:      p.Label,
				IsPrimary:  p.IsPrimary,
			})
		}

		// 写入邮箱
		for _, e := range params.Emails {
			models.DB.Create(&TabCustomerEmail{
				CustomerID: customer.ID,
				Email:      e.Email,
				Label:      e.Label,
				IsPrimary:  e.IsPrimary,
			})
		}

		// 写入单位
		for _, c := range params.Companies {
			models.DB.Create(&TabCustomerCompany{
				CustomerID:  customer.ID,
				CompanyName: c.CompanyName,
				Department:  c.Department,
				Position:    c.Position,
				IsPrimary:   c.IsPrimary,
			})
		}

		// 写日志
		models.DB.Create(&TabCustomerLog{
			CustomerID: customer.ID,
			UserID:     user.ID,
			ActionType: "create",
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", gin.H{"id": customer.ID})
	})

	// POST /update - 编辑客户（创建者或管理员可操作）
	r.POST("/update", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		var params From_customer_update
		if err := decodeJSON(data, &params); err != nil {
			ReturnJson(ctx, "parameErr", nil)
			return
		}
		if params.ID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 查找客户
		var customer TabCustomer
		if err := models.DB.Unscoped().First(&customer, params.ID).Error; err != nil {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 权限校验：只有创建者或管理员可以修改
		if !canModifyCustomer(user.ID, customer.CreatedBy) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		// 更新主表
		models.DB.Model(&customer).Updates(map[string]interface{}{
			"first_name": params.FirstName,
			"last_name":  params.LastName,
			"title":      params.Title,
		})

		// 重建绑定表：删除旧的，写入新的
		models.DB.Where("customer_id = ?", customer.ID).Delete(&TabCustomerPhone{})
		for _, p := range params.Phones {
			models.DB.Create(&TabCustomerPhone{
				CustomerID: customer.ID,
				Prefix:     p.Prefix,
				Phone:      p.Phone,
				Label:      p.Label,
				IsPrimary:  p.IsPrimary,
			})
		}

		models.DB.Where("customer_id = ?", customer.ID).Delete(&TabCustomerEmail{})
		for _, e := range params.Emails {
			models.DB.Create(&TabCustomerEmail{
				CustomerID: customer.ID,
				Email:      e.Email,
				Label:      e.Label,
				IsPrimary:  e.IsPrimary,
			})
		}

		models.DB.Where("customer_id = ?", customer.ID).Delete(&TabCustomerCompany{})
		for _, c := range params.Companies {
			models.DB.Create(&TabCustomerCompany{
				CustomerID:  customer.ID,
				CompanyName: c.CompanyName,
				Department:  c.Department,
				Position:    c.Position,
				IsPrimary:   c.IsPrimary,
			})
		}

		// 写日志
		models.DB.Create(&TabCustomerLog{
			CustomerID: customer.ID,
			UserID:     user.ID,
			ActionType: "update",
			OldContent: "update", // 简化处理
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// POST /delete - 软删除客户（创建者或管理员可操作）
	r.POST("/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		type Req struct {
			ID uint `json:"id"`
		}
		var req Req
		if err := decodeJSON(data, &req); err != nil || req.ID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		var customer TabCustomer
		if err := models.DB.First(&customer, req.ID).Error; err != nil {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 权限校验：只有创建者或管理员可以删除
		if !canModifyCustomer(user.ID, customer.CreatedBy) {
			ReturnJson(ctx, "no_permission", nil)
			return
		}

		models.DB.Delete(&customer)

		// 写日志
		models.DB.Create(&TabCustomerLog{
			CustomerID: customer.ID,
			UserID:     user.ID,
			ActionType: "delete",
			IP:         ctx.ClientIP(),
		})

		ReturnJson(ctx, "apiOK", nil)
	})

	// POST /list - 客户列表（登录用户可读）
	r.POST("/list", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		type Req struct {
			Page     int    `json:"page"`
			PageSize int    `json:"page_size"`
			Search   string `json:"search"`
		}
		var req Req
		if err := decodeJSON(data, &req); err != nil {
			req.Page = 1
			req.PageSize = 20
		}
		if req.Page < 1 {
			req.Page = 1
		}
		if req.PageSize < 1 || req.PageSize > 100 {
			req.PageSize = 20
		}

		offset := (req.Page - 1) * req.PageSize

		// 子查询：获取每个客户的主电话/主邮箱/主单位用于搜索和展示
		query := models.DB.Model(&TabCustomer{})

		if req.Search != "" {
			search := "%" + req.Search + "%"
			query = query.Where(
				models.DB.Where("first_name LIKE ? OR last_name LIKE ?", search, search).
					Or("id IN (?)", models.DB.Table("tab_customer_phones").Select("customer_id").Where("phone LIKE ?", search)).
					Or("id IN (?)", models.DB.Table("tab_customer_emails").Select("customer_id").Where("email LIKE ?", search)).
					Or("id IN (?)", models.DB.Table("tab_customer_companies").Select("customer_id").Where("company_name LIKE ?", search)),
			)
		}

		var total int64
		query.Count(&total)

		var customers []TabCustomer
		query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&customers)

		// 组装列表数据（含主电话/主邮箱/主单位 + 编辑权限）
		type CustomerListItem struct {
			TabCustomer
			PrimaryPhone   string `json:"primary_phone"`
			PrimaryEmail   string `json:"primary_email"`
			PrimaryCompany string `json:"primary_company"`
			Edit           bool   `json:"edit"`
		}

		var list []CustomerListItem
		for _, c := range customers {
			item := CustomerListItem{
				TabCustomer: c,
				Edit:        canModifyCustomer(user.ID, c.CreatedBy),
			}

			var phone TabCustomerPhone
			if err := models.DB.Where("customer_id = ? AND is_primary = ?", c.ID, true).First(&phone).Error; err == nil {
				item.PrimaryPhone = phone.Phone
			} else if err := models.DB.Where("customer_id = ?", c.ID).First(&phone).Error; err == nil {
				item.PrimaryPhone = phone.Phone
			}

			var email TabCustomerEmail
			if err := models.DB.Where("customer_id = ? AND is_primary = ?", c.ID, true).First(&email).Error; err == nil {
				item.PrimaryEmail = email.Email
			} else if err := models.DB.Where("customer_id = ?", c.ID).First(&email).Error; err == nil {
				item.PrimaryEmail = email.Email
			}

			var company TabCustomerCompany
			if err := models.DB.Where("customer_id = ? AND is_primary = ?", c.ID, true).First(&company).Error; err == nil {
				item.PrimaryCompany = company.CompanyName
			} else if err := models.DB.Where("customer_id = ?", c.ID).First(&company).Error; err == nil {
				item.PrimaryCompany = company.CompanyName
			}

			list = append(list, item)
		}

		ReturnJson(ctx, "apiOK", gin.H{
			"customers": list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	})

	// POST /get - 获取客户详情（登录用户可读）
	r.POST("/get", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userNoLogin", nil)
			return
		}

		type Req struct {
			ID uint `json:"id"`
		}
		var req Req
		if err := decodeJSON(data, &req); err != nil || req.ID == 0 {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		var customer TabCustomer
		if err := models.DB.First(&customer, req.ID).Error; err != nil {
			ReturnJson(ctx, "parameErr", nil)
			return
		}

		// 获取电话列表
		var phones []TabCustomerPhone
		models.DB.Where("customer_id = ?", req.ID).Find(&phones)

		// 获取邮箱列表
		var emails []TabCustomerEmail
		models.DB.Where("customer_id = ?", req.ID).Find(&emails)

		// 获取单位列表
		var companies []TabCustomerCompany
		models.DB.Where("customer_id = ?", req.ID).Find(&companies)

		// 写查询日志
		// models.DB.Create(&TabCustomerLog{
		// 	CustomerID: req.ID,
		// 	ActionType: "query",
		// 	IP:         ctx.ClientIP(),
		// })

		ReturnJson(ctx, "apiOK", gin.H{
			"customer":  customer,
			"phones":    phones,
			"emails":    emails,
			"companies": companies,
			"canModify": canModifyCustomer(user.ID, customer.CreatedBy),
		})
	})
}
