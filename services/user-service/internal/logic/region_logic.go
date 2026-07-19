package logic

import (
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type RegionLogic struct {
	svcCtx *svc.ServiceContext
}

func NewRegionLogic(svcCtx *svc.ServiceContext) *RegionLogic {
	return &RegionLogic{svcCtx: svcCtx}
}

func (l *RegionLogic) ListChildren(parentCode string) ([]model.Region, error) {
	return l.svcCtx.Repo.ListRegionsByParent(parentCode)
}

func (l *RegionLogic) Tree() ([]model.RegionTreeNode, error) {
	return l.svcCtx.Repo.BuildRegionTree()
}
