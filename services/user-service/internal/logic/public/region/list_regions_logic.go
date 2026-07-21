package region

import (
	"context"
	"encoding/json"
	"mymall/pkg/appinput"
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

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(ctx context.Context, req *types.RegionListReq) (resp *types.PageListResp, err error) {
	data, err := hpublic.NewRegionHandler(l.svcCtx).List(ctx, appinput.CallInput{Query: url.Values{"parent_code": {req.ParentCode}}})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.PageListResp
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
