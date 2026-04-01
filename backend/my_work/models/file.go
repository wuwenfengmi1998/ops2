package models

import (
	"crypto"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// 计算文件的哈希
func SHA256HashFile(file_head *multipart.FileHeader) (string, error) {
	// 打开文件
	file, err := file_head.Open()
	if err != nil {
		return "foen error", err
	}
	defer file.Close()

	hasher := crypto.SHA256.New()

	// 从文件流中读取并计算哈希
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil

}

// 获取文件mime
func GetFileMime(file_head *multipart.FileHeader) (string, error) {
	file, err := file_head.Open()
	if err != nil {
		return "foen error", err
	}
	defer file.Close()

	// 读取前512字节用于MIME检测
	buffer := make([]byte, 512)
	io.ReadFull(file, buffer)
	mimeType := http.DetectContentType(buffer)

	return mimeType, nil

}
