package controllers

import (
	"mymall/pkg/response"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (ctrl *BaseController) Success(c *gin.Context, data interface{}, msg string) error {
	response.Success(c, data, msg)
	c.Abort()
	return nil
}

func (ctrl *BaseController) Error(c *gin.Context, msg string, code int) error {
	response.Error(c, msg, code)
	c.Abort()
	return nil
}
