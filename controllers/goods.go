package controllers

import "mymall/utils"
import "github.com/gin-gonic/gin"

type GoodsController struct {
}

func (ctrl *GoodsController) GetGoodsList(c *gin.Context) {
	data := make(map[string]interface{})
	utils.SuccessResponse(c, data, "hello world")
}
