package global

import (
	"fmt"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	var dialector gorm.Dialector

	// 1. 根据配置选择驱动
	switch Config.Database.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			Config.Database.Host,
			Config.Database.Port,
			Config.Database.Username,
			Config.Database.Password,
			Config.Database.DBName,
			Config.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", Config.Database.Type)
	}

	// 2. GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境打印SQL，上线改成logger.Error
	}

	// 3. 打开连接
	if DB, err = gorm.Open(dialector, gormConfig); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 4. 配置连接池（高并发IM必备）
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(Config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(Config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.Database.ConnMaxLifetime) * time.Second)

	Log.Info("数据库初始化成功")
	return nil
}
