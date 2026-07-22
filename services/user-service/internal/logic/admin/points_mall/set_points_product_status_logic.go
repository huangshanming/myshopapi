package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type SetPointsProductStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewSetPointsProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPointsProductStatusLogic {
	return &SetPointsProductStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *SetPointsProductStatusLogic) SetPointsProductStatus(ctx context.Context, req *types.PointsProductStatusReq) (resp *types.PointsProductResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	p, err := biz.NewPointsProductLogic(l.svcCtx).SetStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PointsProductResp{Data: p}, nil
}
