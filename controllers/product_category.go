package controllers

import "github.com/gin-gonic/gin"
import "mymall/dao"
import "mymall/common"
import "strconv"

type ProductCategoryController struct {
	BaseController
}

func (ctrl *ProductCategoryController) GetList(c *gin.Context) {
	pageReq, _ := c.Get("pageReq")
	req := pageReq.(*common.PageReq)
	data := dao.ProductCategoryDao.GetList(req)
	ctrl.Success(c, data, "查询成功")
	return
}

func (ctrl *ProductCategoryController) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 10)
	data := dao.ProductCategoryDao.GetDetail(id)
	ctrl.Success(c, data, "查询成功")
	return
}
