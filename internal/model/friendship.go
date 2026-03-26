package model

import "gorm.io/gorm"

type Friendship struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	FriendID uint   `gorm:"not null"`
	Status   string `gorm:"default:'pending'"` // pending/accepted
}
