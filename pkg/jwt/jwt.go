// pkg/jwt/jwt.go
package jwt

import (
	"errors"
	"ginchat/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func Generate(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(config.Cfg.JWT.Expire) * time.Hour),
			),
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
