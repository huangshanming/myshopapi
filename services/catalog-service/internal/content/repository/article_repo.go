package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/content/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	communityArticleColumns = "id, shop_id, author_user_id, category_id, title, cover_url, content, audit_status, reject_reason, status, schedule_publish_at, is_top, view_count, like_count, audience_count, read_count, collect_count, comment_count, published_at, deleted_at, created_by, created_at, updated_at"
	communityArticleCategoryColumns = "id, parent_id, name, sort, status, created_at, updated_at"
	communityArticleCommentColumns = "id, article_id, shop_id, user_id, parent_id, root_id, reply_to_user_id, content, status, created_at, updated_at"
	communityArticleImgColumns = "id, article_id, shop_id, url, sort, created_at"
	communityCommentEmojiColumns = "id, name, image_url, sort, status, created_at, updated_at"
)

type ArticleRepository struct {
	conn sqlx.SqlConn
}

func NewArticleRepository(conn sqlx.SqlConn) *ArticleRepository {
	return &ArticleRepository{conn: conn}
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
	where := make([]string, 0, 8)
	args := make([]any, 0, 12)
	if f.Recycle {
		where = append(where, "status=?")
		args = append(args, model.ArticleDeleted)
	} else {
		where = append(where, "status<>?")
		args = append(args, model.ArticleDeleted)
	}
	if f.FilterShop {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	} else if f.ShopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	}
	if f.Title != "" {
		where = append(where, "title LIKE ?")
		args = append(args, "%"+f.Title+"%")
	}
	if f.AuditStatus != "" {
		where = append(where, "audit_status=?")
		args = append(args, f.AuditStatus)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.HasSchedule != nil {
		if *f.HasSchedule {
			where = append(where, "schedule_publish_at IS NOT NULL")
		} else {
			where = append(where, "schedule_publish_at IS NULL")
		}
	}
	if f.CreatedFrom != nil {
		where = append(where, "created_at >= ?")
		args = append(args, *f.CreatedFrom)
	}
	if f.CreatedTo != nil {
		where = append(where, "created_at <= ?")
		args = append(args, *f.CreatedTo)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_article WHERE "+whereSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.CommunityArticle
	qArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleColumns+" FROM community_article WHERE "+whereSQL+" ORDER BY is_top DESC, id DESC LIMIT ? OFFSET ?",
		qArgs...,
	)
	return list, total, err
}

func (r *ArticleRepository) GetByID(ctx context.Context, id uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.conn.QueryRowCtx(ctx, &a, "SELECT "+communityArticleColumns+" FROM community_article WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) GetByIDShop(ctx context.Context, id, shopID uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.conn.QueryRowCtx(ctx, &a, "SELECT "+communityArticleColumns+" FROM community_article WHERE id=? AND shop_id=? LIMIT 1", id, shopID)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) GetByIDAuthor(ctx context.Context, id, userID uint64) (*model.CommunityArticle, error) {
	var a model.CommunityArticle
	err := r.conn.QueryRowCtx(ctx, &a,
		"SELECT "+communityArticleColumns+" FROM community_article WHERE id=? AND author_user_id=? AND shop_id=0 LIMIT 1",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) ListByAuthor(ctx context.Context, userID uint64, page, pageSize int) ([]model.CommunityArticle, int64, error) {
	whereSQL := "author_user_id=? AND shop_id=0 AND status<>?"
	args := []any{userID, model.ArticleDeleted}
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_article WHERE "+whereSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var list []model.CommunityArticle
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleColumns+" FROM community_article WHERE "+whereSQL+" ORDER BY id DESC LIMIT ? OFFSET ?",
		userID, model.ArticleDeleted, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *ArticleRepository) UpdateAuthor(ctx context.Context, id, userID uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("community_article", updates, "id=? AND author_user_id=? AND shop_id=0", id, userID)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("文章不存在或无权操作")
	}
	return nil
}

func (r *ArticleRepository) SoftDeleteAuthor(ctx context.Context, id, userID uint64) error {
	now := common.LocalTime(time.Now())
	return r.UpdateAuthor(ctx, id, userID, map[string]interface{}{
		"status": model.ArticleDeleted, "deleted_at": now,
	})
}

func (r *ArticleRepository) Create(ctx context.Context, a *model.CommunityArticle) error {
	id, err := lastInsertID(ctx, r.conn, `
INSERT INTO community_article (shop_id, author_user_id, category_id, title, cover_url, content, audit_status, reject_reason, status, schedule_publish_at, is_top, view_count, like_count, audience_count, read_count, collect_count, comment_count, published_at, deleted_at, created_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ShopID, a.AuthorUserID, a.CategoryID, a.Title, a.CoverURL, a.Content,
		a.AuditStatus, a.RejectReason, a.Status, a.SchedulePublishAt, a.IsTop,
		a.ViewCount, a.LikeCount, a.AudienceCount, a.ReadCount, a.CollectCount, a.CommentCount,
		a.PublishedAt, a.DeletedAt, a.CreatedBy,
	)
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (r *ArticleRepository) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("community_article", updates, "id=?", id)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("文章不存在")
	}
	return nil
}

func (r *ArticleRepository) UpdateShop(ctx context.Context, id, shopID uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("community_article", updates, "id=? AND shop_id=?", id, shopID)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("文章不存在或无权操作")
	}
	return nil
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
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := execSession(ctx, session, "DELETE FROM community_article_comment WHERE article_id=?", id); err != nil {
			return err
		}
		if _, err := execSession(ctx, session, "DELETE FROM community_article_img WHERE article_id=?", id); err != nil {
			return err
		}
		_, err := execSession(ctx, session, "DELETE FROM community_article WHERE id=?", id)
		return err
	})
}

func (r *ArticleRepository) ReplaceImages(ctx context.Context, articleID, shopID uint64, urls []string) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := execSession(ctx, session, "DELETE FROM community_article_img WHERE article_id=?", articleID); err != nil {
			return err
		}
		for i, u := range urls {
			if u == "" {
				continue
			}
			if _, err := session.ExecCtx(ctx,
				"INSERT INTO community_article_img (article_id, shop_id, url, sort) VALUES (?, ?, ?, ?)",
				articleID, shopID, u, i,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ArticleRepository) ListImages(ctx context.Context, articleID uint64) ([]model.CommunityArticleImg, error) {
	var list []model.CommunityArticleImg
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleImgColumns+" FROM community_article_img WHERE article_id=? ORDER BY sort ASC, id ASC",
		articleID,
	)
	return list, err
}

func (r *ArticleRepository) ClaimDuePublish(ctx context.Context, limit int) ([]model.CommunityArticle, error) {
	now := time.Now()
	var list []model.CommunityArticle
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleColumns+` FROM community_article
WHERE status=? AND audit_status=? AND schedule_publish_at IS NOT NULL AND schedule_publish_at<=?
LIMIT ?`,
		model.ArticleScheduled, model.ArticleAuditApproved, now, limit,
	)
	if err != nil {
		return nil, err
	}
	out := make([]model.CommunityArticle, 0, len(list))
	pub := common.LocalTime(now)
	for _, a := range list {
		n, err := execAffected(ctx, r.conn,
			"UPDATE community_article SET status=?, published_at=? WHERE id=? AND status=?",
			model.ArticlePublished, pub, a.ID, model.ArticleScheduled,
		)
		if err != nil {
			return nil, err
		}
		if n == 1 {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *ArticleRepository) Stats(ctx context.Context) (map[string]interface{}, error) {
	type statusRow struct {
		Status string `db:"status"`
		Cnt    int64  `db:"cnt"`
	}
	var rows []statusRow
	_ = r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT status, COUNT(*) AS cnt FROM community_article WHERE status<>? GROUP BY status",
		model.ArticleDeleted,
	)

	byStatus := map[string]int64{}
	for _, x := range rows {
		byStatus[x.Status] = x.Cnt
	}

	auditPending, _ := countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM community_article WHERE audit_status=? AND status<>?",
		model.ArticleAuditPending, model.ArticleDeleted,
	)

	var viewSum, likeSum int64
	_ = r.conn.QueryRowCtx(ctx, &viewSum,
		"SELECT COALESCE(SUM(view_count),0) FROM community_article WHERE status<>?",
		model.ArticleDeleted,
	)
	_ = r.conn.QueryRowCtx(ctx, &likeSum,
		"SELECT COALESCE(SUM(like_count),0) FROM community_article WHERE status<>?",
		model.ArticleDeleted,
	)

	type dayRow struct {
		Day string `db:"day" json:"day"`
		Cnt int64  `db:"cnt" json:"cnt"`
	}
	var recent []dayRow
	_ = r.conn.QueryRowsCtx(ctx, &recent, `
		SELECT DATE(created_at) AS day, COUNT(*) AS cnt
		FROM community_article
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND status <> ?
		GROUP BY DATE(created_at) ORDER BY day ASC
	`, model.ArticleDeleted)

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
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleCategoryColumns+" FROM community_article_category ORDER BY sort ASC, id ASC",
	)
	return list, err
}

func (r *ArticleRepository) CreateCategory(ctx context.Context, c *model.CommunityArticleCategory) error {
	id, err := lastInsertID(ctx, r.conn,
		"INSERT INTO community_article_category (parent_id, name, sort, status) VALUES (?, ?, ?, ?)",
		c.ParentID, c.Name, c.Sort, c.Status,
	)
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *ArticleRepository) UpdateCategory(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("community_article_category", updates, "id=?", id)
	if err != nil {
		return err
	}
	n, err := execAffected(ctx, r.conn, query, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("分类不存在")
	}
	return nil
}

func (r *ArticleRepository) DeleteCategory(ctx context.Context, id uint64) error {
	child, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_article_category WHERE parent_id=?", id)
	if child > 0 {
		return errors.New("请先删除子分类")
	}
	used, _ := countCtx(ctx, r.conn,
		"SELECT COUNT(*) FROM community_article WHERE category_id=? AND status<>?",
		id, model.ArticleDeleted,
	)
	if used > 0 {
		return errors.New("分类下仍有文章")
	}
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM community_article_category WHERE id=?", id)
	return err
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
	where := []string{"status<>?"}
	args := []any{model.CommentDeleted}
	if f.ShopID > 0 {
		where = append(where, "shop_id=?")
		args = append(args, f.ShopID)
	}
	if f.ArticleID > 0 {
		where = append(where, "article_id=?")
		args = append(args, f.ArticleID)
	}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	whereSQL := strings.Join(where, " AND ")

	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_article_comment WHERE "+whereSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	var list []model.CommunityArticleComment
	qArgs := append(append([]any{}, args...), f.PageSize, (f.Page-1)*f.PageSize)
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleCommentColumns+" FROM community_article_comment WHERE "+whereSQL+" ORDER BY id DESC LIMIT ? OFFSET ?",
		qArgs...,
	)
	return list, total, err
}

func (r *ArticleRepository) PatchComment(ctx context.Context, id uint64, shopID uint64, status string) error {
	where := "id=?"
	args := []any{status, id}
	if shopID > 0 {
		where = "id=? AND shop_id=?"
		args = []any{status, id, shopID}
	}
	n, err := execAffected(ctx, r.conn, "UPDATE community_article_comment SET status=? WHERE "+where, args...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("评论不存在或无权操作")
	}
	return nil
}

func (r *ArticleRepository) DeleteComment(ctx context.Context, id uint64, shopID uint64) error {
	return r.PatchComment(ctx, id, shopID, model.CommentDeleted)
}

func (r *ArticleRepository) GetComment(ctx context.Context, id uint64) (*model.CommunityArticleComment, error) {
	var c model.CommunityArticleComment
	err := r.conn.QueryRowCtx(ctx, &c,
		"SELECT "+communityArticleCommentColumns+" FROM community_article_comment WHERE id=? AND status=? LIMIT 1",
		id, model.CommentVisible,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ArticleRepository) CreateComment(ctx context.Context, c *model.CommunityArticleComment) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		id, err := lastInsertID(ctx, session, `
INSERT INTO community_article_comment (article_id, shop_id, user_id, parent_id, root_id, reply_to_user_id, content, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			c.ArticleID, c.ShopID, c.UserID, c.ParentID, c.RootID, c.ReplyToUserID, c.Content, c.Status,
		)
		if err != nil {
			return err
		}
		c.ID = id
		if c.RootID == 0 {
			if _, err := execSession(ctx, session, "UPDATE community_article_comment SET root_id=? WHERE id=?", id, id); err != nil {
				return err
			}
			c.RootID = id
		}
		_, err = session.ExecCtx(ctx,
			"UPDATE community_article SET comment_count=comment_count+1 WHERE id=?",
			c.ArticleID,
		)
		return err
	})
}

type CommentUserBrief struct {
	ID       uint64 `db:"id" json:"id"`
	Nickname string `db:"nickname" json:"nickname"`
	Avatar   string `db:"avatar" json:"avatar"`
	Mobile   string `db:"mobile" json:"mobile"`
}

func (r *ArticleRepository) MapUserBriefs(ctx context.Context, ids []uint64) map[uint64]CommentUserBrief {
	out := map[uint64]CommentUserBrief{}
	if len(ids) == 0 {
		return out
	}
	var rows []CommentUserBrief
	_ = r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT id, nickname, COALESCE(avatar,'') AS avatar, mobile FROM users WHERE id IN ("+placeholders(len(ids))+")",
		inArgs(ids)...,
	)
	for _, u := range rows {
		out[u.ID] = u
	}
	return out
}

type ShopBrief struct {
	ID   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Logo string `db:"logo" json:"logo"`
}

func (r *ArticleRepository) MapShopBriefs(ctx context.Context, ids []uint64) map[uint64]ShopBrief {
	out := map[uint64]ShopBrief{}
	if len(ids) == 0 {
		return out
	}
	var rows []ShopBrief
	_ = r.conn.QueryRowsCtx(ctx, &rows,
		"SELECT id, name, COALESCE(logo,'') AS logo FROM shops WHERE id IN ("+placeholders(len(ids))+")",
		inArgs(ids)...,
	)
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
	whereSQL := "article_id=? AND parent_id=0 AND status=?"
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_article_comment WHERE "+whereSQL, articleID, model.CommentVisible)
	if err != nil {
		return nil, 0, err
	}
	var list []model.CommunityArticleComment
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleCommentColumns+" FROM community_article_comment WHERE "+whereSQL+" ORDER BY id DESC LIMIT ? OFFSET ?",
		articleID, model.CommentVisible, pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *ArticleRepository) ListPublicCommentChildren(ctx context.Context, articleID uint64, rootIDs []uint64) ([]model.CommunityArticleComment, error) {
	if len(rootIDs) == 0 {
		return nil, nil
	}
	args := append([]any{articleID}, inArgs(rootIDs)...)
	args = append(args, model.CommentVisible)
	var list []model.CommunityArticleComment
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityArticleCommentColumns+" FROM community_article_comment WHERE article_id=? AND root_id IN ("+placeholders(len(rootIDs))+") AND parent_id>0 AND status=? ORDER BY id ASC",
		args...,
	)
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
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM community_comment_emojis")
	if err != nil {
		return nil, 0, err
	}
	var list []model.CommunityCommentEmoji
	err = r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityCommentEmojiColumns+" FROM community_comment_emojis ORDER BY sort ASC, id DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	return list, total, err
}

func (r *ArticleRepository) ListEmojisPublic(ctx context.Context) ([]model.CommunityCommentEmoji, error) {
	var list []model.CommunityCommentEmoji
	err := r.conn.QueryRowsCtx(ctx, &list,
		"SELECT "+communityCommentEmojiColumns+" FROM community_comment_emojis WHERE status=1 ORDER BY sort ASC, id ASC",
	)
	return list, err
}

func (r *ArticleRepository) GetEmoji(ctx context.Context, id uint64) (*model.CommunityCommentEmoji, error) {
	var e model.CommunityCommentEmoji
	err := r.conn.QueryRowCtx(ctx, &e, "SELECT "+communityCommentEmojiColumns+" FROM community_comment_emojis WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ArticleRepository) CreateEmoji(ctx context.Context, e *model.CommunityCommentEmoji) error {
	id, err := lastInsertID(ctx, r.conn,
		"INSERT INTO community_comment_emojis (name, image_url, sort, status) VALUES (?, ?, ?, ?)",
		e.Name, e.ImageURL, e.Sort, e.Status,
	)
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (r *ArticleRepository) UpdateEmoji(ctx context.Context, id uint64, updates map[string]interface{}) error {
	query, args, err := buildUpdate("community_comment_emojis", updates, "id=?", id)
	if err != nil {
		return err
	}
	_, err = r.conn.ExecCtx(ctx, query, args...)
	return err
}

func (r *ArticleRepository) DeleteEmoji(ctx context.Context, id uint64) error {
	_, err := r.conn.ExecCtx(ctx, "DELETE FROM community_comment_emojis WHERE id=?", id)
	return err
}

func articleColumnsAliased(alias string) string {
	cols := strings.Split(communityArticleColumns, ", ")
	for i, c := range cols {
		cols[i] = alias + "." + c
	}
	return strings.Join(cols, ", ")
}

func (r *ArticleRepository) ListPublic(ctx context.Context, page, pageSize, homeLimit int) ([]model.CommunityArticle, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	now := time.Now()
	_, _ = r.conn.ExecCtx(ctx,
		`UPDATE homepage_slot_orders SET status='expired' WHERE slot_type='article' AND status='active' AND end_at < ?`,
		now,
	)

	type articleBoostRow struct {
		model.CommunityArticle
		Boost int `db:"boost"`
	}

	base := `
FROM community_article a
LEFT JOIN (
  SELECT DISTINCT target_id FROM homepage_slot_orders
  WHERE slot_type='article' AND status='active' AND start_at<=? AND end_at>?
) o ON o.target_id = a.id
WHERE a.status=? AND a.audit_status=? AND a.deleted_at IS NULL`

	var total int64
	if err := r.conn.QueryRowCtx(ctx, &total, "SELECT COUNT(*) "+base, now, now, model.ArticlePublished, model.ArticleAuditApproved); err != nil {
		return nil, 0, err
	}
	limit, offset := pageSize, (page-1)*pageSize
	if homeLimit > 0 {
		limit, offset = homeLimit, 0
		if int64(homeLimit) < total {
			total = int64(homeLimit)
		}
	}
	var rows []articleBoostRow
	sql := fmt.Sprintf("SELECT %s, CASE WHEN o.target_id IS NULL THEN 0 ELSE 1 END AS boost ", articleColumnsAliased("a")) + base + `
ORDER BY boost DESC, a.is_top DESC, a.id DESC LIMIT ? OFFSET ?`
	if err := r.conn.QueryRowsCtx(ctx, &rows, sql, now, now, model.ArticlePublished, model.ArticleAuditApproved, limit, offset); err != nil {
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
	err := r.conn.QueryRowCtx(ctx, &a,
		"SELECT "+communityArticleColumns+" FROM community_article WHERE id=? AND status=? AND audit_status=? LIMIT 1",
		id, model.ArticlePublished, model.ArticleAuditApproved,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) RecordRead(ctx context.Context, articleID, userID uint64) error {
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx,
			"UPDATE community_article SET read_count=read_count+1, view_count=view_count+1 WHERE id=?",
			articleID,
		); err != nil {
			return err
		}
		if userID == 0 {
			return nil
		}
		n, err := execSession(ctx, session,
			"INSERT IGNORE INTO article_audiences (user_id, article_id) VALUES (?, ?)",
			userID, articleID,
		)
		if err != nil {
			return err
		}
		if n > 0 {
			_, err = session.ExecCtx(ctx,
				"UPDATE community_article SET audience_count=audience_count+1 WHERE id=?",
				articleID,
			)
			return err
		}
		return nil
	})
}

func (r *ArticleRepository) ToggleLike(ctx context.Context, userID, articleID uint64, like bool) (changed bool, err error) {
	err = r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if like {
			n, err := execSession(ctx, session,
				"INSERT IGNORE INTO article_likes (user_id, article_id) VALUES (?, ?)",
				userID, articleID,
			)
			if err != nil {
				return err
			}
			if n > 0 {
				changed = true
				_, err = session.ExecCtx(ctx,
					"UPDATE community_article SET like_count=like_count+1 WHERE id=?",
					articleID,
				)
				return err
			}
			return nil
		}
		n, err := execSession(ctx, session,
			"DELETE FROM article_likes WHERE user_id=? AND article_id=?",
			userID, articleID,
		)
		if err != nil {
			return err
		}
		if n > 0 {
			changed = true
			_, err = session.ExecCtx(ctx,
				"UPDATE community_article SET like_count=GREATEST(0, like_count-1) WHERE id=?",
				articleID,
			)
			return err
		}
		return nil
	})
	return changed, err
}

func (r *ArticleRepository) ToggleFavorite(ctx context.Context, userID, articleID uint64, fav bool) (changed bool, err error) {
	err = r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if fav {
			n, err := execSession(ctx, session,
				"INSERT IGNORE INTO article_favorites (user_id, article_id) VALUES (?, ?)",
				userID, articleID,
			)
			if err != nil {
				return err
			}
			if n > 0 {
				changed = true
				_, err = session.ExecCtx(ctx,
					"UPDATE community_article SET collect_count=collect_count+1 WHERE id=?",
					articleID,
				)
				return err
			}
			return nil
		}
		n, err := execSession(ctx, session,
			"DELETE FROM article_favorites WHERE user_id=? AND article_id=?",
			userID, articleID,
		)
		if err != nil {
			return err
		}
		if n > 0 {
			changed = true
			_, err = session.ExecCtx(ctx,
				"UPDATE community_article SET collect_count=GREATEST(0, collect_count-1) WHERE id=?",
				articleID,
			)
			return err
		}
		return nil
	})
	return changed, err
}

func (r *ArticleRepository) EngagementStatus(ctx context.Context, userID, articleID uint64) (liked, favorited bool) {
	n, _ := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM article_likes WHERE user_id=? AND article_id=?", userID, articleID)
	liked = n > 0
	n, _ = countCtx(ctx, r.conn, "SELECT COUNT(*) FROM article_favorites WHERE user_id=? AND article_id=?", userID, articleID)
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
	total, err := countCtx(ctx, r.conn, "SELECT COUNT(*) FROM "+table+" WHERE user_id=?", userID)
	if err != nil {
		return nil, 0, err
	}
	type row struct {
		ArticleID    uint64 `db:"article_id"`
		EngagedAt    string `db:"engaged_at"`
		ID           uint64 `db:"id"`
		ShopID       uint64 `db:"shop_id"`
		Title        string `db:"title"`
		CoverURL     string `db:"cover_url"`
		LikeCount    uint64 `db:"like_count"`
		ReadCount    uint64 `db:"read_count"`
		CollectCount uint64 `db:"collect_count"`
		Status       string `db:"status"`
		AuditStatus  string `db:"audit_status"`
	}
	var rows []row
	err = r.conn.QueryRowsCtx(ctx, &rows, `
SELECT e.article_id, DATE_FORMAT(e.created_at, '%Y-%m-%d %H:%i:%s') AS engaged_at,
       a.id, a.shop_id, COALESCE(a.title,'') AS title, COALESCE(a.cover_url,'') AS cover_url,
       COALESCE(a.like_count,0) AS like_count, COALESCE(a.read_count,0) AS read_count,
       COALESCE(a.collect_count,0) AS collect_count,
       COALESCE(a.status,'') AS status, COALESCE(a.audit_status,'') AS audit_status
FROM `+table+` e
LEFT JOIN community_article a ON a.id = e.article_id
WHERE e.user_id = ?
ORDER BY e.id DESC
LIMIT ? OFFSET ?`, userID, pageSize, (page-1)*pageSize)
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
	n, _ := countCtx(ctx, r.conn,
		`SELECT COUNT(*) FROM homepage_slot_orders WHERE slot_type='article' AND status='active' AND target_id=? AND start_at<=? AND end_at>?`,
		articleID, now, now,
	)
	return n > 0
}

func (r *ArticleRepository) GetHomeArticleLimit(ctx context.Context) int {
	var lim int
	_ = r.conn.QueryRowCtx(ctx, &lim, "SELECT home_limit FROM homepage_slot_settings WHERE slot_type=? LIMIT 1", "article")
	if lim < 1 {
		return 6
	}
	return lim
}
