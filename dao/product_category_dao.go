package dao

import (
	"fmt"
	"gorm.io/gorm"
	"mymall/common"
	"mymall/db"
	"mymall/models"
	"mymall/utils"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type ProductCategoryDAO struct{}

var ProductCategoryDao = &ProductCategoryDAO{}

func (d *ProductCategoryDAO) GetList(page *common.PageReq) *common.PageRes[model.ProductCategory] {
	var res = new(common.PageRes[model.ProductCategory])
	res.Total = 0
	var total int64
	if err := db.GormDB.Model(&model.ProductCategory{}).Count(&total).Error; err != nil {
		return res
	}
	fmt.Println(total)
	if total == 0 {
		return res
	}
	gormDb := db.GormDB.Model(&model.ProductCategory{}).Where("is_show = ?", "1")
	result, err := utils.Paginate[model.ProductCategory](gormDb, page)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res
		}
		return res
	}
	res = result
	return res
}

func (d *ProductCategoryDAO) GetDetail(id uint64) *model.ProductCategory {
	var productCategory model.ProductCategory
	if err := db.GormDB.Model(&model.ProductCategory{}).Where("id = ?", id).First(&productCategory).Error; err != nil {
		return nil
	}
	return &productCategory
}
