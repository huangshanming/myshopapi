package repository

import (
	"errors"

	"mymall/common/password"
	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) HashPassword(plain string) string {
	return password.Hash(plain)
}

func (r *UserRepository) UpdatePassword(id uint64, plain string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", password.Hash(plain)).Error
}

func (r *UserRepository) CreateAdmin(mobile, plain, nickname string) (*model.User, error) {
	var existing model.User
	if err := r.db.Where("mobile = ?", mobile).First(&existing).Error; err == nil {
		return nil, errors.New("用户已存在")
	}
	if nickname == "" {
		nickname = mobile
	}
	user := model.User{
		Mobile:   mobile,
		Password: password.Hash(plain),
		Nickname: nickname,
		Status:   1,
		Role:     "platform_admin",
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	var user model.User
	err := r.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) VerifyLogin(mobile, plain string) (*model.User, error) {
	user, err := r.FindByMobile(mobile)
	if err != nil {
		return nil, err
	}
	if user.Password != password.Hash(plain) {
		return nil, gorm.ErrRecordNotFound
	}
	if user.Status != 1 {
		return nil, errors.New("账号已禁用")
	}
	return user, nil
}

func (r *UserRepository) Create(mobile, plain string) (*model.User, error) {
	var existing model.User
	if err := r.db.Where("mobile = ?", mobile).First(&existing).Error; err == nil {
		return nil, errors.New("用户已存在")
	}

	user := model.User{
		Mobile:   mobile,
		Password: password.Hash(plain),
		Nickname: mobile,
		Status:   1,
		Role:     "user",
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProfile(id uint64, nickname, avatar, mobile string, gender int) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"avatar":   avatar,
		"gender":   gender,
		"mobile":   mobile,
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) MobileTakenByOther(mobile string, excludeID uint64) bool {
	var count int64
	r.db.Model(&model.User{}).Where("mobile = ? AND id <> ?", mobile, excludeID).Count(&count)
	return count > 0
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
