// internal/handler/group.go
package handler

import (
	"ginchat/internal/model"
	"ginchat/pkg/database"
	"ginchat/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建群
func CreateGroup(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Name   string `json:"name" binding:"required,min=2,max=20"`
		Notice string `json:"notice"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	group := model.Group{
		Name:    req.Name,
		OwnerID: userID,
		Notice:  req.Notice,
	}
	if err := database.DB.Create(&group).Error; err != nil {
		response.Fail(c, response.CodeServerError)
		return
	}

	// 创建者自动加入群，角色 owner
	database.DB.Create(&model.GroupMember{
		GroupID: group.ID,
		UserID:  userID,
		Role:    "owner",
	})

	response.OK(c, gin.H{"groupId": group.ID})
}

// 群详情
func GetGroup(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	var group model.Group
	if err := database.DB.First(&group, groupID).Error; err != nil {
		response.Fail(c, response.CodeGroupNotFound)
		return
	}

	response.OK(c, group)
}

// 我的群列表
func GetGroupList(c *gin.Context) {
	userID := c.GetUint("userID")

	var members []model.GroupMember
	database.DB.Where("user_id = ?", userID).Find(&members)

	if len(members) == 0 {
		response.OK(c, []model.Group{})
		return
	}

	var groupIDs []uint
	for _, m := range members {
		groupIDs = append(groupIDs, m.GroupID)
	}

	var groups []model.Group
	database.DB.Where("id IN ?", groupIDs).Find(&groups)

	response.OK(c, groups)
}

// 邀请成员
func AddGroupMember(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	var req struct {
		UserID uint `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	// 检查操作人是否在群里
	var self model.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&self).Error; err != nil {
		response.Fail(c, response.CodeGroupNoPermission)
		return
	}

	// 检查目标用户是否已在群里
	var exist model.GroupMember
	if database.DB.Where("group_id = ? AND user_id = ?", groupID, req.UserID).
		First(&exist).Error == nil {
		response.Fail(c, response.CodeGroupMemberExist)
		return
	}

	database.DB.Create(&model.GroupMember{
		GroupID: uint(groupID),
		UserID:  req.UserID,
		Role:    "member",
	})

	response.OK(c, nil)
}

// 踢出成员（只有 owner/admin 能操作）
func RemoveGroupMember(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	var req struct {
		UserID uint `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	// 检查操作人权限
	var self model.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&self).Error; err != nil {
		response.Fail(c, response.CodeGroupNoPermission)
		return
	}
	if self.Role == "member" {
		response.Fail(c, response.CodeGroupNoPermission)
		return
	}

	database.DB.Where("group_id = ? AND user_id = ?", groupID, req.UserID).
		Delete(&model.GroupMember{})

	response.OK(c, nil)
}

// 成员列表
func GetGroupMembers(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	var members []model.GroupMember
	database.DB.Where("group_id = ?", groupID).Find(&members)

	if len(members) == 0 {
		response.OK(c, []model.User{})
		return
	}

	var userIDs []uint
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}

	var users []model.User
	database.DB.Select("id, username, phone, email").
		Where("id IN ?", userIDs).Find(&users)

	response.OK(c, users)
}
