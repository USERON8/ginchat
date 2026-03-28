// pkg/redis/online.go
package redis

import (
	"context"
	"fmt"
)

func SetOnline(userID uint) {
	RDB.Set(context.Background(),
		fmt.Sprintf("online:%d", userID),
		1,
		0, // 不过期，下线时手动删
	)
}

func SetOffline(userID uint) {
	RDB.Del(context.Background(),
		fmt.Sprintf("online:%d", userID),
	)
}

func IsOnline(userID uint) bool {
	val, err := RDB.Exists(context.Background(),
		fmt.Sprintf("online:%d", userID),
	).Result()
	return err == nil && val > 0
}
