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

func toBizPointsProductSave(req *types.PointsProductSaveReq) biz.PointsProductSaveReq {
	out := biz.PointsProductSaveReq{
		Name:        req.Name,
		CoverURL:    req.CoverURL,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.PointsPrice != 0 {
		v := req.PointsPrice
		out.PointsPrice = &v
	}
	if req.Stock != 0 {
		v := req.Stock
		out.Stock = &v
	}
	if req.PerUserLimit != 0 {
		v := req.PerUserLimit
		out.PerUserLimit = &v
	}
	if req.Sort != 0 {
		v := req.Sort
		out.Sort = &v
	}
	return out
}

func toBizPointsProductUpdate(req *types.PointsProductUpdateReq) biz.PointsProductSaveReq {
	out := biz.PointsProductSaveReq{
		Name:        req.Name,
		CoverURL:    req.CoverURL,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.PointsPrice != 0 {
		v := req.PointsPrice
		out.PointsPrice = &v
	}
	if req.Stock != 0 {
		v := req.Stock
		out.Stock = &v
	}
	if req.PerUserLimit != 0 {
		v := req.PerUserLimit
		out.PerUserLimit = &v
	}
	if req.Sort != 0 {
		v := req.Sort
		out.Sort = &v
	}
	return out
}

type CreatePointsProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreatePointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePointsProductLogic {
	return &CreatePointsProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreatePointsProductLogic) CreatePointsProduct(ctx context.Context, req *types.PointsProductSaveReq) (resp *types.AnyResp, err error) {
	p, err := biz.NewPointsProductLogic(l.svcCtx).Create(ctx, toBizPointsProductSave(req))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: p}, nil
}
