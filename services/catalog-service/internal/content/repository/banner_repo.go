package repository

import (
	"context"
	"errors"
	"time"

	"mymall/services/catalog-service/internal/content/model"
)

const homepageBannerColumns = "id, title, image_url, link_type, link_id, sort, status, start_at, end_at, created_at, updated_at"

func (r *ArticleRepository) ListBannersAdmin(ctx context.Context, page, pageSize int) ([]model.HomepageBanner, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM homepage_banners")
	if err != nil {
		return nil, 0, err
	}
	var list []model.HomepageBanner
	err = r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+homepageBannerColumns+" FROM homepage_banners ORDER BY sort ASC, id DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *ArticleRepository) ListBannersPublic(ctx context.Context) ([]model.HomepageBanner, error) {
	now := time.Now()
	var list []model.HomepageBanner
	err := r.conn.QueryRowsPartialCtx(ctx, &list,
		"SELECT "+homepageBannerColumns+` FROM homepage_banners
WHERE status=? AND (start_at IS NULL OR start_at<=?) AND (end_at IS NULL OR end_at>?)
ORDER BY sort ASC, id DESC`,
		model.BannerOn, now, now,
	)
	return list, err
}

func (r *ArticleRepository) GetBanner(ctx context.Context, id uint64) (*model.HomepageBanner, error) {
	var b model.HomepageBanner
	err := r.conn.QueryRowPartialCtx(ctx, &b, "SELECT "+homepageBannerColumns+" FROM homepage_banners WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *ArticleRepository) CreateBanner(ctx context.Context, b *model.HomepageBanner) error {
	id, err := lastInsertID(ctx, r.conn, `
INSERT INTO homepage_banners (title, image_url, link_type, link_id, sort, status, start_at, end_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		b.Title, b.ImageURL, b.LinkType, b.LinkID, b.Sort, b.Status, b.StartAt, b.EndAt,
	)
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (r *ArticleRepository) UpdateBanner(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("homepage_banners", updates, "id=?", id)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("Banner 不存在")
	}
	return nil
}

func (r *ArticleRepository) DeleteBanner(ctx context.Context, id uint64) error {
	n, err := execAffected(ctx, r.conn, "DELETE FROM homepage_banners WHERE id=?", id)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("Banner 不存在")
	}
	return nil
}

func (r *ArticleRepository) FillBannerLinkNames(ctx context.Context, list []model.HomepageBanner) {
	for i := range list {
		switch list[i].LinkType {
		case model.BannerLinkProduct:
			if list[i].LinkID == 0 {
				continue
			}
			var name string
			_ = r.conn.QueryRowPartialCtx(ctx, &name, "SELECT name FROM products WHERE id=? LIMIT 1", list[i].LinkID)
			list[i].LinkName = name
		case model.BannerLinkArticle:
			if list[i].LinkID == 0 {
				continue
			}
			var title string
			_ = r.conn.QueryRowPartialCtx(ctx, &title, "SELECT title FROM community_article WHERE id=? LIMIT 1", list[i].LinkID)
			list[i].LinkName = title
		}
	}
}

func (r *ArticleRepository) ProductExistsOnSale(ctx context.Context, id uint64) bool {
	n, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM products WHERE id=? AND status=?", id, "on_sale")
	return n > 0
}

func (r *ArticleRepository) ArticleExistsPublished(ctx context.Context, id uint64) bool {
	n, _ := countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM community_article WHERE id=? AND status=? AND audit_status=?",
		id, model.ArticlePublished, model.ArticleAuditApproved,
	)
	return n > 0
}
