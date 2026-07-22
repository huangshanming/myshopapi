package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"mymall/services/user-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const regionColumns = "id, IFNULL(code,'') AS code, IFNULL(name,'') AS name, IFNULL(parent_code,'') AS parent_code, level, sort, created_at"

type pcaNode struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Children []pcaNode `json:"children"`
}

func (r *UserRepository) CountRegions(ctx context.Context) (int64, error) {
	return countQuery(ctx, r.conn, "SELECT COUNT(*) FROM regions")
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
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		const batch = 500
		for i := 0; i < len(rows); i += batch {
			end := i + batch
			if end > len(rows) {
				end = len(rows)
			}
			chunk := rows[i:end]
			placeholders := make([]string, 0, len(chunk))
			args := make([]any, 0, len(chunk)*5)
			for _, reg := range chunk {
				placeholders = append(placeholders, "(?,?,?,?,?)")
				args = append(args, reg.Code, reg.Name, reg.ParentCode, reg.Level, reg.Sort)
			}
			query := fmt.Sprintf(
				"INSERT INTO regions (code, name, parent_code, level, sort) VALUES %s "+
					"ON DUPLICATE KEY UPDATE name=VALUES(name), parent_code=VALUES(parent_code), level=VALUES(level), sort=VALUES(sort)",
				strings.Join(placeholders, ","),
			)
			if _, err := session.ExecCtx(ctx, query, args...); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepository) ListRegionsByParent(ctx context.Context, parentCode string) ([]model.Region, error) {
	var list []model.Region
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+regionColumns+" FROM regions WHERE parent_code=? ORDER BY sort ASC, code ASC",
		parentCode,
	)
	return list, err
}

func (r *UserRepository) ListRegionsByLevel(ctx context.Context, level int) ([]model.Region, error) {
	var list []model.Region
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+regionColumns+" FROM regions WHERE level=? ORDER BY sort ASC, code ASC",
		level,
	)
	return list, err
}

func (r *UserRepository) GetRegionByCode(ctx context.Context, code string) (*model.Region, error) {
	var reg model.Region
	err := r.conn.QueryRowPartialCtx(ctx, &reg,
		"SELECT "+regionColumns+" FROM regions WHERE code=? LIMIT 1", code,
	)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *UserRepository) BuildRegionTree(ctx context.Context) ([]model.RegionTreeNode, error) {
	var all []model.Region
	if err := r.conn.QueryRowsPartialCtx(ctx, &all,
		"SELECT "+regionColumns+" FROM regions ORDER BY level ASC, sort ASC, code ASC",
	); err != nil {
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
