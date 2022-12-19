package jwt

import (
	"log"
	"testing"
)

func TestAccessToken(t *testing.T) {
	// 由于测试相对路径问题，settings.Init()函数执行失败，只能手动设置
	//_ = settings.Init()
	//Init()
	secret = []byte("test")
	accessHour = 2
	tokenStr, _ := GetAccessToken(123, "123")
	payload, err := ParseToken(tokenStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("userID:%d,username:%s", payload.UserID, payload.Username)
	}
}
