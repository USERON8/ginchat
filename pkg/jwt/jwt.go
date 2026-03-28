// pkg/jwt/jwt.go
package jwt

import (
	"errors"
	"ginchat/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint
	TokenType string // access / refresh
	jwt.RegisteredClaims
}

// 生成 access token（2小时）
func GenerateAccessToken(userID uint) (string, error) {
	return generate(userID, "access", 2*time.Hour)
}

// 生成 refresh token（7天）
func GenerateRefreshToken(userID uint) (string, error) {
	return generate(userID, "refresh", 7*24*time.Hour)
}

func generate(userID uint, tokenType string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.Cfg.JWT.Secret))
}

func Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token 无效")
	}
	return token.Claims.(*Claims), nil
}
