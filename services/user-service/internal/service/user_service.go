package service

import (
	"errors"

	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/repository"
)

type UserService struct {
	repo   *repository.UserRepository
	jwtCfg jwt.Config
}

func NewUserService(repo *repository.UserRepository, jwtCfg jwt.Config) *UserService {
	return &UserService{repo: repo, jwtCfg: jwtCfg}
}

func (s *UserService) Login(mobile, password string) (string, *model.User, error) {
	return s.LoginWithShop(mobile, password, 0)
}

// LoginWithShop 登录；shopID>0 时写入 JWT（商家切换店铺）
func (s *UserService) LoginWithShop(mobile, password string, shopID uint64) (string, *model.User, error) {
	user, err := s.repo.VerifyLogin(mobile, password)
	if err != nil {
		return "", nil, errors.New("手机号或密码错误")
	}
	role := user.Role
	if role == "" {
		role = jwt.RoleUser
	}
	if shopID == 0 && jwt.IsMerchant(role) {
		shopID = s.repo.FirstShopID(user.ID)
	}
	token, err := jwt.GenerateTokenWithShop(user.ID, role, shopID, s.jwtCfg)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// SwitchShopToken 已登录用户切换店铺（重新签发带 shop_id 的 token）
func (s *UserService) SwitchShopToken(userID uint64, role string, shopID uint64) (string, error) {
	return jwt.GenerateTokenWithShop(userID, role, shopID, s.jwtCfg)
}

func (s *UserService) Register(mobile, password string) (*model.User, error) {
	return s.repo.Create(mobile, password)
}

func (s *UserService) GetProfile(userID uint64) (*model.User, error) {
	return s.repo.FindByID(userID)
}
