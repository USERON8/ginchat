// pkg/redis/token.go
package redis

import (
	"context"
	"fmt"
	"time"
)

func SetRefreshToken(userID uint, token string) error {
	return RDB.Set(context.Background(),
		fmt.Sprintf("refresh_token:%d", userID),
		token,
		7*24*time.Hour,
	).Err()
}

func GetRefreshToken(userID uint) (string, error) {
	return RDB.Get(context.Background(),
		fmt.Sprintf("refresh_token:%d", userID),
	).Result()
}

func DelRefreshToken(userID uint) error {
	return RDB.Del(context.Background(),
		fmt.Sprintf("refresh_token:%d", userID),
	).Err()
}
