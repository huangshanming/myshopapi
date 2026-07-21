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

type UpdatePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePointsProductLogic {
	return &UpdatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePointsProductLogic) UpdatePointsProduct(ctx context.Context, req *types.PointsProductUpdateReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	p, err := biz.NewPointsProductLogic(l.svcCtx).Update(ctx, req.Id, toBizPointsProductUpdate(req))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: p}, nil
}
