package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

func Error(c *gin.Context, msg string, code int) {
	if code == 0 {
		code = 400
	}
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
