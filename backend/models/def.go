package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"regexp"

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

func HashUserPass(user *TabUser_) {
	switch ConfigsUser.PassHashType {
	case "text":
		break
	case "md5":
		user.Pass = Md5Str(user.Pass)

	case "md5salt":
		if user.Salt == "" {
			user.Salt = RandStr32()
		}
		user.Pass = Md5Str(Md5Str(user.Pass) + user.Salt)

	}

}

func IsExpired(expireTime time.Time) bool {
	return expireTime.Before(time.Now())
}

func CheckCookiesAndUpdate(cookie *TabCookie_) bool {
	if !IsExpired(cookie.ExpiresAt) {
		if cookie.Remember {
			cookiewhere := TabCookie_{
				ID: cookie.ID,
			}
			cookieupdata := TabCookie_{
				UpdatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Duration(ConfigsUser.CookieTimeout) * time.Second),
			}
			DB.Where(&cookiewhere).Updates(&cookieupdata)

		}
		return true
	} else {
		//以过期
		return false
	}
	//return false
}

// 判断邮箱是否合法
func IsEmailValid(email string) bool {
	// 正则表达式（覆盖 99% 常见邮箱格式）
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
