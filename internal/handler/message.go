// internal/handler/message.go
package handler

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	"ginchat/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 拉取历史消息（离线消息也靠这个）
func GetMessageHistory(c *gin.Context) {
	userID := c.GetUint("userID")

	friendID, err := strconv.Atoi(c.Query("friendId"))
	if err != nil || friendID == 0 {
		response.Fail(c, "参数错误")
		return
	}

	var messages []model.PrivateMessage
	database.DB.Where(
		"(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
		userID, friendID, friendID, userID,
	).Order("created_at ASC").Find(&messages)

	// 把对方发给自己的消息标记为已读
	now := time.Now().Unix()
	database.DB.Model(&model.PrivateMessage{}).
		Where("from_id = ? AND to_id = ? AND read_at IS NULL", friendID, userID).
		Update("read_at", now)

	response.OK(c, messages)
}

// 查询未读消息数（可选，首页展示用）
func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	var count int64
	database.DB.Model(&model.PrivateMessage{}).
		Where("to_id = ? AND read_at IS NULL", userID).
		Count(&count)

	response.OK(c, gin.H{"unread": count})
}
