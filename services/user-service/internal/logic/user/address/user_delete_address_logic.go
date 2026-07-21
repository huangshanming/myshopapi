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

type UserDeleteAddressLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteAddressLogic {
	return &UserDeleteAddressLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserDeleteAddressLogic) UserDeleteAddress(ctx context.Context, req *types.IdPathReq) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := biz.NewAddressLogic(l.svcCtx).Delete(ctx, userID, req.Id); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
