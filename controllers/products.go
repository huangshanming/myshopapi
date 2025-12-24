package controllers

import "github.com/gin-gonic/gin"
import "mymall/dao"
import "net/http"
import "mymall/common"

type ProductsController struct {
}

func (ctrl *ProductsController) GetList(c *gin.Context) {
	pageReq, _ := c.Get("pageReq")
	req := pageReq.(*common.PageReq)
	data := dao.ProductDao.GetList(req)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "查询成功", "data": data})
}
