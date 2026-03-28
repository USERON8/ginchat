// internal/model/group.go
package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name    string `gorm:"not null"`
	OwnerID uint   `gorm:"not null"`
	Avatar  string
	Notice  string
}

type GroupMember struct {
	gorm.Model
	GroupID uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
	Role    string `gorm:"default:'member'"` // owner/admin/member
}

type GroupMessage struct {
	gorm.Model
	GroupID uint   `gorm:"not null;index"`
	FromID  uint   `gorm:"not null"`
	Content string `gorm:"not null"`
	Type    string `gorm:"default:'text'"`
}
