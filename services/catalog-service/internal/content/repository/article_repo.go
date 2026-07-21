package repository

import (
	"context"
	"errors"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/content/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

type ArticleListFilter struct {
	ShopID      uint64
	FilterShop  bool // true 时按 ShopID 精确筛选（含平台文章 shop_id=0）
	Title       string
	AuditStatus string
	Status      string
	HasSchedule *bool // nil=不限, true=有定时, false=无定时
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Recycle     bool
	Page        int
	PageSize    int
}

func (r *ArticleRepository) List(ctx context.Context, f ArticleListFilter) ([]model.CommunityArticle, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.CommunityArticle{})
	if f.Recycle {
		q = q.Where("status = ?", model.ArticleDeleted)
	} else {
		q = q.Where("status <> ?", model.ArticleDeleted)
	}
	if f.FilterShop {
		q = q.Where("shop_id = ?", f.ShopID)
	} else if f.ShopID > 0 {
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

func (r *ArticleRepository) GetByID(ctx context.Context, id uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) GetByIDShop(ctx context.Context, id, shopID uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.WithContext(ctx).Where("id = ? AND shop_id = ?", id, shopID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) GetByIDAuthor(ctx context.Context, id, userID uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.WithContext(ctx).Where("id = ? AND author_user_id = ? AND shop_id = 0", id, userID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) ListByAuthor(ctx context.Context, userID uint64, page, pageSize int) ([]model.CommunityArticle, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.CommunityArticle{}).
		Where("author_user_id = ? AND shop_id = 0 AND status <> ?", userID, model.ArticleDeleted)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var list []model.CommunityArticle
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) UpdateAuthor(ctx context.Context, id, userID uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.CommunityArticle{}).
		Where("id = ? AND author_user_id = ? AND shop_id = 0", id, userID).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("文章不存在或无权操作")
	}
	return res.Error
}

func (r *ArticleRepository) SoftDeleteAuthor(ctx context.Context, id, userID uint64) error {
	now := common.LocalTime(time.Now())
	return r.UpdateAuthor(ctx, id, userID, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) Create(ctx context.Context, a *model.CommunityArticle) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *ArticleRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.CommunityArticle{}).Where("id = ?", id).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return res.Error
}

func (r *ArticleRepository) UpdateShop(ctx context.Context, id, shopID uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.CommunityArticle{}).Where("id = ? AND shop_id = ?", id, shopID).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("文章不存在或无权操作")
	}
	return res.Error
}

func (r *ArticleRepository) SoftDelete(ctx context.Context, id uint64) error {
	now := common.LocalTime(time.Now())
	return r.Update(ctx, id, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) SoftDeleteShop(ctx context.Context, id, shopID uint64) error {
	now := common.LocalTime(time.Now())
	return r.UpdateShop(ctx, id, shopID, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) Restore(ctx context.Context, id uint64) error {
	return r.Update(ctx, id, map[string]interface{}{
		"status": model.ArticleOffline, "deleted_at": nil,
	})
}

func (r *ArticleRepository) PermanentDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", id).Delete(&model.CommunityArticleComment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("article_id = ?", id).Delete(&model.CommunityArticleImg{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&model.CommunityArticle{}).Error
	})
}

func (r *ArticleRepository) ReplaceImages(ctx context.Context, articleID, shopID uint64, urls []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (r *ArticleRepository) ListImages(ctx context.Context, articleID uint64) ([]model.CommunityArticleImg, error) {
	var list []model.CommunityArticleImg
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ArticleRepository) ClaimDuePublish(ctx context.Context, limit int) ([]model.CommunityArticle, error) {
	var list []model.CommunityArticle
	now := time.Now()
	err := r.db.WithContext(ctx).Where("status = ? AND audit_status = ? AND schedule_publish_at IS NOT NULL AND schedule_publish_at <= ?",
		model.ArticleScheduled, model.ArticleAuditApproved, now).
		Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make([]model.CommunityArticle, 0, len(list))
	pub := common.LocalTime(now)
	for _, a := range list {
		res := r.db.WithContext(ctx).Model(&model.CommunityArticle{}).
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

func (r *ArticleRepository) Stats(ctx context.Context) (map[string]interface{}, error) {
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

func (r *ArticleRepository) ListCategories(ctx context.Context) ([]model.CommunityArticleCategory, error) {
	var list []model.CommunityArticleCategory
	err := r.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ArticleRepository) CreateCategory(ctx context.Context, c *model.CommunityArticleCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ArticleRepository) UpdateCategory(ctx context.Context, id uint64, updates map[string]interface{}) error {
	res := r.db.WithContext(ctx).Model(&model.CommunityArticleCategory{}).Where("id = ?", id).Updates(updates)
	if res.RowsAffected == 0 {
		return errors.New("分类不存在")
	}
	return res.Error
}

func (r *ArticleRepository) DeleteCategory(ctx context.Context, id uint64) error {
	var child int64
	_ = r.db.WithContext(ctx).Model(&model.CommunityArticleCategory{}).Where("parent_id = ?", id).Count(&child).Error
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	var used int64
	_ = r.db.WithContext(ctx).Model(&model.CommunityArticle{}).
		Where("category_id = ? AND status <> ?", id, model.ArticleDeleted).Count(&used).Error
	if used > 0 {
		return errors.New("分类下仍有文章")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.CommunityArticleCategory{}).Error
}

// ---- comments ----

type CommentListFilter struct {
	ShopID    uint64
	ArticleID uint64
	Status    string
	Page      int
	PageSize  int
}

func (r *ArticleRepository) ListComments(ctx context.Context, f CommentListFilter) ([]model.CommunityArticleComment, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.CommunityArticleComment{}).Where("status <> ?", model.CommentDeleted)
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
	err := q.Order("id DESC").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) PatchComment(ctx context.Context, id uint64, shopID uint64, status string) error {
	q := r.db.WithContext(ctx).Model(&model.CommunityArticleComment{}).Where("id = ?", id)
	if shopID > 0 {
		q = q.Where("shop_id = ?", shopID)
	}
	res := q.Update("status", status)
	if res.RowsAffected == 0 {
		return errors.New("评论不存在或无权操作")
	}
	return res.Error
}

func (r *ArticleRepository) DeleteComment(ctx context.Context, id uint64, shopID uint64) error {
	return r.PatchComment(ctx, id, shopID, model.CommentDeleted)
}

func (r *ArticleRepository) GetComment(ctx context.Context, id uint64) (*model.CommunityArticleComment, error) {
	var c model.CommunityArticleComment
	if err := r.db.WithContext(ctx).Where("id = ? AND status = ?", id, model.CommentVisible).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ArticleRepository) CreateComment(ctx context.Context, c *model.CommunityArticleComment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if c.RootID == 0 {
			if err := tx.Model(c).Update("root_id", c.ID).Error; err != nil {
				return err
			}
			c.RootID = c.ID
		}
		return tx.Model(&model.CommunityArticle{}).Where("id = ?", c.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
	})
}

type CommentUserBrief struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
}

func (r *ArticleRepository) MapUserBriefs(ctx context.Context, ids []uint64) map[uint64]CommentUserBrief {
	out := map[uint64]CommentUserBrief{}
	if len(ids) == 0 {
		return out
	}
	var rows []CommentUserBrief
	_ = r.db.WithContext(ctx).Table("users").Select("id, nickname, COALESCE(avatar,'') AS avatar, mobile").Where("id IN ?", ids).Scan(&rows).Error
	for _, u := range rows {
		out[u.ID] = u
	}
	return out
}

type ShopBrief struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

func (r *ArticleRepository) MapShopBriefs(ctx context.Context, ids []uint64) map[uint64]ShopBrief {
	out := map[uint64]ShopBrief{}
	if len(ids) == 0 {
		return out
	}
	var rows []ShopBrief
	_ = r.db.WithContext(ctx).Table("shops").Select("id, name, COALESCE(logo,'') AS logo").Where("id IN ?", ids).Scan(&rows).Error
	for _, s := range rows {
		out[s.ID] = s
	}
	return out
}

func (r *ArticleRepository) ListPublicCommentRoots(ctx context.Context, articleID uint64, page, pageSize int) ([]model.CommunityArticleComment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := r.db.WithContext(ctx).Model(&model.CommunityArticleComment{}).
		Where("article_id = ? AND parent_id = 0 AND status = ?", articleID, model.CommentVisible)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CommunityArticleComment
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) ListPublicCommentChildren(ctx context.Context, articleID uint64, rootIDs []uint64) ([]model.CommunityArticleComment, error) {
	if len(rootIDs) == 0 {
		return nil, nil
	}
	var list []model.CommunityArticleComment
	err := r.db.WithContext(ctx).Where("article_id = ? AND root_id IN ? AND parent_id > 0 AND status = ?",
		articleID, rootIDs, model.CommentVisible).
		Order("id ASC").Find(&list).Error
	return list, err
}

// ---- emojis ----

func (r *ArticleRepository) ListEmojisAdmin(ctx context.Context, page, pageSize int) ([]model.CommunityCommentEmoji, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	q := r.db.WithContext(ctx).Model(&model.CommunityCommentEmoji{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CommunityCommentEmoji
	err := q.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) ListEmojisPublic(ctx context.Context) ([]model.CommunityCommentEmoji, error) {
	var list []model.CommunityCommentEmoji
	err := r.db.WithContext(ctx).Where("status = 1").Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ArticleRepository) GetEmoji(ctx context.Context, id uint64) (*model.CommunityCommentEmoji, error) {
	var e model.CommunityCommentEmoji
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ArticleRepository) CreateEmoji(ctx context.Context, e *model.CommunityCommentEmoji) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *ArticleRepository) UpdateEmoji(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.CommunityCommentEmoji{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ArticleRepository) DeleteEmoji(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.CommunityCommentEmoji{}).Error
}

func (r *ArticleRepository) ListPublic(ctx context.Context, page, pageSize, homeLimit int) ([]model.CommunityArticle, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	now := time.Now()
	_ = r.db.WithContext(ctx).Exec(`UPDATE homepage_slot_orders SET status='expired' WHERE slot_type='article' AND status='active' AND end_at < ?`, now)

	type row struct {
		model.CommunityArticle
		Boost int `gorm:"column:boost"`
	}
	base := `
FROM community_article a
LEFT JOIN (
  SELECT DISTINCT target_id FROM homepage_slot_orders
  WHERE slot_type='article' AND status='active' AND start_at<=? AND end_at>?
) o ON o.target_id = a.id
WHERE a.status=? AND a.audit_status=? AND a.deleted_at IS NULL`

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*) "+base, now, now, model.ArticlePublished, model.ArticleAuditApproved).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	limit, offset := pageSize, (page-1)*pageSize
	if homeLimit > 0 {
		limit, offset = homeLimit, 0
		if int64(homeLimit) < total {
			total = int64(homeLimit)
		}
	}
	var rows []row
	sql := `SELECT a.*, CASE WHEN o.target_id IS NULL THEN 0 ELSE 1 END AS boost ` + base + `
ORDER BY boost DESC, a.is_top DESC, a.id DESC LIMIT ? OFFSET ?`
	if err := r.db.WithContext(ctx).Raw(sql, now, now, model.ArticlePublished, model.ArticleAuditApproved, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]model.CommunityArticle, 0, len(rows))
	for _, x := range rows {
		out = append(out, x.CommunityArticle)
	}
	return out, total, nil
}

func (r *ArticleRepository) GetPublished(ctx context.Context, id uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.db.WithContext(ctx).Where("id = ? AND status = ? AND audit_status = ?", id, model.ArticlePublished, model.ArticleAuditApproved).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) RecordRead(ctx context.Context, articleID, userID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
			Updates(map[string]interface{}{
				"read_count": gorm.Expr("read_count + 1"),
				"view_count": gorm.Expr("view_count + 1"),
			}).Error; err != nil {
			return err
		}
		if userID == 0 {
			return nil
		}
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ArticleAudience{
			UserID: userID, ArticleID: articleID,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			return tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
				UpdateColumn("audience_count", gorm.Expr("audience_count + 1")).Error
		}
		return nil
	})
}

func (r *ArticleRepository) ToggleLike(ctx context.Context, userID, articleID uint64, like bool) (changed bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if like {
			res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ArticleLike{UserID: userID, ArticleID: articleID})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				changed = true
				return tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
					UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
			}
			return nil
		}
		res := tx.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&model.ArticleLike{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			changed = true
			return tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
				UpdateColumn("like_count", gorm.Expr("GREATEST(0, like_count - 1)")).Error
		}
		return nil
	})
	return changed, err
}

func (r *ArticleRepository) ToggleFavorite(ctx context.Context, userID, articleID uint64, fav bool) (changed bool, err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if fav {
			res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ArticleFavorite{UserID: userID, ArticleID: articleID})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				changed = true
				return tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
					UpdateColumn("collect_count", gorm.Expr("collect_count + 1")).Error
			}
			return nil
		}
		res := tx.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&model.ArticleFavorite{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected > 0 {
			changed = true
			return tx.Model(&model.CommunityArticle{}).Where("id = ?", articleID).
				UpdateColumn("collect_count", gorm.Expr("GREATEST(0, collect_count - 1)")).Error
		}
		return nil
	})
	return changed, err
}

func (r *ArticleRepository) EngagementStatus(ctx context.Context, userID, articleID uint64) (liked, favorited bool) {
	var n int64
	r.db.WithContext(ctx).Model(&model.ArticleLike{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&n)
	liked = n > 0
	r.db.WithContext(ctx).Model(&model.ArticleFavorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&n)
	favorited = n > 0
	return
}

// UserArticleItem 用户收藏/点赞列表项
type UserArticleItem struct {
	ID           uint64 `json:"id"`
	ShopID       uint64 `json:"shop_id"`
	Title        string `json:"title"`
	CoverURL     string `json:"cover_url"`
	LikeCount    uint64 `json:"like_count"`
	ReadCount    uint64 `json:"read_count"`
	CollectCount uint64 `json:"collect_count"`
	Status       string `json:"status"`
	AuditStatus  string `json:"audit_status"`
	Invalid      bool   `json:"invalid"`
	EngagedAt    string `json:"engaged_at,omitempty"`
}

func (r *ArticleRepository) listUserArticles(ctx context.Context, userID uint64, table string, page, pageSize int) ([]UserArticleItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var total int64
	if err := r.db.WithContext(ctx).Table(table).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	type row struct {
		ArticleID    uint64 `gorm:"column:article_id"`
		EngagedAt    string `gorm:"column:engaged_at"`
		ID           uint64 `gorm:"column:id"`
		ShopID       uint64 `gorm:"column:shop_id"`
		Title        string `gorm:"column:title"`
		CoverURL     string `gorm:"column:cover_url"`
		LikeCount    uint64 `gorm:"column:like_count"`
		ReadCount    uint64 `gorm:"column:read_count"`
		CollectCount uint64 `gorm:"column:collect_count"`
		Status       string `gorm:"column:status"`
		AuditStatus  string `gorm:"column:audit_status"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
SELECT e.article_id, DATE_FORMAT(e.created_at, '%Y-%m-%d %H:%i:%s') AS engaged_at,
       a.id, a.shop_id, COALESCE(a.title,'') AS title, COALESCE(a.cover_url,'') AS cover_url,
       COALESCE(a.like_count,0) AS like_count, COALESCE(a.read_count,0) AS read_count,
       COALESCE(a.collect_count,0) AS collect_count,
       COALESCE(a.status,'') AS status, COALESCE(a.audit_status,'') AS audit_status
FROM `+table+` e
LEFT JOIN community_article a ON a.id = e.article_id
WHERE e.user_id = ?
ORDER BY e.id DESC
LIMIT ? OFFSET ?`, userID, pageSize, (page-1)*pageSize).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]UserArticleItem, 0, len(rows))
	for _, x := range rows {
		invalid := x.ID == 0 || x.Status != model.ArticlePublished || x.AuditStatus != model.ArticleAuditApproved
		title := x.Title
		if x.ID == 0 {
			title = "文章已失效"
		} else if invalid && title == "" {
			title = "文章已失效"
		}
		out = append(out, UserArticleItem{
			ID: x.ArticleID, ShopID: x.ShopID, Title: title, CoverURL: x.CoverURL,
			LikeCount: x.LikeCount, ReadCount: x.ReadCount, CollectCount: x.CollectCount,
			Status: x.Status, AuditStatus: x.AuditStatus, Invalid: invalid, EngagedAt: x.EngagedAt,
		})
	}
	return out, total, nil
}

func (r *ArticleRepository) ListUserFavorites(ctx context.Context, userID uint64, page, pageSize int) ([]UserArticleItem, int64, error) {
	return r.listUserArticles(ctx, userID, "article_favorites", page, pageSize)
}

func (r *ArticleRepository) ListUserLikes(ctx context.Context, userID uint64, page, pageSize int) ([]UserArticleItem, int64, error) {
	return r.listUserArticles(ctx, userID, "article_likes", page, pageSize)
}

func (r *ArticleRepository) IsArticleBoosted(ctx context.Context, articleID uint64) bool {
	now := time.Now()
	var n int64
	r.db.WithContext(ctx).Table("homepage_slot_orders").
		Where("slot_type='article' AND status='active' AND target_id=? AND start_at<=? AND end_at>?", articleID, now, now).
		Count(&n)
	return n > 0
}

func (r *ArticleRepository) GetHomeArticleLimit(ctx context.Context) int {
	var lim int
	_ = r.db.WithContext(ctx).Table("homepage_slot_settings").Select("home_limit").Where("slot_type=?", "article").Scan(&lim).Error
	if lim < 1 {
		return 6
	}
	return lim
}
