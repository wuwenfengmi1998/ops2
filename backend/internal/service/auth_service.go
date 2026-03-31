package service

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ops/internal/database"
	"ops/internal/repository"
)

// AuthService 用户认证服务结构
type AuthService struct {
	userRepo     repository.UserRepository
	userInfoRepo repository.UserInfoRepository
	cookieRepo   repository.CookieRepository
	db           *gorm.DB
}

// UserWithInfo 用户信息结构
type UserWithInfo struct {
	UserID      uint   `json:"userID"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatarURL"`
	CookieValue string `json:"cookieValue"`
}

// CookieInfo Cookie信息结构
type CookieInfo struct {
	Value      string `json:"value"`
	ExpireDate time.Time `json:"expireDate"`
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:     repository.NewUserRepository(db),
		userInfoRepo: repository.NewUserInfoRepository(db),
		cookieRepo:   repository.NewCookieRepository(db),
		db:           db,
	}
}

// Login 用户登录
func (s *AuthService) Login(name, password, deviceID, ip, remember string) (*UserWithInfo, *CookieInfo, error) {
	if name == "" || password == "" {
		return nil, nil, errors.New("username and password are required")
	}

	// 查找用户
	user, err := s.userRepo.FindByName(name)
	if err != nil {
		return nil, nil, fmt.Errorf("find user error: %w", err)
	}
	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	// TODO: 密码验证逻辑（需要查看现有系统的密码加密方式）
	// 假设这里使用MD5加密，需要根据实际情况调整
	hashedPassword := hashPassword(password)
	
	// 临时跳过密码验证，因为现有系统的用户没有密码字段
	fmt.Printf("DEBUG: Trying to login user %s (password: %s, hashed: %s)\n", name, password, hashedPassword)

	// 生成Cookie
	cookieValue := generateCookieValue(user.ID, name, deviceID)
	
	// 设置过期时间
	expiresAt := time.Now()
	if remember == "1" || remember == "true" {
		expiresAt = expiresAt.Add(30 * 24 * time.Hour) // 30天
	} else {
		expiresAt = expiresAt.Add(24 * time.Hour) // 24小时
	}

	cookie := &database.TabCookie{
		Value:     cookieValue,
		UserID:    user.ID,
		ExpiresAt: expiresAt.Unix(),
		CreateAt:  time.Now().Unix(),
		Remember:  (remember == "1" || remember == "true"),
	}

	// 保存Cookie到数据库
	if err := s.cookieRepo.Create(cookie); err != nil {
		return nil, nil, fmt.Errorf("create cookie error: %w", err)
	}

	// 获取用户信息
	userInfo, err := s.userInfoRepo.FindByUserID(user.ID)
	if err != nil {
		fmt.Printf("WARN: user info not found for user %s: %v\n", name, err)
	}

	// 构建头像URL
	avatarURL := "/static/default_avatar.png"
	if userInfo != nil && userInfo.AvatarPath != "" {
		avatarURL = "/static/uploads/" + userInfo.AvatarPath
	}

	// 返回用户信息和Cookie
	userWithInfo := &UserWithInfo{
		UserID:      user.ID,
		Name:        user.Name,
		AvatarURL:   avatarURL,
		CookieValue: cookieValue,
	}

	cookieInfo := &CookieInfo{
		Value:      cookieValue,
		ExpireDate: expiresAt,
	}

	return userWithInfo, cookieInfo, nil
}

// Register 用户注册
func (s *AuthService) Register(name, password, email, phone string) (*UserWithInfo, *CookieInfo, error) {
	if name == "" || password == "" {
		return nil, nil, errors.New("username and password are required")
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByName(name)
	if err != nil {
		return nil, nil, fmt.Errorf("check username exists error: %w", err)
	}
	if exists {
		return nil, nil, errors.New("username already exists")
	}

	// 创建用户
	user := &database.TabUser{
		Name: name,
		// 注意：现有TabUser表只有ID和Name字段，没有密码字段
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("create user error: %w", err)
	}

	// 创建用户信息
	userInfo := &database.TabUserInfo{
		UserID:    user.ID,
		AvatarPath: "", // 默认空
		Birthdate: "",
		Gender:    0,
		Introduction: "",
	}

	if err := s.userInfoRepo.Create(userInfo); err != nil {
		// 如果创建用户信息失败，删除用户（可选）
		s.userRepo.Delete(user.ID)
		return nil, nil, fmt.Errorf("create user info error: %w", err)
	}

	// 生成Cookie
	cookieValue := generateCookieValue(user.ID, name, "register")
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天

	cookie := &database.TabCookie{
		Value:     cookieValue,
		UserID:    user.ID,
		ExpiresAt: expiresAt.Unix(),
		CreateAt:  time.Now().Unix(),
		Remember:  true,
	}

	if err := s.cookieRepo.Create(cookie); err != nil {
		return nil, nil, fmt.Errorf("create cookie error: %w", err)
	}

	// 返回用户信息和Cookie
	userWithInfo := &UserWithInfo{
		UserID:      user.ID,
		Name:        user.Name,
		AvatarURL:   "/static/default_avatar.png",
		CookieValue: cookieValue,
	}

	cookieInfo := &CookieInfo{
		Value:      cookieValue,
		ExpireDate: expiresAt,
	}

	return userWithInfo, cookieInfo, nil
}

// ForgotPassword 忘记密码
func (s *AuthService) ForgotPassword(name, email, phone string) (string, error) {
	if name == "" {
		return "", errors.New("username is required")
	}

	// 查找用户
	user, err := s.userRepo.FindByName(name)
	if err != nil {
		return "", fmt.Errorf("find user error: %w", err)
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// 生成重置令牌
	resetToken := generateResetToken(user.ID, name)

	// TODO: 发送重置密码邮件或短信
	// 这里应该实现邮件发送或短信发送逻辑

	fmt.Printf("DEBUG: Password reset token for user %s: %s\n", name, resetToken)

	return resetToken, nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(token, newPassword string) error {
	if token == "" || newPassword == "" {
		return errors.New("token and new password are required")
	}

	// TODO: 验证重置令牌并获取用户ID
	// 这里应该解析token获取用户ID
	userID := parseResetToken(token)
	if userID == 0 {
		return errors.New("invalid reset token")
	}

	// 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("find user error: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// TODO: 更新密码
	// 注意：现有TabUser表没有密码字段，这里可能需要扩展表结构或使用其他方式存储密码

	fmt.Printf("DEBUG: Password reset for user %s (ID: %d)\n", user.Name, user.ID)

	return nil
}

// Logout 用户退出登录
func (s *AuthService) Logout(cookieValue, deviceID string) error {
	if cookieValue == "" {
		return errors.New("cookie value is required")
	}

	return s.cookieRepo.DeleteByValue(cookieValue)
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(userID uint) (*UserWithInfo, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 获取增强的用户信息
	enhancedUser, err := repository.GetEnhancedUserInfo(s.db, userID)
	if err != nil {
		return nil, fmt.Errorf("get enhanced user info error: %w", err)
	}
	if enhancedUser == nil {
		return nil, errors.New("user not found")
	}

	return &UserWithInfo{
		UserID:    enhancedUser.TabUser.ID,
		Name:      enhancedUser.TabUser.Name,
		AvatarURL: enhancedUser.AvatarURL,
	}, nil
}

// UpdateProfile 更新用户信息
func (s *AuthService) UpdateProfile(userID uint, updateData map[string]interface{}) (*UserWithInfo, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 获取用户信息
	enhancedUser, err := repository.GetEnhancedUserInfo(s.db, userID)
	if err != nil {
		return nil, fmt.Errorf("get enhanced user info error: %w", err)
	}
	if enhancedUser == nil {
		return nil, errors.New("user not found")
	}

	// 更新用户信息
	// 检查是否有avatar字段
	if avatarPath, ok := updateData["avatar"]; ok {
		avatarStr, isString := avatarPath.(string)
		if isString && avatarStr != "" {
			enhancedUser.UserInfo.AvatarPath = avatarStr
			if err := s.userInfoRepo.Update(&enhancedUser.UserInfo); err != nil {
				return nil, fmt.Errorf("update user info error: %w", err)
			}
		}
	}

	// 检查其他可更新字段
	if gender, ok := updateData["gender"]; ok {
		if genderNum, isNum := gender.(float64); isNum {
			enhancedUser.UserInfo.Gender = int(genderNum)
		}
	}

	if birthdate, ok := updateData["birthdate"]; ok {
		if birthdateStr, isString := birthdate.(string); isString {
			enhancedUser.UserInfo.Birthdate = birthdateStr
		}
	}

	if intro, ok := updateData["introduction"]; ok {
		if introStr, isString := intro.(string); isString {
			enhancedUser.UserInfo.Introduction = introStr
		}
	}

	// 保存更新后的用户信息
	if err := s.userInfoRepo.Update(&enhancedUser.UserInfo); err != nil {
		return nil, fmt.Errorf("update user info error: %w", err)
	}

	// 构建头像URL
	avatarURL := "/static/default_avatar.png"
	if enhancedUser.UserInfo.AvatarPath != "" {
		avatarURL = "/static/uploads/" + enhancedUser.UserInfo.AvatarPath
	}

	return &UserWithInfo{
		UserID:    userID,
		Name:      enhancedUser.TabUser.Name,
		AvatarURL: avatarURL,
	}, nil
}

// ValidateCookie 验证Cookie有效性
func (s *AuthService) ValidateCookie(cookieValue string) (uint, error) {
	if cookieValue == "" {
		return 0, errors.New("cookie value is required")
	}

	cookie, err := s.cookieRepo.FindByValue(cookieValue)
	if err != nil {
		return 0, fmt.Errorf("find cookie error: %w", err)
	}
	if cookie == nil {
		return 0, errors.New("cookie not found")
	}

	// 检查是否过期
	if cookie.ExpiresAt < time.Now().Unix() {
		// 删除过期的Cookie
		s.cookieRepo.DeleteByValue(cookieValue)
		return 0, errors.New("cookie expired")
	}

	return cookie.UserID, nil
}

// 辅助函数
func hashPassword(password string) string {
	// 使用MD5哈希（根据现有系统可能使用其他方式）
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateCookieValue(userID uint, username, deviceID string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%d%s%s%d", userID, username, deviceID, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateResetToken(userID uint, username string) string {
	timestamp := time.Now().UnixNano()
	random := fmt.Sprintf("%d", timestamp)
	data := fmt.Sprintf("%d%s%s%d", userID, username, random, timestamp)
	hash := sha256.Sum256([]byte(data))
	token := hex.EncodeToString(hash[:])
	
	// 存储到数据库或Redis（这里简化处理）
	// 在实际应用中应该存储token并设置过期时间
	return token
}

func parseResetToken(token string) uint {
	// 简化的token解析，实际应该从数据库或Redis验证
	// 这里返回0表示无效
	if len(token) < 32 {
		return 0
	}
	
	// TODO: 实现token解析逻辑
	// 暂时返回0，需要根据具体token格式实现
	return 0
}

// CleanupExpiredCookies 清理过期Cookie
func (s *AuthService) CleanupExpiredCookies() error {
	return s.cookieRepo.DeleteExpired()
}

// GetUserByCookie 通过Cookie获取用户信息
func (s *AuthService) GetUserByCookie(cookieValue string) (*UserWithInfo, error) {
	userID, err := s.ValidateCookie(cookieValue)
	if err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}

// UpdateUserPassword 更新用户密码
func (s *AuthService) UpdateUserPassword(userID uint, oldPassword, newPassword string) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}

	if oldPassword == "" || newPassword == "" {
		return errors.New("old password and new password are required")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("find user error: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// TODO: 验证旧密码
	// 现有系统没有密码字段，需要扩展
	
	// TODO: 更新密码
	// 现有系统没有密码字段，需要扩展

	return errors.New("password update not supported in current schema")
}