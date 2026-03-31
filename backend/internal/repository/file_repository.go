package repository

import (
	"ops/models"

	"gorm.io/gorm"
)

type FileRepository interface {
	CreateFile(file *models.TabFileInfo_) error
	GetFileByID(fileID uint) (*models.TabFileInfo_, error)
	GetFileByHash(hash string) (*models.TabFileInfo_, error)
	GetFilesByUser(userID uint, fileType string, page, entries int) ([]models.TabFileInfo_, int64, error)
	UpdateFile(file *models.TabFileInfo_) error
	DeleteFile(fileID uint) error
	IncrementFileUsage(fileID uint) error
	GetFilesByType(fileType string, limit int) ([]models.TabFileInfo_, error)
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) CreateFile(file *models.TabFileInfo_) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) GetFileByID(fileID uint) (*models.TabFileInfo_, error) {
	var file models.TabFileInfo_
	if err := r.db.First(&file, fileID).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetFileByHash(hash string) (*models.TabFileInfo_, error) {
	var file models.TabFileInfo_
	if err := r.db.Where("sha256 = ?", hash).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetFilesByUser(userID uint, fileType string, page, entries int) ([]models.TabFileInfo_, int64, error) {
	var files []models.TabFileInfo_
	var total int64

	query := r.db.Model(&models.TabFileInfo_{}).Where("user_id = ?", userID)

	if fileType != "" {
		query = query.Where("type = ?", fileType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := entries * (page - 1)
	if err := query.Order("date DESC").Offset(offset).Limit(entries).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

func (r *fileRepository) UpdateFile(file *models.TabFileInfo_) error {
	return r.db.Save(file).Error
}

func (r *fileRepository) DeleteFile(fileID uint) error {
	return r.db.Delete(&models.TabFileInfo_{}, fileID).Error
}

func (r *fileRepository) IncrementFileUsage(fileID uint) error {
	return r.db.Model(&models.TabFileInfo_{}).
		Where("id = ?", fileID).
		Update("const", gorm.Expr("const + ?", 1)).Error
}

func (r *fileRepository) GetFilesByType(fileType string, limit int) ([]models.TabFileInfo_, error) {
	var files []models.TabFileInfo_
	query := r.db.Where("type = ?", fileType).Order("const DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}