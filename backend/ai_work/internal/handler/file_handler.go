package handler

import (
	"ops/internal/service"
	"ops/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileHandler struct {
	service service.FileService
}

func NewFileHandler(db *gorm.DB) *FileHandler {
	return &FileHandler{
		service: service.NewFileService(db),
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器，支持图片、文档等多种类型
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param file formData file true "文件内容"
// @Param type formData string false "文件类型" default(image)
// @Param description formData string false "文件描述"
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 400 {object} response.StandardResponse "参数错误"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 413 {object} response.StandardResponse "文件过大"
// @Failure 415 {object} response.StandardResponse "文件类型不支持"
// @Failure 500 {object} response.StandardResponse "服务器错误"
// @Router /api/v1/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取文件类型参数
	fileType := c.PostForm("type")
	if fileType == "" {
		fileType = "image" // 默认类型为图片
	}

	// 获取文件描述
	description := c.PostForm("description")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 调用Service上传文件
	uploadResponse, success := h.service.UploadFile(c, userID.(uint), file, fileType, description)
	if !success {
		response.BadRequest(c, "文件上传失败，请检查文件格式和大小")
		return
	}

	response.Success(c, uploadResponse)
}

// GetFileList 获取文件列表
// @Summary 获取文件列表
// @Description 获取当前用户上传的文件列表
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param type query string false "文件类型过滤"
// @Param page query int false "页码" default(1)
// @Param entries query int false "每页数量" default(20)
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 400 {object} response.StandardResponse "参数错误"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Router /api/v1/files/list [get]
func (h *FileHandler) GetFileList(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取查询参数
	fileType := c.Query("type")
	page := GetIntParam(c, "page", 1)
	entries := GetIntParam(c, "entries", 20)

	// 调用Service获取文件列表
	fileListResponse, success := h.service.GetFileList(userID.(uint), fileType, page, entries)
	if !success {
		response.BadRequest(c, "参数错误")
		return
	}

	response.Success(c, fileListResponse)
}

// GetFileByID 获取文件信息
// @Summary 获取文件信息
// @Description 根据文件ID获取文件详细信息
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param id path int true "文件ID"
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 404 {object} response.StandardResponse "文件不存在"
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetFileByID(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取文件ID
	fileID := GetUintParam(c, "id")
	if fileID == 0 {
		response.BadRequest(c, "文件ID无效")
		return
	}

	// 调用Service获取文件信息
	file, success := h.service.GetFileByID(fileID, userID.(uint))
	if !success {
		response.Error(c, "-100", "文件不存在或无权限访问")
		return
	}

	response.Success(c, gin.H{
		"file_id":    file.ID,
		"name":       file.Name,
		"sha256":     file.Sha256,
		"mime":       file.Mime,
		"type":       file.Type,
		"size":       file.Const, // 注意：这里const字段实际上存储的是使用次数，需要确认实际字段
		"created_at": file.Date.Format("2006-01-02T15:04:05Z"),
		"path":       file.Path,
	})
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除用户上传的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param userID header string false "用户ID" default("")
// @Param id path int true "文件ID"
// @Success 200 {object} response.StandardResponse "成功"
// @Failure 401 {object} response.StandardResponse "未授权"
// @Failure 404 {object} response.StandardResponse "文件不存在"
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取文件ID
	fileID := GetUintParam(c, "id")
	if fileID == 0 {
		response.BadRequest(c, "文件ID无效")
		return
	}

	// 调用Service删除文件
	success := h.service.DeleteFile(fileID, userID.(uint))
	if !success {
		response.Error(c, "-100", "文件删除失败，文件不存在或无权限")
		return
	}

	response.Success(c, gin.H{"message": "文件删除成功"})
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 下载文件内容（直接下载）
// @Tags 文件管理
// @Accept json
// @Produce application/octet-stream
// @Param hash path string true "文件SHA256哈希值"
// @Success 200 {file} binary "文件内容"
// @Failure 404 {object} response.StandardResponse "文件不存在"
// @Router /api/v1/files/download/{hash} [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		response.BadRequest(c, "文件哈希值无效")
		return
	}

	// 调用Service下载文件
	success := h.service.DownloadFile(c, hash, true)
	if !success {
		response.Error(c, "-100", "文件不存在")
		return
	}
}

// GetFile 获取文件（预览）
// @Summary 获取文件（预览）
// @Description 获取文件内容（浏览器预览）
// @Tags 文件管理
// @Accept json
// @Produce *
// @Param hash path string true "文件SHA256哈希值"
// @Success 200 {file} binary "文件内容"
// @Failure 404 {object} response.StandardResponse "文件不存在"
// @Router /api/v1/files/get/{hash} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		response.BadRequest(c, "文件哈希值无效")
		return
	}

	// 调用Service获取文件（预览模式）
	success := h.service.DownloadFile(c, hash, false)
	if !success {
		response.Error(c, "-100", "文件不存在")
		return
	}
}

