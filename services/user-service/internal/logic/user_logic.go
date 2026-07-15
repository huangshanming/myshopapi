package logic

import (
	"errors"

	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type UserLogic struct {
	svcCtx *svc.ServiceContext
}

func NewUserLogic(svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{svcCtx: svcCtx}
}

func (l *UserLogic) Login(mobile, password string) (string, *model.User, error) {
	return l.LoginWithShop(mobile, password, 0)
}

// LoginWithShop 登录；shopID>0 时写入 JWT（商家切换店铺）
func (l *UserLogic) LoginWithShop(mobile, password string, shopID uint64) (string, *model.User, error) {
	user, err := l.svcCtx.Repo.VerifyLogin(mobile, password)
	if err != nil {
		return "", nil, errors.New("手机号或密码错误")
	}
	role := user.Role
	if role == "" {
		role = jwt.RoleUser
	}
	if shopID == 0 && jwt.IsMerchant(role) {
		shopID = l.svcCtx.Repo.FirstShopID(user.ID)
	}
	token, err := jwt.GenerateTokenWithShop(user.ID, role, shopID, l.svcCtx.JWT)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// SwitchShopToken 已登录用户切换店铺（重新签发带 shop_id 的 token）
func (l *UserLogic) SwitchShopToken(userID uint64, role string, shopID uint64) (string, error) {
	return jwt.GenerateTokenWithShop(userID, role, shopID, l.svcCtx.JWT)
}

func (l *UserLogic) Register(mobile, password string) (*model.User, error) {
	return l.svcCtx.Repo.Create(mobile, password)
}

func (l *UserLogic) GetProfile(userID uint64) (*model.User, error) {
	return l.svcCtx.Repo.FindByID(userID)
}
