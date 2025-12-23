package controllers

import "github.com/gin-gonic/gin"
import "mymall/dao"
import "net/http"

type ProductsController struct {
}

func (ctrl *ProductsController) GetList(c *gin.Context) {
	data := dao.ProductDao.GetList(1)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "查询成功", "data": data})
}
