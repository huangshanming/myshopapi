package dao

import (
	"gorm.io/gorm"
	"mymall/common"
	"mymall/db"
	"mymall/models"
	"mymall/utils"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type ProductDAO struct{}

var ProductDao = &ProductDAO{}

func (d *ProductDAO) GetList(page *common.PageReq) map[string]interface{} {
	var res = make(map[string]interface{})
	res["total"] = 0
	res["data"] = []string{}
	var total int64
	if err := db.GormDB.Model(&model.Product{}).Count(&total).Error; err != nil {
		return res
	}
	if total == 0 {
		return res
	}
	gormDb := db.GormDB.Model(&model.Product{}).Where("status = ?", "on_sale")
	result, err := utils.Paginate[model.ProductListResp](gormDb, page)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res
		}
		return res
	}
	res["data"] = result
	return res
}

func (d *ProductDAO) GetDetail(id uint64) *model.Product {
	var product model.Product
	if err := db.GormDB.Model(&model.Product{}).Where("id = ?", id).First(&product).Error; err != nil {
		return nil
	}
	return &product
}
