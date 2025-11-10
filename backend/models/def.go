package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// 获取当前时间字符串
// 参数格式可选，默认"2006-01-02 15:04:05"
func GetCurrentTimeString(format ...string) string {
	// 默认格式
	layout := "2006_01_02-15_04_05.999999999"

	// 如果传入了格式参数则使用自定义格式
	if len(format) > 0 {
		layout = format[0]
	}

	return time.Now().Format(layout)
}

func RandStr32() string {
	// 生成 32 字节 (256 位) 随机数据
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// 转换为 16 进制字符串 (长度 64)
	cookie := hex.EncodeToString(b)
	return cookie
}

func Md5Str(str string) string {
	hashBytes2 := md5.Sum([]byte(str))
	hashString2 := hex.EncodeToString(hashBytes2[:]) // 注意数组转切片的[:]
	return hashString2
}

func HashUserPass(str string) string {
	switch ConfigsUser.PassHashType {
	case "text":
		return str
	case "md5":
		return Md5Str(str)
	}

	return GetCurrentTimeString() + RandStr32() //如果转换失败返回当前时间，避免撞库
}
