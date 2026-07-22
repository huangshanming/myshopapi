package shop

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminResetOwnerPasswordLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminResetOwnerPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminResetOwnerPasswordLogic {
	return &AdminResetOwnerPasswordLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminResetOwnerPasswordLogic) AdminResetOwnerPassword(ctx context.Context, req *types.OwnerPasswordBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).ResetOwnerPassword(ctx, id, req.Password); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
