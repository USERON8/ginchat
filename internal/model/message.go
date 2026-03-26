package model

import "gorm.io/gorm"

// 私聊消息
type PrivateMessage struct {
	gorm.Model
	FromID  uint   `gorm:"not null"`
	ToID    uint   `gorm:"not null"`
	Content string `gorm:"not null"`
	Type    string `gorm:"default:'text'"` // text/image
	ReadAt  *int64 // 为 nil 表示未读
}
