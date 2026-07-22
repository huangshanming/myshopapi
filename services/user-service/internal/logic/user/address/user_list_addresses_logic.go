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

type UserListAddressesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListAddressesLogic {
	return &UserListAddressesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserListAddressesLogic) UserListAddresses(ctx context.Context) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	list, err := biz.NewAddressLogic(l.svcCtx).List(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: int64(len(list))}, nil
}
