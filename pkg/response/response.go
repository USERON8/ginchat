// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: 0, Msg: "success", Data: data})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Result{Code: -1, Msg: msg, Data: nil})
}
