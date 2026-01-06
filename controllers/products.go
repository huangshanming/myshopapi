package controllers

import "github.com/gin-gonic/gin"
import "mymall/dao"
import "mymall/common"
import "strconv"
import "fmt"

type ProductsController struct {
	BaseController
}

func (ctrl *ProductsController) GetList(c *gin.Context) {
	pageReq, _ := c.Get("pageReq")
	req := pageReq.(*common.PageReq)
	data := dao.ProductDao.GetList(req)
	ctrl.Success(c, data, "查询成功")
	return
}

func (ctrl *ProductsController) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 10)
	fmt.Println(id)
	data := dao.ProductDao.GetDetail(id)
	ctrl.Success(c, data, "查询成功")
	return
}
