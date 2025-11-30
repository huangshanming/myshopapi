package controllers
import utils

type GoodsController struct {
}

func (ctrl *GoodsController) GetGoodsList(c *gin.Context) {
	utils.SuccessResponse(c, "hello world")
}
