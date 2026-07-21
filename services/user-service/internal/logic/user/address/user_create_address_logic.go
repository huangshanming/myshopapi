package address

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserCreateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateAddressLogic {
	return &UserCreateAddressLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserCreateAddressLogic) UserCreateAddress(ctx context.Context, req *types.AddressReq) (*types.AnyResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	a, err := biz.NewAddressLogic(l.svcCtx).Create(ctx, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: a}, nil
}
