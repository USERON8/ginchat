// internal/handler/message.go 替换
package handler

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	"ginchat/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 私聊历史消息（游标分页）
func GetMessageHistory(c *gin.Context) {
	userID := c.GetUint("userID")

	friendID, err := strconv.Atoi(c.Query("friendId"))
	if err != nil || friendID == 0 {
		response.Fail(c, response.CodeParamError)
		return
	}

	// lastId=0 表示第一次加载，取最新的
	lastID, _ := strconv.Atoi(c.DefaultQuery("lastId", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	query := database.DB.Where(
		"(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
		userID, friendID, friendID, userID,
	).Order("id DESC").Limit(size)

	// 游标：只取比 lastId 更小的（更早的消息）
	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}

	var messages []model.PrivateMessage
	query.Find(&messages)

	// 标记已读
	now := time.Now().Unix()
	database.DB.Model(&model.PrivateMessage{}).
		Where("from_id = ? AND to_id = ? AND read_at IS NULL", friendID, userID).
		Update("read_at", now)

	response.OK(c, gin.H{
		"list":    messages,
		"hasMore": len(messages) == size, // 返回的数量等于 size，说明可能还有更多
	})
}

// 标记已读
func MarkRead(c *gin.Context) {
	userID := c.GetUint("userID")

	friendID, err := strconv.Atoi(c.Query("friendId"))
	if err != nil || friendID == 0 {
		response.Fail(c, response.CodeParamError)
		return
	}

	now := time.Now().Unix()
	database.DB.Model(&model.PrivateMessage{}).
		Where("from_id = ? AND to_id = ? AND read_at IS NULL", friendID, userID).
		Update("read_at", now)

	response.OK(c, nil)
}

// 未读消息数
func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	var count int64
	database.DB.Model(&model.PrivateMessage{}).
		Where("to_id = ? AND read_at IS NULL", userID).
		Count(&count)

	response.OK(c, gin.H{"unread": count})
}

// 群聊历史消息（游标分页）
func GetGroupMessageHistory(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	lastID, _ := strconv.Atoi(c.DefaultQuery("lastId", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	query := database.DB.Where("group_id = ?", groupID).
		Order("id DESC").Limit(size)

	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}

	var messages []model.GroupMessage
	query.Find(&messages)

	response.OK(c, gin.H{
		"list":    messages,
		"hasMore": len(messages) == size,
	})
}
