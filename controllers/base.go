package controllers

import (
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (ctrl *BaseController) Success(c *gin.Context, data interface{}, msg string) error {
	if msg == "" {
		msg = "success"
	}
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
	c.Abort()
	return nil
}

func (ctrl *BaseController) Error(c *gin.Context, msg string, code int) error {
	if code == 0 {
		code = 400
	}
	c.JSON(200, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
	return nil
}
