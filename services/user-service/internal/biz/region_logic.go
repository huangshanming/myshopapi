package biz

import (
	"context"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type RegionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegionLogic {
	return &RegionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegionLogic) ListChildren(parentCode string) ([]model.Region, error) {
	return l.svcCtx.Repo.ListRegionsByParent(parentCode)
}

func (l *RegionLogic) Tree() ([]model.RegionTreeNode, error) {
	return l.svcCtx.Repo.BuildRegionTree()
}
