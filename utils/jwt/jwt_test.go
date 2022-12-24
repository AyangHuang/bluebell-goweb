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
	accessHours = 2
	refreshHours = 3
	ac, re, _ := GetTowToken(123, "ayang")
	acP, _ := ParseToken(ac)
	reP, _ := ParseToken(re)
	log.Printf("ac :userID:%d,username:%s\nre :userID:%d,username:%s",
		acP.UserID, acP.Username,
		reP.UserID, reP.Username)
}
