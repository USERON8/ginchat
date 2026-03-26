// pkg/database/db.go
package database

import (
	"fmt"
	"ginchat/internal/model"
	"ginchat/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	cfg := config.Cfg.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动建表
	DB.AutoMigrate(&model.User{}, &model.PrivateMessage{}, &model.Friendship{})
}
