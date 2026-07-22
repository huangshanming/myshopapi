package region

import (
	"context"
	"net/http"
	"github.com/zeromicro/go-zero/core/logx"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type RegionTreeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRegionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegionTreeLogic {
	return &RegionTreeLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *RegionTreeLogic) RegionTree(ctx context.Context) (*types.TreeResp, error) {
	tree, err := biz.NewRegionLogic(l.svcCtx).Tree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.TreeResp{Tree: tree}, nil
}
