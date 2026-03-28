// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code Code        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 成功
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// 失败（用错误码）
func Fail(c *gin.Context, code Code) {
	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// 失败（自定义消息，比如 binding 校验的详细错误）
func FailWithMsg(c *gin.Context, code Code, msg string) {
	c.JSON(http.StatusOK, Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
