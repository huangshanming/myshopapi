package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"mymall/services/user-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pcaNode struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Children []pcaNode `json:"children"`
}

func (r *UserRepository) CountRegions(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Region{}).Count(&n).Error
	return n, err
}

func (r *UserRepository) SeedRegionsFromPCA(ctx context.Context, raw []byte) error {
	var roots []pcaNode
	if err := json.Unmarshal(raw, &roots); err != nil {
		return fmt.Errorf("parse pca json: %w", err)
	}
	rows := make([]model.Region, 0, 4096)
	var walk func(nodes []pcaNode, parent string, level int)
	walk = func(nodes []pcaNode, parent string, level int) {
		for i, n := range nodes {
			if n.Code == "" || n.Name == "" {
				continue
			}
			rows = append(rows, model.Region{
				Code:       n.Code,
				Name:       n.Name,
				ParentCode: parent,
				Level:      level,
				Sort:       i + 1,
			})
			if len(n.Children) > 0 {
				walk(n.Children, n.Code, level+1)
			}
		}
	}
	walk(roots, "", model.RegionLevelProvince)
	if len(rows) == 0 {
		return fmt.Errorf("pca data empty")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		const batch = 500
		for i := 0; i < len(rows); i += batch {
			end := i + batch
			if end > len(rows) {
				end = len(rows)
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "code"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "parent_code", "level", "sort"}),
			}).Create(rows[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepository) ListRegionsByParent(ctx context.Context, parentCode string) ([]model.Region, error) {
	var list []model.Region
	err := r.db.WithContext(ctx).Where("parent_code = ?", parentCode).Order("sort ASC, code ASC").Find(&list).Error
	return list, err
}

func (r *UserRepository) ListRegionsByLevel(ctx context.Context, level int) ([]model.Region, error) {
	var list []model.Region
	err := r.db.WithContext(ctx).Where("level = ?", level).Order("sort ASC, code ASC").Find(&list).Error
	return list, err
}

func (r *UserRepository) GetRegionByCode(ctx context.Context, code string) (*model.Region, error) {
	var reg model.Region
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&reg).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *UserRepository) BuildRegionTree(ctx context.Context) ([]model.RegionTreeNode, error) {
	var all []model.Region
	if err := r.db.WithContext(ctx).Order("level ASC, sort ASC, code ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	byParent := map[string][]model.Region{}
	for _, reg := range all {
		byParent[reg.ParentCode] = append(byParent[reg.ParentCode], reg)
	}
	var build func(parent string) []model.RegionTreeNode
	build = func(parent string) []model.RegionTreeNode {
		children := byParent[parent]
		out := make([]model.RegionTreeNode, 0, len(children))
		for _, c := range children {
			node := model.RegionTreeNode{Code: c.Code, Name: c.Name, Level: c.Level}
			if kids := build(c.Code); len(kids) > 0 {
				node.Children = kids
			}
			out = append(out, node)
		}
		return out
	}
	return build(""), nil
}
