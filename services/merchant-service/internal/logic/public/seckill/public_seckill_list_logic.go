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

func (l *PublicSeckillListLogic) PublicSeckillList(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	p, ps := in.Page()
	data, err := biz.NewMerchantLogic(l.svcCtx).PublicSeckillList(p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
