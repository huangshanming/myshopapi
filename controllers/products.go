package controllers

import (
	"fmt"
	"mymall/dao"
	"mymall/pkg/pagination"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductsController struct {
	BaseController
}

func (ctrl *ProductsController) GetList(c *gin.Context) {
	pageReq, _ := c.Get("pageReq")
	req := pageReq.(*pagination.PageReq)
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
