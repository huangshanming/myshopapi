package repository

import (
	"errors"
	"time"

	"mymall/services/catalog-service/internal/content/model"
)

func (r *ArticleRepository) ListBannersAdmin(page, pageSize int) ([]model.HomepageBanner, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.Model(&model.HomepageBanner{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.HomepageBanner
	err := q.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) ListBannersPublic() ([]model.HomepageBanner, error) {
	now := time.Now()
	var list []model.HomepageBanner
	err := r.db.Model(&model.HomepageBanner{}).
		Where("status = ?", model.BannerOn).
		Where("(start_at IS NULL OR start_at <= ?)", now).
		Where("(end_at IS NULL OR end_at > ?)", now).
		Order("sort ASC, id DESC").
		Find(&list).Error
	return list, err
}

func (r *ArticleRepository) GetBanner(id uint64) (*model.HomepageBanner, error) {
	var b model.HomepageBanner
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *ArticleRepository) CreateBanner(b *model.HomepageBanner) error {
	return r.db.Create(b).Error
}

func (r *ArticleRepository) UpdateBanner(id uint64, updates map[string]interface{}) error {
	res := r.db.Model(&model.HomepageBanner{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("Banner 不存在")
	}
	return nil
}

func (r *ArticleRepository) DeleteBanner(id uint64) error {
	res := r.db.Delete(&model.HomepageBanner{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("Banner 不存在")
	}
	return nil
}

func (r *ArticleRepository) FillBannerLinkNames(list []model.HomepageBanner) {
	for i := range list {
		switch list[i].LinkType {
		case model.BannerLinkProduct:
			if list[i].LinkID == 0 {
				continue
			}
			var name string
			_ = r.db.Table("products").Select("name").Where("id = ?", list[i].LinkID).Scan(&name).Error
			list[i].LinkName = name
		case model.BannerLinkArticle:
			if list[i].LinkID == 0 {
				continue
			}
			var title string
			_ = r.db.Table("community_article").Select("title").Where("id = ?", list[i].LinkID).Scan(&title).Error
			list[i].LinkName = title
		}
	}
}

func (r *ArticleRepository) ProductExistsOnSale(id uint64) bool {
	var n int64
	r.db.Table("products").Where("id = ? AND status = ?", id, "on_sale").Count(&n)
	return n > 0
}

func (r *ArticleRepository) ArticleExistsPublished(id uint64) bool {
	var n int64
	r.db.Model(&model.CommunityArticle{}).
		Where("id = ? AND status = ? AND audit_status = ?", id, model.ArticlePublished, model.ArticleAuditApproved).
		Count(&n)
	return n > 0
}
