package address

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type InternalGetLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalGetLogic {
	return &InternalGetLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalGetLogic) InternalGet(ctx context.Context, req *types.InternalAddressReq) (resp *types.AnyResp, err error) {
	if req.UserID == 0 || req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数无效")
	}
	a, err := biz.NewAddressLogic(l.svcCtx).Get(ctx, req.UserID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: a}, nil
}
