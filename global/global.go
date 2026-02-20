package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *ConfigStruct // 全局配置
	DB     *gorm.DB      // 数据库实例
	RDB    *redis.Client // Redis实例
	Log    *zap.Logger   // 日志实例
)
