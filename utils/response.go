package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, msg string, code int) {
	if code == 0 {
		code = 400
	}
	c.JSON(200, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
