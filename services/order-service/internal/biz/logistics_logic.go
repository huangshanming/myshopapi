package biz

import (
	"context"
	"errors"
	"strings"

	"mymall/services/order-service/internal/model"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"gorm.io/gorm"
)

type LogisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogisticsLogic {
	return &LogisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogisticsLogic) List(f repository.LogisticsListFilter) ([]model.LogisticsCompany, int64, error) {
	return l.svcCtx.LogisticsRepo.List(l.ctx, f)
}

func (l *LogisticsLogic) Options(keyword string) ([]model.LogisticsCompany, error) {
	return l.svcCtx.LogisticsRepo.Options(l.ctx, keyword, 50)
}

func (l *LogisticsLogic) Create(req types.LogisticsSaveReq) (*model.LogisticsCompany, error) {
	name := strings.TrimSpace(req.Name)
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if name == "" || code == "" {
		return nil, errors.New("名称与编码必填")
	}
	c := &model.LogisticsCompany{
		Name:   name,
		Code:   code,
		Sort:   req.Sort,
		Status: 1,
	}
	if req.Status != nil {
		c.Status = *req.Status
	}
	if err := l.svcCtx.LogisticsRepo.Create(l.ctx, c); err != nil {
		if isDuplicate(err) {
			return nil, errors.New("编码已存在")
		}
		return nil, err
	}
	return c, nil
}

func (l *LogisticsLogic) Update(id uint64, req types.LogisticsSaveReq) error {
	name := strings.TrimSpace(req.Name)
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if name == "" || code == "" {
		return errors.New("名称与编码必填")
	}
	if err := l.svcCtx.LogisticsRepo.Update(l.ctx, id, name, code, req.Sort); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("记录不存在")
		}
		if isDuplicate(err) {
			return errors.New("编码已存在")
		}
		return err
	}
	return nil
}

func (l *LogisticsLogic) UpdateStatus(id uint64, status int8) error {
	if status != 0 && status != 1 {
		return errors.New("状态无效")
	}
	if err := l.svcCtx.LogisticsRepo.UpdateStatus(l.ctx, id, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("记录不存在")
		}
		return err
	}
	return nil
}

func (l *LogisticsLogic) Delete(id uint64) error {
	if err := l.svcCtx.LogisticsRepo.Delete(l.ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("记录不存在")
		}
		return err
	}
	return nil
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
