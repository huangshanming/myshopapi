package seckill

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicSeckillListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillListLogic {
	return &PublicSeckillListLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillListLogic) PublicSeckillList(ctx context.Context, req *types.PageReq) (resp *types.SeckillListResp, err error) {
	p, ps := req.Page, req.PageSize
	data, err := biz.NewMerchantLogic(l.svcCtx).PublicSeckillList(p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	// Flat body: {session_id, start_at, end_at, total, list} for mall-uni
	return &types.SeckillListResp{Data: data}, nil
}
