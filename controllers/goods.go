package controllers

import "mymall/utils"
import "github.com/gin-gonic/gin"
import "mymall/dao"
import "net/http"

type GoodsController struct {
}

func (ctrl *GoodsController) GetGoodsList(c *gin.Context) {
	data := make(map[string]interface{})
	utils.SuccessResponse(c, data, "hello world")
	a, _ := dao.DemoDao.GetByID(1)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "查询成功", "data": a})
}
