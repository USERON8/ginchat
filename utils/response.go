package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 状态码 200成功 其他失败
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 400,
		Msg:  msg,
		Data: nil,
	})
}
