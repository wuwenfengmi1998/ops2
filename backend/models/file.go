package models

import "os"

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
