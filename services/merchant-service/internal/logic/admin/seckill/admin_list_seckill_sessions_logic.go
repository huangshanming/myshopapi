package seckill

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSeckillSessionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSeckillSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSeckillSessionsLogic {
	return &AdminListSeckillSessionsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSeckillSessionsLogic) AdminListSeckillSessions(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	p, ps := in.Page()
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListSeckillSessions(p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
