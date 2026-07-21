package region

import (
	"context"
	"mymall/pkg/appinput"
	hpublic "mymall/services/user-service/internal/app/public"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegionTreeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRegionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegionTreeLogic {
	return &RegionTreeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *RegionTreeLogic) RegionTree(ctx context.Context) (resp *types.AnyResp, err error) {
	data, err := hpublic.NewRegionHandler(l.svcCtx).Tree(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
