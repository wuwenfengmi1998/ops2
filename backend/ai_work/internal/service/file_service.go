package service

import (
	"mime"
	"mime/multipart"
	"ops/internal/repository"
	"ops/models"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileService interface {
	UploadFile(c *gin.Context, userID uint, fileHeader *multipart.FileHeader, fileType, description string) (UploadResponse, bool)
	GetFileList(userID uint, fileType string, page, entries int) (FileListResponse, bool)
	GetFileByID(fileID uint, userID uint) (*models.TabFileInfo_, bool)
	GetFileByHash(hash string) (*models.TabFileInfo_, bool)
	DeleteFile(fileID uint, userID uint) bool
	DownloadFile(c *gin.Context, hash string, download bool) bool
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(db *gorm.DB) FileService {
	return &fileService{
		repo: repository.NewFileRepository(db),
	}
}

// 响应结构体
type UploadResponse struct {
	FileID      uint   `json:"file_id"`
	Name        string `json:"name"`
	SHA256      string `json:"sha256"`
	Mime        string `json:"mime"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url"`
	PreviewURL  string `json:"preview_url"`
	CreatedAt   string `json:"created_at"`
}

type FileListResponse struct {
	Files []FileInfo `json:"files"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Pages int        `json:"pages"`
}

type FileInfo struct {
	FileID    uint   `json:"file_id"`
	Name      string `json:"name"`
	SHA256    string `json:"sha256"`
	Mime      string `json:"mime"`
	Size      int64  `json:"size"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

func (s *fileService) UploadFile(c *gin.Context, userID uint, fileHeader *multipart.FileHeader, fileType, description string) (UploadResponse, bool) {
	// 验证文件大小
	if fileHeader.Size > int64(models.ConfigsFile.MaxSize) {
		return UploadResponse{}, false
	}

	// 验证文件最小大小
	if fileHeader.Size < 512 {
		return UploadResponse{}, false
	}

	// 验证文件名
	if fileHeader.Filename == "" {
		return UploadResponse{}, false
	}

	// 安全处理文件名
	filename := filepath.Base(fileHeader.Filename)

	// 计算文件哈希
	hashStr, err := models.SHA256HashFile(fileHeader)
	if err != nil {
		return UploadResponse{}, false
	}

	// 获取文件MIME类型
	mimeType, err := models.GetFileMime(fileHeader)
	if err != nil {
		return UploadResponse{}, false
	}

	// 验证MIME类型（如果是图片）
	if fileType == "image" {
		if models.ConfigsFile.AllowImageMime[mimeType] == "" {
			return UploadResponse{}, false
		}
	}

	// 构建文件保存路径
	var savePath string
	switch fileType {
	case "image":
		savePath = filepath.Join(models.ConfigsFile.Pahts["image"], hashStr)
	default:
		savePath = filepath.Join(models.ConfigsFile.Pahts["default"], hashStr)
	}

	// 检查文件是否已存在
	if models.FileExists(savePath) {
		// 如果文件已存在，增加使用计数
		existingFile, err := s.repo.GetFileByHash(hashStr)
		if err == nil && existingFile != nil {
			s.repo.IncrementFileUsage(existingFile.ID)
		}
	} else {
		// 保存文件到磁盘
		if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
			return UploadResponse{}, false
		}
	}

	// 检查数据库中是否已存在该文件
	existingFile, _ := s.repo.GetFileByHash(hashStr)
	if existingFile != nil {
		// 更新使用计数
		s.repo.IncrementFileUsage(existingFile.ID)
		
		return UploadResponse{
			FileID:      existingFile.ID,
			Name:        filename,
			SHA256:      hashStr,
			Mime:        mimeType,
			Size:        fileHeader.Size,
			DownloadURL: "/api/v1/files/download/" + hashStr,
			PreviewURL:  "/api/v1/files/get/" + hashStr,
			CreatedAt:   existingFile.Date.Format("2006-01-02T15:04:05Z"),
		}, true
	}

	// 创建新的文件记录
	newFile := &models.TabFileInfo_{
		Name:   filename,
		Path:   savePath,
		Sha256: hashStr,
		Mime:   mimeType,
		Type:   fileType,
		UserID: userID,
		Date:   time.Now(),
	}

	if err := s.repo.CreateFile(newFile); err != nil {
		return UploadResponse{}, false
	}

	return UploadResponse{
		FileID:      newFile.ID,
		Name:        filename,
		SHA256:      hashStr,
		Mime:        mimeType,
		Size:        fileHeader.Size,
		DownloadURL: "/api/v1/files/download/" + hashStr,
		PreviewURL:  "/api/v1/files/get/" + hashStr,
		CreatedAt:   newFile.Date.Format("2006-01-02T15:04:05Z"),
	}, true
}

func (s *fileService) GetFileList(userID uint, fileType string, page, entries int) (FileListResponse, bool) {
	// 验证分页参数
	if entries <= 0 || entries > 100 {
		return FileListResponse{}, false
	}
	if page <= 0 {
		return FileListResponse{}, false
	}

	files, total, err := s.repo.GetFilesByUser(userID, fileType, page, entries)
	if err != nil {
		return FileListResponse{}, false
	}

	// 计算总页数
	pages := int(total) / entries
	if int(total)%entries > 0 {
		pages++
	}

	// 转换文件信息
	fileInfos := make([]FileInfo, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			FileID:    file.ID,
			Name:      file.Name,
			SHA256:    file.Sha256,
			Mime:      file.Mime,
			Type:      file.Type,
			CreatedAt: file.Date.Format("2006-01-02T15:04:05Z"),
		})
	}

	return FileListResponse{
		Files: fileInfos,
		Total: total,
		Page:  page,
		Pages: pages,
	}, true
}

func (s *fileService) GetFileByID(fileID uint, userID uint) (*models.TabFileInfo_, bool) {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return nil, false
	}

	// 检查文件所有权
	if file.UserID != userID {
		return nil, false
	}

	return file, true
}

func (s *fileService) GetFileByHash(hash string) (*models.TabFileInfo_, bool) {
	file, err := s.repo.GetFileByHash(hash)
	if err != nil {
		return nil, false
	}
	return file, true
}

func (s *fileService) DeleteFile(fileID uint, userID uint) bool {
	// 首先检查文件所有权
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return false
	}

	if file.UserID != userID {
		return false
	}

	// 删除文件记录
	if err := s.repo.DeleteFile(fileID); err != nil {
		return false
	}

	// 注意：这里不删除物理文件，因为可能还有其他引用
	// 如果需要删除物理文件，需要检查引用计数

	return true
}

func (s *fileService) DownloadFile(c *gin.Context, hash string, download bool) bool {
	file, err := s.repo.GetFileByHash(hash)
	if err != nil {
		return false
	}

	// 检查文件是否存在
	if !models.FileExists(file.Path) {
		return false
	}

	// 设置响应头
	if download {
		// 下载模式
		c.Header("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	} else {
		// 预览模式
		ext := filepath.Ext(file.Name)
		if ext != "" {
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				c.Header("Content-Type", mimeType)
			}
		}
	}

	c.Header("Content-Type", "application/octet-stream")
	c.File(file.Path)

	// 增加使用计数
	s.repo.IncrementFileUsage(file.ID)

	return true
}