package dao

import (
	"fmt"
	"gorm.io/gorm"
	"mymall/db"
	"mymall/models"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type ProductDAO struct{}

var ProductDao = &ProductDAO{}

func (d *ProductDAO) GetList(page int) map[string]interface{} {
	var product model.Product
	pagesize := 20
	offset := (page - 1) * pagesize
	var res = make(map[string]interface{})
	res["total"] = 0
	res["data"] = []string{}
	var total int64
	if err := db.GormDB.Model(&model.Product{}).Count(&total).Error; err != nil {
		return res
	}
	fmt.Println(total)
	if total == 0 {
		return res
	}
	result := db.GormDB.Model(&model.Product{}).Where("status = ?", "on_sale").Limit(pagesize).Offset(offset).Find(&product)
	fmt.Println(result.Error)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return res
		}
		return res
	}
	res["data"] = product
	return res
}
