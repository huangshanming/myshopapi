package region

import (
	"context"
	"mymall/pkg/httpinvoke"
	hpublic "mymall/services/user-service/internal/app/public"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegionTreeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRegionTreeLogic(svcCtx *svc.ServiceContext) *RegionTreeLogic {
	return &RegionTreeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *RegionTreeLogic) RegionTree(ctx context.Context) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/regions/tree", nil, nil, nil, hpublic.NewRegionHandler(l.svcCtx).Tree)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
