package repository

import (
	"errors"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/model"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

type ArticleListFilter struct {
	ShopID       uint64
	Title        string
	AuditStatus  string
	Status       string
	HasSchedule  *bool // nil=不限, true=有定时, false=无定时
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	Recycle      bool
	Page         int
	PageSize     int
}

func (r *ArticleRepository) List(f ArticleListFilter) ([]model.CommunityArticle, int64, error) {
	q := r.db.Model(&model.CommunityArticle{})
	if f.Recycle {
		q = q.Where("status = ?", model.ArticleDeleted)
	} else {
		q = q.Where("status <> ?", model.ArticleDeleted)
	}
	if f.ShopID > 0 {
		q = q.Where("shop_id = ?", f.ShopID)
	}
	if f.Title != "" {
		q = q.Where("title LIKE ?", "%"+f.Title+"%")
	}
	if f.AuditStatus != "" {
		q = q.Where("audit_status = ?", f.AuditStatus)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.HasSchedule != nil {
		if *f.HasSchedule {
			q = q.Where("schedule_publish_at IS NOT NULL")
		} else {
			q = q.Where("schedule_publish_at IS NULL")
		}
	}
	if f.CreatedFrom != nil {
		q = q.Where("created_at >= ?", *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		q = q.Where("created_at <= ?", *f.CreatedTo)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.CommunityArticle
	err := q.Order("is_top DESC, id DESC").
		Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).
		Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) GetByID(id uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) GetByIDShop(id, shopID uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.Where("id = ? AND shop_id = ?", id, shopID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) Create(a *model.CommunityArticle) error {
	return r.db.Create(a).Error
}

func (r *ArticleRepository) Update(id uint64, updates map[string]interface{}) error {
	res := r.db.Model(&model.CommunityArticle{}).Where("id = ?", id).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return res.Error
}

func (r *ArticleRepository) UpdateShop(id, shopID uint64, updates map[string]interface{}) error {
	res := r.db.Model(&model.CommunityArticle{}).Where("id = ? AND shop_id = ?", id, shopID).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("文章不存在或无权操作")
	}
	return res.Error
}

func (r *ArticleRepository) SoftDelete(id uint64) error {
	now := common.LocalTime(time.Now())
	return r.Update(id, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) SoftDeleteShop(id, shopID uint64) error {
	now := common.LocalTime(time.Now())
	return r.UpdateShop(id, shopID, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) Restore(id uint64) error {
	return r.Update(id, map[string]interface{}{
		"status": model.ArticleOffline, "deleted_at": nil,
	})
}

func (r *ArticleRepository) PermanentDelete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", id).Delete(&model.CommunityArticleComment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("article_id = ?", id).Delete(&model.CommunityArticleImg{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&model.CommunityArticle{}).Error
	})
}

func (r *ArticleRepository) ReplaceImages(articleID, shopID uint64, urls []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&model.CommunityArticleImg{}).Error; err != nil {
			return err
		}
		for i, u := range urls {
			if u == "" {
				continue
			}
			img := model.CommunityArticleImg{ArticleID: articleID, ShopID: shopID, URL: u, Sort: i}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ArticleRepository) ListImages(articleID uint64) ([]model.CommunityArticleImg, error) {
	var list []model.CommunityArticleImg
	err := r.db.Where("article_id = ?", articleID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ArticleRepository) ClaimDuePublish(limit int) ([]model.CommunityArticle, error) {
	var list []model.CommunityArticle
	now := time.Now()
	err := r.db.Where("status = ? AND audit_status = ? AND schedule_publish_at IS NOT NULL AND schedule_publish_at <= ?",
		model.ArticleScheduled, model.ArticleAuditApproved, now).
		Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make([]model.CommunityArticle, 0, len(list))
	pub := common.LocalTime(now)
	for _, a := range list {
		res := r.db.Model(&model.CommunityArticle{}).
			Where("id = ? AND status = ?", a.ID, model.ArticleScheduled).
			Updates(map[string]interface{}{
				"status": model.ArticlePublished, "published_at": pub,
			})
		if res.RowsAffected == 1 {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *ArticleRepository) Stats() (map[string]interface{}, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	_ = r.db.Model(&model.CommunityArticle{}).
		Select("status, count(*) as cnt").
		Where("status <> ?", model.ArticleDeleted).
		Group("status").Scan(&rows).Error

	byStatus := map[string]int64{}
	for _, x := range rows {
		byStatus[x.Status] = x.Cnt
	}

	var auditPending int64
	_ = r.db.Model(&model.CommunityArticle{}).
		Where("audit_status = ? AND status <> ?", model.ArticleAuditPending, model.ArticleDeleted).
		Count(&auditPending).Error

	var viewSum, likeSum int64
	_ = r.db.Model(&model.CommunityArticle{}).
		Where("status <> ?", model.ArticleDeleted).
		Select("COALESCE(SUM(view_count),0), COALESCE(SUM(like_count),0)").
		Row().Scan(&viewSum, &likeSum)

	type dayRow struct {
		Day string `json:"day"`
		Cnt int64  `json:"cnt"`
	}
	var recent []dayRow
	_ = r.db.Raw(`
		SELECT DATE(created_at) AS day, COUNT(*) AS cnt
		FROM community_article
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND status <> ?
		GROUP BY DATE(created_at) ORDER BY day ASC
	`, model.ArticleDeleted).Scan(&recent).Error

	return map[string]interface{}{
		"by_status":     byStatus,
		"audit_pending": auditPending,
		"view_sum":      viewSum,
		"like_sum":      likeSum,
		"recent_7d":     recent,
	}, nil
}

// ---- categories ----

func (r *ArticleRepository) ListCategories() ([]model.CommunityArticleCategory, error) {
	var list []model.CommunityArticleCategory
	err := r.db.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ArticleRepository) CreateCategory(c *model.CommunityArticleCategory) error {
	return r.db.Create(c).Error
}

func (r *ArticleRepository) UpdateCategory(id uint64, updates map[string]interface{}) error {
	res := r.db.Model(&model.CommunityArticleCategory{}).Where("id = ?", id).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("分类不存在")
	}
	return res.Error
}

func (r *ArticleRepository) DeleteCategory(id uint64) error {
	var child int64
	_ = r.db.Model(&model.CommunityArticleCategory{}).Where("parent_id = ?", id).Count(&child).Error
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	var used int64
	_ = r.db.Model(&model.CommunityArticle{}).
		Where("category_id = ? AND status <> ?", id, model.ArticleDeleted).Count(&used).Error
	if used > 0 {
		return errors.New("分类下仍有文章")
	}
	return r.db.Where("id = ?", id).Delete(&model.CommunityArticleCategory{}).Error
}

// ---- comments ----

type CommentListFilter struct {
	ShopID    uint64
	ArticleID uint64
	Status    string
	Page      int
	PageSize  int
}

func (r *ArticleRepository) ListComments(f CommentListFilter) ([]model.CommunityArticleComment, int64, error) {
	q := r.db.Model(&model.CommunityArticleComment{}).Where("status <> ?", model.CommentDeleted)
	if f.ShopID > 0 {
		q = q.Where("shop_id = ?", f.ShopID)
	}
	if f.ArticleID > 0 {
		q = q.Where("article_id = ?", f.ArticleID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.CommunityArticleComment
	err := q.Order("id DESC").Offset((f.Page-1)*f.PageSize).Limit(f.PageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) PatchComment(id uint64, shopID uint64, status string) error {
	q := r.db.Model(&model.CommunityArticleComment{}).Where("id = ?", id)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Update("status", status)
	if res.RowsAffected == 0 {
		return errors.New("评论不存在或无权操作")
	}
	return res.Error
}

func (r *ArticleRepository) DeleteComment(id uint64, shopID uint64) error {
	return r.PatchComment(id, shopID, model.CommentDeleted)
}
