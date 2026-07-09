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
	user, err := s.repo.VerifyLogin(mobile, password)
	if err != nil {
		return "", nil, errors.New("手机号或密码错误")
	}
	token, err := jwt.GenerateToken(user.ID, jwt.RoleUser, s.jwtCfg)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *UserService) Register(mobile, password string) (*model.User, error) {
	return s.repo.Create(mobile, password)
}

func (s *UserService) GetProfile(userID uint64) (*model.User, error) {
	return s.repo.FindByID(userID)
}
