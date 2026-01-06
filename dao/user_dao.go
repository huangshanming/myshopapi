package dao

import (
	"crypto/md5"
	"encoding/hex"
	"mymall/db"
	"mymall/models"
)

// UserDAO 仍保持原有结构，仅修改内部实现
type UserDAO struct{}

var UserDao = &UserDAO{}

func (d *UserDAO) GetInfo(mobile string, password string) model.User {
	user := model.User{}
	db.GormDB.Model(&model.User{}).Where("mobile = ?", mobile).Find(&user)
	md5Password := d.GetMd5Password(password)
	if string(md5Password) != user.Password {
		return model.User{}
	}
	return user
}

func (d *UserDAO) GetMd5Password(password string) string {
	hash := md5.New()
	hash.Write([]byte(password + "this is my mall")) // 加盐拼接后写入
	// 2. 计算哈希值（二进制）并编码为十六进制字符串
	md5Password := hex.EncodeToString(hash.Sum(nil))
	return md5Password
}

func (d *UserDAO) Register(mobile string, password string) (model.User, string) {
	user := model.User{}
	db.GormDB.Model(&model.User{}).Where("mobile = ?", mobile).Find(&user)
	if user.ID > 0 {
		return model.User{}, "用户已存在"
	}
	md5Password := d.GetMd5Password(password)
	user = model.User{
		Mobile:   mobile,
		Password: md5Password,
		Nickname: mobile,
		Status:   1,
	}
	db.GormDB.Model(&model.User{}).Create(&user)
	return user, ""
}
