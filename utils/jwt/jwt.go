package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"web_app/settings"
)

type AccessPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	secret     = []byte(settings.Conf.TokenConfig.Secret)
	accessHour = settings.Conf.TokenConfig.AccessHour
)

// GetAccessToken 返回一个 accessToken
func GetAccessToken(userID int64, username string) (string, error) {

	claims := AccessPayload{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发人
			Issuer: "bluebell",
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessHour) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*AccessPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessPayload{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if claims, ok := token.Claims.(*AccessPayload); ok {
		if token.Valid {
			return claims, nil
		}
		return nil, errors.New("invalid token")
	} else {
		return nil, err
	}

}
