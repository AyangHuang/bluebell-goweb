package jwt

import (
	"bluebell/settings"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Payload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const (
	AccessToken int = iota
	RefreshToken
)

var (
	secret         []byte
	accessHours    int
	refreshHours   int
	ErrorTokenType error = errors.New("tokenType err, 0 or 1 is ok")
)

func Init() {
	secret = []byte(settings.Conf.TokenConfig.Secret)
	accessHours = settings.Conf.TokenConfig.AccessHour
	refreshHours = settings.Conf.TokenConfig.RefreshHour
}

func GetTowToken(userID int64, username string) (accessToken, refreshToken string, err error) {
	accessToken, err = GetToken(userID, username, AccessToken)
	if err == nil {
		refreshToken, err = GetToken(userID, username, RefreshToken)
	}
	return
}

// GetToken 返回一个 Token
func GetToken(userID int64, username string, tokenType int) (string, error) {
	var times int
	if tokenType == AccessToken {
		times = accessHours
	} else if tokenType == RefreshToken {
		times = refreshHours
	} else {
		return "", ErrorTokenType
	}
	claims := Payload{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发人
			Issuer: "bluebell",
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(times) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if claims, ok := token.Claims.(*Payload); ok {
		if token.Valid {
			return claims, nil
		}
		return nil, errors.New("invalid token")
	} else {
		return nil, err
	}

}
