package seckill

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSeckillEntriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSeckillEntriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSeckillEntriesLogic {
	return &AdminListSeckillEntriesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSeckillEntriesLogic) AdminListSeckillEntries(ctx context.Context, req *types.SeckillEntryListReq) (resp *types.PageListResp, err error) {
	p, ps := req.Page, req.PageSize
	sid := req.SessionId
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
