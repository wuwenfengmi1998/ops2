package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ops/internal/database"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(user *database.TabUser) error
	FindByID(id uint) (*database.TabUser, error)
	FindByName(name string) (*database.TabUser, error)
	FindByEmail(email string) (*database.TabUser, error)
	FindByPhone(phone string) (*database.TabUser, error)
	Update(user *database.TabUser) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
}

// UserInfoRepository 用户信息数据访问接口
type UserInfoRepository interface {
	Create(userInfo *database.TabUserInfo) error
	FindByUserID(userID uint) (*database.TabUserInfo, error)
	Update(userInfo *database.TabUserInfo) error
	Delete(userID uint) error
}

// CookieRepository Cookie数据访问接口
type CookieRepository interface {
	Create(cookie *database.TabCookie) error
	FindByValue(cookieValue string) (*database.TabCookie, error)
	FindByUserID(userID uint) ([]*database.TabCookie, error)
	DeleteByValue(cookieValue string) error
	DeleteByUserID(userID uint) error
	DeleteExpired() error
}

// userRepo 用户仓库实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Create 创建用户
func (r *userRepo) Create(user *database.TabUser) error {
	if user == nil {
		return errors.New("user is nil")
	}
	
	if user.Name == "" {
		return errors.New("username is required")
	}
	
	return r.db.Create(user).Error
}

// FindByID 通过ID查找用户
func (r *userRepo) FindByID(id uint) (*database.TabUser, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}
	
	var user database.TabUser
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// FindByName 通过用户名查找用户
func (r *userRepo) FindByName(name string) (*database.TabUser, error) {
	if name == "" {
		return nil, errors.New("username is required")
	}
	
	var user database.TabUser
	err := r.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// FindByEmail 通过邮箱查找用户
func (r *userRepo) FindByEmail(email string) (*database.TabUser, error) {
	// TabUser表目前没有email字段，这里返回nil
	return nil, nil
}

// FindByPhone 通过手机号查找用户
func (r *userRepo) FindByPhone(phone string) (*database.TabUser, error) {
	// TabUser表目前没有phone字段，这里返回nil
	return nil, nil
}

// Update 更新用户信息
func (r *userRepo) Update(user *database.TabUser) error {
	if user == nil {
		return errors.New("user is nil")
	}
	
	if user.ID == 0 {
		return errors.New("user ID is required")
	}
	
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepo) Delete(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}
	
	return r.db.Delete(&database.TabUser{}, id).Error
}

// ExistsByName 检查用户名是否存在
func (r *userRepo) ExistsByName(name string) (bool, error) {
	if name == "" {
		return false, errors.New("username is required")
	}
	
	var count int64
	err := r.db.Model(&database.TabUser{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// userInfoRepo 用户信息仓库实现
type userInfoRepo struct {
	db *gorm.DB
}

// NewUserInfoRepository 创建用户信息仓库实例
func NewUserInfoRepository(db *gorm.DB) UserInfoRepository {
	return &userInfoRepo{db: db}
}

// Create 创建用户信息
func (r *userInfoRepo) Create(userInfo *database.TabUserInfo) error {
	if userInfo == nil {
		return errors.New("user info is nil")
	}
	
	if userInfo.UserID == 0 {
		return errors.New("user ID is required")
	}
	
	return r.db.Create(userInfo).Error
}

// FindByUserID 通过用户ID查找用户信息
func (r *userInfoRepo) FindByUserID(userID uint) (*database.TabUserInfo, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	
	var userInfo database.TabUserInfo
	err := r.db.Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &userInfo, nil
}

// Update 更新用户信息
func (r *userInfoRepo) Update(userInfo *database.TabUserInfo) error {
	if userInfo == nil {
		return errors.New("user info is nil")
	}
	
	if userInfo.UserID == 0 {
		return errors.New("user ID is required")
	}
	
	return r.db.Save(userInfo).Error
}

// Delete 删除用户信息
func (r *userInfoRepo) Delete(userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	
	return r.db.Where("user_id = ?", userID).Delete(&database.TabUserInfo{}).Error
}

// cookieRepo Cookie仓库实现
type cookieRepo struct {
	db *gorm.DB
}

// NewCookieRepository 创建Cookie仓库实例
func NewCookieRepository(db *gorm.DB) CookieRepository {
	return &cookieRepo{db: db}
}

// Create 创建Cookie
func (r *cookieRepo) Create(cookie *database.TabCookie) error {
	if cookie == nil {
		return errors.New("cookie is nil")
	}
	
	if cookie.Value == "" {
		return errors.New("cookie value is required")
	}
	
	if cookie.UserID == 0 {
		return errors.New("user ID is required")
	}
	
	if cookie.ExpiresAt == 0 {
		cookie.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	
	if cookie.CreateAt == 0 {
		cookie.CreateAt = time.Now().Unix()
	}
	
	return r.db.Create(cookie).Error
}

// FindByValue 通过Cookie值查找
func (r *cookieRepo) FindByValue(cookieValue string) (*database.TabCookie, error) {
	if cookieValue == "" {
		return nil, errors.New("cookie value is required")
	}
	
	var cookie database.TabCookie
	err := r.db.Where("value = ?", cookieValue).First(&cookie).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &cookie, nil
}

// FindByUserID 通过用户ID查找所有Cookie
func (r *cookieRepo) FindByUserID(userID uint) ([]*database.TabCookie, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	
	var cookies []*database.TabCookie
	err := r.db.Where("user_id = ?", userID).Find(&cookies).Error
	if err != nil {
		return nil, err
	}
	
	return cookies, nil
}

// DeleteByValue 通过Cookie值删除
func (r *cookieRepo) DeleteByValue(cookieValue string) error {
	if cookieValue == "" {
		return errors.New("cookie value is required")
	}
	
	return r.db.Where("value = ?", cookieValue).Delete(&database.TabCookie{}).Error
}

// DeleteByUserID 通过用户ID删除所有Cookie
func (r *cookieRepo) DeleteByUserID(userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	
	return r.db.Where("user_id = ?", userID).Delete(&database.TabCookie{}).Error
}

// DeleteExpired 删除过期的Cookie
func (r *cookieRepo) DeleteExpired() error {
	now := time.Now().Unix()
	return r.db.Where("expires_at < ?", now).Delete(&database.TabCookie{}).Error
}

// EnhancedUserInfo 增强的用户信息结构
type EnhancedUserInfo struct {
	database.TabUser
	UserInfo database.TabUserInfo
	AvatarURL string
}

// GetEnhancedUserInfo 获取增强的用户信息
func GetEnhancedUserInfo(db *gorm.DB, userID uint) (*EnhancedUserInfo, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	
	var user database.TabUser
	var userInfo database.TabUserInfo
	
	// 获取用户基本信息
	err := db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	// 获取用户详细信息
	err = db.Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// 构建头像URL
	avatarURL := "/static/default_avatar.png"
	if userInfo.AvatarPath != "" {
		avatarURL = "/file/" + userInfo.AvatarPath
	}
	
	return &EnhancedUserInfo{
		TabUser:   user,
		UserInfo:  userInfo,
		AvatarURL: avatarURL,
	}, nil
}