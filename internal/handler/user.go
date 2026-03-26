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
		response.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var exist model.User
	if database.DB.Where("username = ?", req.Username).First(&exist).Error == nil {
		response.Fail(c, "用户名已存在")
		return
	}

	// 密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "服务器错误")
		return
	}

	user := model.User{Username: req.Username, Password: string(hash)}
	if err := database.DB.Create(&user).Error; err != nil {
		response.Fail(c, "注册失败")
		return
	}

	response.OK(c, gin.H{"id": user.ID})
}

func Login(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Fail(c, "用户不存在")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Fail(c, "密码错误")
		return
	}

	// 记录登录信息（你的字段）
	now := uint64(time.Now().Unix())
	database.DB.Model(&user).Updates(map[string]interface{}{
		"client_ip":    c.ClientIP(),
		"login_time":   now,
		"is_logged_in": true,
	})
	token, err := jwtpkg.Generate(user.ID)
	if err != nil {
		response.Fail(c, "生成token失败")
		return
	}
	response.OK(c, gin.H{"token": token, "userID": user.ID})

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

	response.OK(c, nil)
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
		response.Fail(c, "用户不存在")
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
		response.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"phone": req.Phone,
			"email": req.Email,
		}).Error; err != nil {
		response.Fail(c, "更新失败")
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
		response.Fail(c, "请输入用户名")
		return
	}

	var users []model.User
	database.DB.Where("username LIKE ?", "%"+username+"%").
		Select("id, username, phone, email"). // 不返回密码！
		Find(&users)

	response.OK(c, users)
}
