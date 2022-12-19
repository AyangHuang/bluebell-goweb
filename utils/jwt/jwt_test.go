package jwt

import (
	"log"
	"testing"
	"web_app/settings"
)

func TestAccessToken(t *testing.T) {
	_ = settings.Init()
	tokenStr, _ := GetAccessToken(123, "123")
	payload, err := ParseToken(tokenStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("userID:%d,username:%s", payload.UserID, payload.Username)
	}
}
