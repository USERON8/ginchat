// internal/handler/user.go
package handler

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	jwtpkg "ginchat/pkg/jwt"
	"ginchat/pkg/redis"
	"ginchat/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeParamError, err.Error())
		return
	}

	// 检查用户名是否已存在
	var exist model.User
	if database.DB.Where("username = ?", req.Username).First(&exist).Error == nil {
		response.Fail(c, response.CodeUserExist)
		return
	}

	// 密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	user := model.User{Username: req.Username, Password: string(hash)}
	if err := database.DB.Create(&user).Error; err != nil {
		response.Fail(c, response.CodeRegisterFail)
		return
	}

	response.OK(c, gin.H{"id": user.ID})
}

func Login(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeParamError, err.Error())
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Fail(c, response.CodeUserExist)
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Fail(c, response.CodePasswordError)
		return
	}

	// 记录登录信息（你的字段）
	now := uint64(time.Now().Unix())
	database.DB.Model(&user).Updates(map[string]interface{}{
		"client_ip":    c.ClientIP(),
		"login_time":   now,
		"is_logged_in": true,
	})
	accessToken, err := jwtpkg.GenerateAccessToken(user.ID)
	if err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(user.ID)
	if err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	// refresh token 存 Redis
	err = redis.SetRefreshToken(user.ID, refreshToken)
	if err != nil {
		return
	}

	response.OK(c, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"userID":       user.ID,
	})
}
func Logout(c *gin.Context) {
	userID := c.GetUint("userID")

	now := uint64(time.Now().Unix())
	database.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_logged_in": false,
			"log_out_time": now,
		})

	// 删 refresh token，让刷新失效
	redis.DelRefreshToken(userID)

	response.OK(c, nil)
}

// handler/user.go 新增
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	// 解析 refresh token
	claims, err := jwtpkg.Parse(req.RefreshToken)
	if err != nil || claims.TokenType != "refresh" {
		response.Fail(c, response.CodeTokenInvalid)
		return
	}

	// 校验 Redis 里是否存在（退出登录后就没了）
	stored, err := redis.GetRefreshToken(claims.UserID)
	if err != nil || stored != req.RefreshToken {
		response.Fail(c, response.CodeTokenInvalid)
		return
	}

	// 生成新的 access token
	newAccessToken, err := jwtpkg.GenerateAccessToken(claims.UserID)
	if err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	response.OK(c, gin.H{"accessToken": newAccessToken})
}
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	// 定义返回结构（不能把整个 User 缓存，密码会泄漏）
	type UserVO struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}

	var vo UserVO

	// 1. 先查 Redis
	if err := redis.GetUserInfo(userID, &vo); err == nil {
		// 缓存命中，直接返回
		response.OK(c, vo)
		return
	}

	// 2. 缓存没有，查数据库
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.Fail(c, response.CodeUserExist)
		return
	}

	vo = UserVO{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Email:    user.Email,
	}

	// 3. 写入 Redis
	err := redis.SetUserInfo(userID, vo)
	if err != nil {
		return
	}

	response.OK(c, vo)
}
func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeParamError, err.Error())
		return
	}

	if err := database.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"phone": req.Phone,
			"email": req.Email,
		}).Error; err != nil {
		response.Fail(c, response.CodeUpdateFail)
		return
	}

	// 删掉旧缓存，下次查询自动重建
	err := redis.DelUserInfo(userID)
	if err != nil {
		return
	}

	response.OK(c, nil)
}
func SearchUser(c *gin.Context) {
	username := c.Query("username") // /api/user/search?username=xxx
	if username == "" {
		response.Fail(c, response.CodeServerError)
		return
	}

	var users []model.User
	database.DB.Where("username LIKE ?", "%"+username+"%").
		Select("id, username, phone, email"). // 不返回密码！
		Find(&users)

	response.OK(c, users)
}

// handler/user.go 新增
func ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	// 查当前用户
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.Fail(c, response.CodeUserNotFound)
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(req.OldPassword),
	); err != nil {
		response.Fail(c, response.CodePasswordError)
		return
	}

	// 新密码不能和旧密码一样
	if req.OldPassword == req.NewPassword {
		response.FailWithMsg(c, response.CodeParamError, "新密码不能与旧密码相同")
		return
	}

	// 加密新密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	// 更新密码
	database.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", string(hash))

	// 删缓存、删 refresh token（强制重新登录）
	redis.DelUserInfo(userID)
	redis.DelRefreshToken(userID)

	response.OK(c, nil)
}
