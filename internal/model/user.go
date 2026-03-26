package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex;not null"`
	Password      string `gorm:"not null"`
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	LogOutTime    uint64
	IsLoggedIn    bool
	IsAdmin       bool
	DeviceInfo    string
	HeartbeatTime uint64
}

func (User) TableName() string {
	return "user"
}
