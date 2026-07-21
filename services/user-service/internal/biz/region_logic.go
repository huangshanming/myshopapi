package biz

import (
	"context"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type RegionLogic struct {
	svcCtx *svc.ServiceContext
}

func NewRegionLogic(svcCtx *svc.ServiceContext) *RegionLogic {
	return &RegionLogic{svcCtx: svcCtx}
}

func (l *RegionLogic) ListChildren(ctx context.Context, parentCode string) ([]model.Region, error) {
	return l.svcCtx.Repo.ListRegionsByParent(ctx, parentCode)
}

func (l *RegionLogic) Tree(ctx context.Context) ([]model.RegionTreeNode, error) {
	return l.svcCtx.Repo.BuildRegionTree(ctx)
}
