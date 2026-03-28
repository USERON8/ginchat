// internal/handler/friend.go
package handler

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	"ginchat/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ApplyFriend(c *gin.Context) {
	userID := c.GetUint("userID")

	// 从路由参数取
	friendIDStr := c.Param("id")
	friendIDInt, err := strconv.Atoi(friendIDStr)
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}
	friendID := uint(friendIDInt)

	if userID == friendID {
		response.Fail(c, response.CodeCannotAddSelf)
		return
	}

	var target model.User
	if err := database.DB.First(&target, friendID).Error; err != nil {
		response.Fail(c, response.CodeUserNotFound)
		return
	}

	var exist model.Friendship
	if database.DB.Where(
		"user_id = ? AND friend_id = ?", userID, friendID,
	).First(&exist).Error == nil {
		response.Fail(c, response.CodeFriendExist)
		return
	}

	database.DB.Create(&model.Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   "pending",
	})

	response.OK(c, nil)
}

// 同意/拒绝申请
func HandleApply(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		ApplyID uint   `json:"applyId" binding:"required"`
		Action  string `json:"action" binding:"required"` // accepted / rejected
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	var friendship model.Friendship
	if err := database.DB.First(&friendship, req.ApplyID).Error; err != nil {
		response.Fail(c, response.CodeFriendNotFound)
		return
	}

	// 只有被申请人才能操作
	if friendship.FriendID != userID {
		response.Fail(c, response.CodeFriendNoPermission)
		return
	}

	if req.Action == "accepted" {
		// 同意：双向建立好友关系
		database.DB.Model(&friendship).Update("status", "accepted")
		// 反向也插一条，这样查好友列表只需要查 user_id = 自己
		reverse := model.Friendship{
			UserID:   userID,
			FriendID: friendship.UserID,
			Status:   "accepted",
		}
		database.DB.Create(&reverse)
	} else {
		// 拒绝：直接删除申请
		database.DB.Delete(&friendship)
	}

	response.OK(c, nil)
}

// 查看收到的好友申请
func GetApplyList(c *gin.Context) {
	userID := c.GetUint("userID")

	var applies []model.Friendship
	database.DB.Where("friend_id = ? AND status = ?", userID, "pending").
		Find(&applies)

	response.OK(c, applies)
}

// 好友列表
func GetFriendList(c *gin.Context) {
	userID := c.GetUint("userID")

	var friendships []model.Friendship
	database.DB.Where("user_id = ? AND status = ?", userID, "accepted").
		Find(&friendships)
	// 没有好友直接返回空数组
	if len(friendships) == 0 {
		response.OK(c, []model.User{})
		return
	}
	// 拿到所有好友 ID
	var friendIDs []uint
	for _, f := range friendships {
		friendIDs = append(friendIDs, f.FriendID)
	}

	// 查好友的用户信息
	var friends []model.User
	database.DB.Select("id, username, phone, email").
		Where("id IN ?", friendIDs).
		Find(&friends)

	response.OK(c, friends)
}

// 删除好友
func DeleteFriend(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		FriendID uint `json:"friendId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	// 双向删除
	database.DB.Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, req.FriendID, req.FriendID, userID,
	).Delete(&model.Friendship{})

	response.OK(c, nil)
}
