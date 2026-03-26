// pkg/redis/user_cache.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// key 格式：user:info:1
func userKey(userID uint) string {
	return fmt.Sprintf("user:info:%d", userID)
}

// 写缓存
func SetUserInfo(userID uint, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return RDB.Set(context.Background(),
		userKey(userID),
		bytes,
		30*time.Minute, // 30分钟过期
	).Err()
}

// 读缓存
func GetUserInfo(userID uint, dest interface{}) error {
	val, err := RDB.Get(context.Background(), userKey(userID)).Bytes()
	if err != nil {
		return err // redis.Nil 说明缓存不存在
	}
	return json.Unmarshal(val, dest)
}

// 删缓存（更新用户信息时调用）
func DelUserInfo(userID uint) error {
	return RDB.Del(context.Background(), userKey(userID)).Err()
}
