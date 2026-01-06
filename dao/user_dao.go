package dao

import (
	"crypto/md5"
	"mymall/db"
	"mymall/models"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type UserDAO struct{}

var UserDao = &UserDAO{}

func (d *UserDAO) GetInfo(mobile string, password string) model.User {
	user := model.User{}
	db.GormDB.Model(&model.User{}).Where("mobile = ?", mobile).Find(&user)
	md5Password := md5.New().Sum([]byte(password + "this is my mall"))
	if string(md5Password) != user.Password {
		return model.User{}
	}
	return user
}
