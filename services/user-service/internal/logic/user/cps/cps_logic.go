package cps

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListActsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListActsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActsLogic {
	return &ListActsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListActsLogic) ListActs(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	list, total, err := biz.NewCpsLogic(l.svcCtx).ListActs(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	if list == nil {
		list = []biz.CpsActVO{}
	}
	return &types.PageListResp{Total: total, List: list}, nil
}

type ConvertLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(ctx context.Context, req *types.CpsConvertReq) (*types.CpsConvertResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	vo, err := biz.NewCpsLogic(l.svcCtx).Convert(ctx, userID, req.ActID)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "未登录" {
			code = http.StatusUnauthorized
		}
		return nil, xerr.New(code, err.Error())
	}
	return &types.CpsConvertResp{Data: vo}, nil
}

type ListGoodsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGoodsLogic {
	return &ListGoodsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListGoodsLogic) ListGoods(ctx context.Context, req *types.CpsGoodsListReq) (*types.PageListResp, error) {
	list, total, err := biz.NewCpsLogic(l.svcCtx).ListGoods(ctx, req.Platform, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	if list == nil {
		list = []biz.CpsGoodsVO{}
	}
	return &types.PageListResp{Total: total, List: list}, nil
}

type ConvertGoodsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewConvertGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertGoodsLogic {
	return &ConvertGoodsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ConvertGoodsLogic) ConvertGoods(ctx context.Context, req *types.CpsGoodsConvertReq) (*types.CpsConvertResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	vo, err := biz.NewCpsLogic(l.svcCtx).ConvertGoods(ctx, userID, req.Platform, req.ItemID, req.RawURL)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "未登录" {
			code = http.StatusUnauthorized
		}
		return nil, xerr.New(code, err.Error())
	}
	return &types.CpsConvertResp{Data: vo}, nil
}
