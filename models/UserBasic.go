package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Username      string
	PassWord      string
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

func (UserBasic) TableName() string {
	return "user_basic"
}
