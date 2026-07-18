package repository

import (
	"mymall/services/order-service/internal/model"

	"gorm.io/gorm"
)

type LogisticsListFilter struct {
	Name     string
	Code     string
	Status   *int8
	Keyword  string
	Page     int
	PageSize int
	EnabledOnly bool
}

type LogisticsRepository struct {
	db *gorm.DB
}

func NewLogisticsRepository(db *gorm.DB) *LogisticsRepository {
	return &LogisticsRepository{db: db}
}

func (r *LogisticsRepository) List(f LogisticsListFilter) ([]model.LogisticsCompany, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	q := r.db.Model(&model.LogisticsCompany{})
	if f.EnabledOnly {
		q = q.Where("status = 1")
	} else if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Name != "" {
		q = q.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		q = q.Where("code LIKE ?", "%"+f.Code+"%")
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.LogisticsCompany
	err := q.Order("sort ASC, id ASC").
		Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&list).Error
	return list, total, err
}

func (r *LogisticsRepository) Options(keyword string, limit int) ([]model.LogisticsCompany, error) {
	if limit < 1 {
		limit = 50
	}
	q := r.db.Model(&model.LogisticsCompany{}).Where("status = 1")
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	var list []model.LogisticsCompany
	err := q.Order("sort ASC, id ASC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *LogisticsRepository) Create(c *model.LogisticsCompany) error {
	return r.db.Create(c).Error
}

func (r *LogisticsRepository) Update(id uint64, name, code string, sort int) error {
	res := r.db.Model(&model.LogisticsCompany{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": name,
		"code": code,
		"sort": sort,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *LogisticsRepository) UpdateStatus(id uint64, status int8) error {
	res := r.db.Model(&model.LogisticsCompany{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *LogisticsRepository) Delete(id uint64) error {
	res := r.db.Delete(&model.LogisticsCompany{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *LogisticsRepository) FindByCode(code string) (*model.LogisticsCompany, error) {
	var c model.LogisticsCompany
	err := r.db.Where("code = ?", code).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *LogisticsRepository) SeedDefaults() error {
	seeds := []model.LogisticsCompany{
		{Name: "顺丰速运", Code: "SF", Sort: 10, Status: 1},
		{Name: "中通快递", Code: "ZTO", Sort: 20, Status: 1},
		{Name: "圆通速递", Code: "YTO", Sort: 30, Status: 1},
		{Name: "韵达快递", Code: "YD", Sort: 40, Status: 1},
		{Name: "申通快递", Code: "STO", Sort: 50, Status: 1},
		{Name: "EMS", Code: "EMS", Sort: 60, Status: 1},
		{Name: "京东物流", Code: "JD", Sort: 70, Status: 1},
		{Name: "德邦快递", Code: "DBL", Sort: 80, Status: 1},
	}
	for _, s := range seeds {
		var n int64
		r.db.Model(&model.LogisticsCompany{}).Where("code = ?", s.Code).Count(&n)
		if n == 0 {
			if err := r.db.Create(&s).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
