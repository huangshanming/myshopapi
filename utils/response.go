package utils

import "github.com/gin-gonic/gin"
struct Response {
  Code int `json:"code"`
  Msg string `json:"msg"`
  Data interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{} , msg string = "success") {
	c.JSON(200, Response {
		Code: 200,
		Msg: "success",
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, msg string, code int = 400) {
	c.JSON(200, Response {
		Code: code,
		Msg: msg,
		Data: nil,
	})
}