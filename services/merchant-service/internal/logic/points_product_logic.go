package logic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/uploadpath"
)

type PointsProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointsProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointsProductLogic {
	return &PointsProductLogic{ctx: ctx, svcCtx: svcCtx}
}

type PointsProductSaveReq struct {
	Name         string `json:"name"`
	CoverURL     string `json:"cover_url"`
	Description  string `json:"description"`
	PointsPrice  *int   `json:"points_price"`
	Stock        *int   `json:"stock"`
	PerUserLimit *int   `json:"per_user_limit"`
	Status       string `json:"status"`
	Sort         *int   `json:"sort"`
}

func (l *PointsProductLogic) List(page, pageSize int, status, keyword string) ([]model.PointsProduct, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.PointsProducts.List(page, pageSize, status, keyword)
}

func (l *PointsProductLogic) Get(id uint64) (*model.PointsProduct, error) {
	p, err := l.svcCtx.PointsProducts.GetByID(id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	return p, nil
}

func (l *PointsProductLogic) Create(req PointsProductSaveReq) (*model.PointsProduct, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("请填写商品名称")
	}
	p := &model.PointsProduct{
		Name: name, CoverURL: strings.TrimSpace(req.CoverURL), Description: strings.TrimSpace(req.Description),
		Status: model.PointsProductStatusOff,
	}
	if req.PointsPrice != nil {
		if *req.PointsPrice < 0 {
			return nil, errors.New("积分价不能为负")
		}
		p.PointsPrice = *req.PointsPrice
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("库存不能为负")
		}
		p.Stock = *req.Stock
	}
	if req.PerUserLimit != nil {
		if *req.PerUserLimit < 0 {
			return nil, errors.New("限兑次数无效")
		}
		p.PerUserLimit = *req.PerUserLimit
	}
	if req.Sort != nil {
		p.Sort = *req.Sort
	}
	if s := strings.TrimSpace(req.Status); s != "" {
		if s != model.PointsProductStatusOn && s != model.PointsProductStatusOff {
			return nil, errors.New("状态无效")
		}
		p.Status = s
	}
	if err := l.svcCtx.PointsProducts.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (l *PointsProductLogic) Update(id uint64, req PointsProductSaveReq) (*model.PointsProduct, error) {
	if _, err := l.svcCtx.PointsProducts.GetByID(id); err != nil {
		return nil, errors.New("商品不存在")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("请填写商品名称")
	}
	updates := map[string]interface{}{
		"name": name, "cover_url": strings.TrimSpace(req.CoverURL), "description": strings.TrimSpace(req.Description),
	}
	if req.PointsPrice != nil {
		if *req.PointsPrice < 0 {
			return nil, errors.New("积分价不能为负")
		}
		updates["points_price"] = *req.PointsPrice
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("库存不能为负")
		}
		updates["stock"] = *req.Stock
	}
	if req.PerUserLimit != nil {
		if *req.PerUserLimit < 0 {
			return nil, errors.New("限兑次数无效")
		}
		updates["per_user_limit"] = *req.PerUserLimit
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if s := strings.TrimSpace(req.Status); s != "" {
		if s != model.PointsProductStatusOn && s != model.PointsProductStatusOff {
			return nil, errors.New("状态无效")
		}
		updates["status"] = s
	}
	if err := l.svcCtx.PointsProducts.Update(id, updates); err != nil {
		return nil, err
	}
	return l.svcCtx.PointsProducts.GetByID(id)
}

func (l *PointsProductLogic) SetStatus(id uint64, status string) (*model.PointsProduct, error) {
	status = strings.TrimSpace(status)
	if status != model.PointsProductStatusOn && status != model.PointsProductStatusOff {
		return nil, errors.New("状态无效")
	}
	if _, err := l.svcCtx.PointsProducts.GetByID(id); err != nil {
		return nil, errors.New("商品不存在")
	}
	if err := l.svcCtx.PointsProducts.Update(id, map[string]interface{}{"status": status}); err != nil {
		return nil, err
	}
	return l.svcCtx.PointsProducts.GetByID(id)
}

func (l *PointsProductLogic) Delete(id uint64) error {
	if _, err := l.svcCtx.PointsProducts.GetByID(id); err != nil {
		return errors.New("商品不存在")
	}
	return l.svcCtx.PointsProducts.Delete(id)
}

func (l *PointsProductLogic) SaveUpload(filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	dir := uploadpath.Abs("points-mall")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/points-mall/" + name, nil
}
