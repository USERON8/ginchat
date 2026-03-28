// pkg/redis/dedup.go
package redis

import (
	"context"
	"fmt"
	"time"
)

func IsDuplicate(msgID string) bool {
	key := fmt.Sprintf("msg:dedup:%s", msgID)

	// SetNX：key 不存在才设置
	// 返回 true = 设置成功 = 第一次见到这个消息 = 不重复
	// 返回 false = key 已存在 = 重复消息
	ok, err := RDB.SetNX(context.Background(),
		key,
		1,
		60*time.Second, // 60秒内同一个 msgId 只处理一次
	).Result()
	if err != nil {
		return false // Redis 出错，放行消息
	}
	return !ok
}
