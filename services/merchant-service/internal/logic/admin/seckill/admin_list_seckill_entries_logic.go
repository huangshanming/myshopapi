package seckill

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"net/url"
	"strconv"

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

func (l *AdminListSeckillEntriesLogic) AdminListSeckillEntries(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	p, ps := in.Page()
	sid, _ := strconv.ParseUint(in.QueryGet("session_id"), 10, 64)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
