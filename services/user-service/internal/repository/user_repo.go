package repository

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
)

const passwordSalt = "this is my mall"

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) hashPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password + passwordSalt))
	return hex.EncodeToString(hash.Sum(nil))
}

func (r *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	var user model.User
	err := r.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) VerifyLogin(mobile, password string) (*model.User, error) {
	user, err := r.FindByMobile(mobile)
	if err != nil {
		return nil, err
	}
	if user.Password != r.hashPassword(password) {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (r *UserRepository) Create(mobile, password string) (*model.User, error) {
	var existing model.User
	if err := r.db.Where("mobile = ?", mobile).First(&existing).Error; err == nil {
		return nil, errors.New("用户已存在")
	}

	user := model.User{
		Mobile:   mobile,
		Password: r.hashPassword(password),
		Nickname: mobile,
		Status:   1,
		Role:     "user",
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FirstShopID 取用户所属第一家店铺（同库 shop_members）
func (r *UserRepository) FirstShopID(userID uint64) uint64 {
	var shopID uint64
	_ = r.db.Table("shop_members").Select("shop_id").Where("user_id = ?", userID).
		Order("id ASC").Limit(1).Scan(&shopID).Error
	return shopID
}
