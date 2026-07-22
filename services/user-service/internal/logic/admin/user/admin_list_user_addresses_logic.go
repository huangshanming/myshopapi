package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AdminListUserAddressesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListUserAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserAddressesLogic {
	return &AdminListUserAddressesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListUserAddressesLogic) AdminListUserAddresses(ctx context.Context, req *types.IdPathReq) (resp *types.PageListResp, err error) {
	list, err := biz.NewAddressLogic(l.svcCtx).List(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: int64(len(list))}, nil
}
