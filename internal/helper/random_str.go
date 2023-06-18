package helper

import (
	"math/rand"
	"time"
)

func RandomStr(n int) string {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子

	// 可选的字符集合
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机字符串
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
