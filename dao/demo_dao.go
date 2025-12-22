package dao

import (
	"gorm.io/gorm"
	"mymall/db"
	"mymall/models"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type DemoDAO struct{}

var DemoDao = &DemoDAO{}

// GetByID 基于 GORM 根据 ID 查询单条用户数据（替代原生 QueryRow）
func (d *DemoDAO) GetByID(id int) (*model.Demo, error) {
	var demo model.Demo
	// GORM 链式调用：Model指定模型 -> Where条件 -> First查询第一条
	result := db.GormDB.Model(&model.Demo{}).Where("id = ?", id).First(&demo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 无数据返回 nil
		}
		return nil, result.Error // 查询异常返回错误
	}
	return &demo, nil
}

//
//// GetList 基于 GORM 实现条件+分页查询（替代原生 Query）
//// 参数：age 年龄筛选（0表示不筛选）；pageNum 页码；pageSize 每页条数
//func (d *UserDAO) GetList(age int, pageNum, pageSize int) ([]*model.User, int64, error) {
//	var userList []*model.User
//	var total int64
//
//	// 1. 构建查询条件
//	query := db.GormDB.Model(&model.User{})
//	if age > 0 {
//		query = query.Where("age = ?", age) // 动态添加年龄条件
//	}
//
//	// 2. 统计总条数
//	if err := query.Count(&total).Error; err != nil {
//		return nil, 0, err
//	}
//
//	// 3. 分页查询（Offset：偏移量，Limit：每页条数）
//	offset := (pageNum - 1) * pageSize
//	if err := query.Offset(offset).Limit(pageSize).Find(&userList).Error; err != nil {
//		return nil, 0, err
//	}
//
//	return userList, total, nil
//}
//
//// 额外扩展：GORM 条件查询示例（模糊查询用户名）
//func (d *UserDAO) SearchByName(name string) ([]*model.User, error) {
//	var userList []*model.User
//	// 模糊查询：LIKE %name%
//	err := db.GormDB.Where("name LIKE ?", "%"+name+"%").Find(&userList).Error
//	return userList, err
//}
