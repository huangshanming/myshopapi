package region

import (
	"context"
	"mymall/pkg/httpinvoke"
	hpublic "mymall/services/user-service/internal/app/public"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(ctx context.Context, req *types.RegionListReq) (resp *types.PageListResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/regions", nil, url.Values{"parent_code": {req.ParentCode}}, nil, hpublic.NewRegionHandler(l.svcCtx).List)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
