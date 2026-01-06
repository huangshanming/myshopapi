package controllers

import (
	"github.com/gin-gonic/gin"
	"mymall/dao"
)

type UsersController struct {
	BaseController
}

func (ctrl *UsersController) Login(c *gin.Context) {
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")
	data := dao.UserDao.GetInfo(mobile, password)
	if data.ID == 0 {
		ctrl.Error(c, "检验失败", 0)
		return

	}
	ctrl.Success(c, data, "查询成功")
	return

}
