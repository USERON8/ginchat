package api

import (
	"ginchat/dao"
	"ginchat/models"
	"ginchat/utils"

	"github.com/gin-gonic/gin"
)

var userDao = dao.NewBaseDao[models.UserBasic]()

// GetUserInfo 获取用户信息接口
func GetUserInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "用户ID不能为空")
		return
	}

	// 这里简化处理，实际项目要做类型转换
	user, err := userDao.GetByID(c.Request.Context(), 1)
	if err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}
