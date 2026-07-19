package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/client/userhttp"
	"mymall/services/catalog-service/internal/content/model"
	"mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/content/types"
	notifymodel "mymall/services/catalog-service/internal/notify/model"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/uploadpath"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func parseScheduleAt(s string) (*common.LocalTime, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return nil, errors.New("定时发布时间格式应为 2006-01-02 15:04:05")
	}
	lt := common.LocalTime(t)
	return &lt, nil
}

func (l *ArticleLogic) resolvePublishStatus(scheduleAt *common.LocalTime) (status string, publishedAt *common.LocalTime) {
	now := time.Now()
	if scheduleAt != nil && time.Time(*scheduleAt).After(now) {
		return model.ArticleScheduled, nil
	}
	pub := common.LocalTime(now)
	return model.ArticlePublished, &pub
}

func (l *ArticleLogic) List(f repository.ArticleListFilter) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.List(f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *ArticleLogic) Detail(id, shopID uint64) (map[string]interface{}, error) {
	var (
		a   *model.CommunityArticle
		err error
	)
	if shopID > 0 {
		a, err = l.svcCtx.Articles.GetByIDShop(id, shopID)
	} else {
		a, err = l.svcCtx.Articles.GetByID(id)
	}
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	imgs, _ := l.svcCtx.Articles.ListImages(id)
	return map[string]interface{}{"article": a, "images": imgs}, nil
}

// AdminCreate 管理员创建：跳过审核；shop_id=0 表示平台官方文章
func (l *ArticleLogic) AdminCreate(operatorID uint64, req types.ArticleSaveReq) (*model.CommunityArticle, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("标题不能为空")
	}
	scheduleAt, err := parseScheduleAt(req.SchedulePublishAt)
	if err != nil {
		return nil, err
	}
	status, publishedAt := l.resolvePublishStatus(scheduleAt)
	var isTop int8
	if req.IsTop != nil {
		isTop = *req.IsTop
	}
	a := &model.CommunityArticle{
		ShopID:            req.ShopID,
		CategoryID:        req.CategoryID,
		Title:             strings.TrimSpace(req.Title),
		CoverURL:          req.CoverURL,
		Content:           req.Content,
		AuditStatus:       model.ArticleAuditApproved,
		Status:            status,
		SchedulePublishAt: scheduleAt,
		IsTop:             isTop,
		PublishedAt:       publishedAt,
		CreatedBy:         operatorID,
	}
	if err := l.svcCtx.Articles.Create(a); err != nil {
		return nil, err
	}
	_ = l.svcCtx.Articles.ReplaceImages(a.ID, a.ShopID, req.ImageURLs)
	return a, nil
}

// AdminUpdate 管理员更新
func (l *ArticleLogic) AdminUpdate(id, operatorID uint64, req types.ArticleSaveReq) error {
	a, err := l.svcCtx.Articles.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已删除文章不可编辑")
	}
	scheduleAt, err := parseScheduleAt(req.SchedulePublishAt)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"shop_id":             req.ShopID,
		"title":               strings.TrimSpace(req.Title),
		"cover_url":           req.CoverURL,
		"content":             req.Content,
		"category_id":         req.CategoryID,
		"schedule_publish_at": scheduleAt,
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	// 管理员编辑后保持已审；按定时重算状态（已下架除外）
	if a.Status != model.ArticleOffline {
		st, pub := l.resolvePublishStatus(scheduleAt)
		updates["audit_status"] = model.ArticleAuditApproved
		updates["status"] = st
		updates["published_at"] = pub
		updates["reject_reason"] = ""
	}
	_ = operatorID
	if err := l.svcCtx.Articles.Update(id, updates); err != nil {
		return err
	}
	if req.ImageURLs != nil {
		_ = l.svcCtx.Articles.ReplaceImages(id, req.ShopID, req.ImageURLs)
	}
	return nil
}

func (l *ArticleLogic) notifyShop(shopID uint64, typ, title, content, link, refType string, refID uint64) {
	if shopID == 0 || l.svcCtx.Notifications == nil {
		return
	}
	_ = l.svcCtx.Notifications.Create(&notifymodel.ShopNotification{
		ShopID:  shopID,
		Type:    typ,
		Title:   title,
		Content: content,
		Link:    link,
		RefType: refType,
		RefID:   refID,
	})
}

func (l *ArticleLogic) Audit(id uint64, req types.ArticleAuditReq) error {
	a, err := l.svcCtx.Articles.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已删除文章不可审核")
	}
	if !req.Pass {
		reason := strings.TrimSpace(req.RejectReason)
		if reason == "" {
			return errors.New("请填写驳回理由")
		}
		if err := l.svcCtx.Articles.Update(id, map[string]interface{}{
			"audit_status":  model.ArticleAuditRejected,
			"reject_reason": reason,
			"status":        model.ArticleDraft,
			"published_at":  nil,
		}); err != nil {
			return err
		}
		l.notifyShop(a.ShopID, notifymodel.NotifArticleRejected, "文章审核未通过",
			fmt.Sprintf("文章「%s」(ID:%d) 审核未通过。理由：%s", a.Title, a.ID, reason),
			"/merchant/articles", "article", a.ID)
		return nil
	}
	st, pub := l.resolvePublishStatus(a.SchedulePublishAt)
	return l.svcCtx.Articles.Update(id, map[string]interface{}{
		"audit_status":  model.ArticleAuditApproved,
		"reject_reason": "",
		"status":        st,
		"published_at":  pub,
	})
}

func (l *ArticleLogic) BatchAudit(req types.ArticleBatchAuditReq) error {
	if len(req.IDs) == 0 {
		return errors.New("请选择文章")
	}
	for _, id := range req.IDs {
		if err := l.Audit(id, types.ArticleAuditReq{Pass: req.Pass, RejectReason: req.RejectReason}); err != nil {
			return fmt.Errorf("文章 %d: %w", id, err)
		}
	}
	return nil
}

func (l *ArticleLogic) SetTop(id uint64, isTop int8) error {
	return l.svcCtx.Articles.Update(id, map[string]interface{}{"is_top": isTop})
}

func (l *ArticleLogic) Offline(id uint64, remark string) error {
	remark = strings.TrimSpace(remark)
	a, err := l.svcCtx.Articles.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已删除")
	}
	if a.ShopID > 0 && remark == "" {
		return errors.New("请填写备注，将通知商家")
	}
	if err := l.svcCtx.Articles.Update(id, map[string]interface{}{"status": model.ArticleOffline}); err != nil {
		return err
	}
	if a.ShopID > 0 {
		l.notifyShop(a.ShopID, notifymodel.NotifArticleOffline, "文章被平台下架",
			fmt.Sprintf("文章「%s」(ID:%d) 已被平台下架。备注：%s", a.Title, a.ID, remark),
			"/merchant/articles", "article", a.ID)
	}
	return nil
}

func (l *ArticleLogic) SoftDelete(id uint64, remark string) error {
	remark = strings.TrimSpace(remark)
	a, err := l.svcCtx.Articles.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已在回收站")
	}
	if a.ShopID > 0 && remark == "" {
		return errors.New("请填写备注，将通知商家")
	}
	if err := l.svcCtx.Articles.SoftDelete(id); err != nil {
		return err
	}
	if a.ShopID > 0 {
		l.notifyShop(a.ShopID, notifymodel.NotifArticleDeleted, "文章被平台删除",
			fmt.Sprintf("文章「%s」(ID:%d) 已被平台移入回收站。备注：%s", a.Title, a.ID, remark),
			"/merchant/articles", "article", a.ID)
	}
	return nil
}

func (l *ArticleLogic) Restore(id uint64) error {
	return l.svcCtx.Articles.Restore(id)
}

func (l *ArticleLogic) PermanentDelete(id uint64) error {
	return l.svcCtx.Articles.PermanentDelete(id)
}

func (l *ArticleLogic) Stats() (map[string]interface{}, error) {
	return l.svcCtx.Articles.Stats()
}

// MerchantCreate 商家创建：强制待审
func (l *ArticleLogic) MerchantCreate(shopID, operatorID uint64, req types.ArticleSaveReq) (*model.CommunityArticle, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errors.New("标题不能为空")
	}
	scheduleAt, err := parseScheduleAt(req.SchedulePublishAt)
	if err != nil {
		return nil, err
	}
	a := &model.CommunityArticle{
		ShopID:            shopID,
		CategoryID:        req.CategoryID,
		Title:             strings.TrimSpace(req.Title),
		CoverURL:          req.CoverURL,
		Content:           req.Content,
		AuditStatus:       model.ArticleAuditPending,
		Status:            model.ArticleDraft,
		SchedulePublishAt: scheduleAt,
		IsTop:             0,
		CreatedBy:         operatorID,
	}
	if err := l.svcCtx.Articles.Create(a); err != nil {
		return nil, err
	}
	_ = l.svcCtx.Articles.ReplaceImages(a.ID, shopID, req.ImageURLs)
	return a, nil
}

// MerchantUpdate 仅 pending 可改；已审/驳回只读
func (l *ArticleLogic) MerchantUpdate(shopID, id uint64, req types.ArticleSaveReq) error {
	a, err := l.svcCtx.Articles.GetByIDShop(id, shopID)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已删除")
	}
	if a.AuditStatus != model.ArticleAuditPending {
		return errors.New("已审核或已驳回的文章不可编辑，请联系平台")
	}
	scheduleAt, err := parseScheduleAt(req.SchedulePublishAt)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"title":               strings.TrimSpace(req.Title),
		"cover_url":           req.CoverURL,
		"content":             req.Content,
		"category_id":         req.CategoryID,
		"schedule_publish_at": scheduleAt,
		"audit_status":        model.ArticleAuditPending,
		"status":              model.ArticleDraft,
		"reject_reason":       "",
	}
	if err := l.svcCtx.Articles.UpdateShop(id, shopID, updates); err != nil {
		return err
	}
	if req.ImageURLs != nil {
		_ = l.svcCtx.Articles.ReplaceImages(id, shopID, req.ImageURLs)
	}
	return nil
}

// MerchantDelete 仅待审可删
func (l *ArticleLogic) MerchantDelete(shopID, id uint64) error {
	a, err := l.svcCtx.Articles.GetByIDShop(id, shopID)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.AuditStatus != model.ArticleAuditPending {
		return errors.New("仅可删除待审核文章")
	}
	return l.svcCtx.Articles.SoftDeleteShop(id, shopID)
}

func (l *ArticleLogic) RunPublishSchedules() {
	_, _ = l.svcCtx.Articles.ClaimDuePublish(20)
}

func (l *ArticleLogic) SaveUpload(shopID uint64, filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	owner := "platform"
	if shopID > 0 {
		owner = fmt.Sprintf("%d", shopID)
	}
	dir := uploadpath.Abs("articles", owner)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/articles/" + owner + "/" + name, nil
}

func (l *ArticleLogic) CategoryTree() ([]map[string]interface{}, error) {
	list, err := l.svcCtx.Articles.ListCategories()
	if err != nil {
		return nil, err
	}
	byParent := map[uint64][]model.CommunityArticleCategory{}
	for _, c := range list {
		byParent[c.ParentID] = append(byParent[c.ParentID], c)
	}
	var build func(parent uint64) []map[string]interface{}
	build = func(parent uint64) []map[string]interface{} {
		nodes := byParent[parent]
		out := make([]map[string]interface{}, 0, len(nodes))
		for _, n := range nodes {
			out = append(out, map[string]interface{}{
				"id": n.ID, "parent_id": n.ParentID, "name": n.Name,
				"sort": n.Sort, "status": n.Status, "children": build(n.ID),
			})
		}
		return out
	}
	return build(0), nil
}

func (l *ArticleLogic) SaveCategory(id uint64, req types.ArticleCategorySaveReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("分类名不能为空")
	}
	if id == 0 {
		c := &model.CommunityArticleCategory{
			ParentID: req.ParentID, Name: strings.TrimSpace(req.Name), Sort: req.Sort, Status: 1,
		}
		if req.Status != nil {
			c.Status = *req.Status
		}
		return l.svcCtx.Articles.CreateCategory(c)
	}
	updates := map[string]interface{}{
		"parent_id": req.ParentID, "name": strings.TrimSpace(req.Name), "sort": req.Sort,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	return l.svcCtx.Articles.UpdateCategory(id, updates)
}

func (l *ArticleLogic) DeleteCategory(id uint64) error {
	return l.svcCtx.Articles.DeleteCategory(id)
}

func (l *ArticleLogic) ListComments(f repository.CommentListFilter) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.ListComments(f)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *ArticleLogic) PatchComment(id, shopID uint64, status string) error {
	switch status {
	case model.CommentVisible, model.CommentHidden, model.CommentDeleted:
	default:
		return errors.New("status 无效")
	}
	return l.svcCtx.Articles.PatchComment(id, shopID, status)
}

func (l *ArticleLogic) DeleteComment(id, shopID uint64) error {
	return l.svcCtx.Articles.DeleteComment(id, shopID)
}

func (l *ArticleLogic) PublicList(page, pageSize int, home bool) (map[string]interface{}, error) {
	homeLimit := 0
	if home {
		homeLimit = l.svcCtx.Articles.GetHomeArticleLimit()
	}
	list, total, err := l.svcCtx.Articles.ListPublic(page, pageSize, homeLimit)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(list))
	for _, a := range list {
		items = append(items, map[string]interface{}{
			"id": a.ID, "shop_id": a.ShopID, "title": a.Title, "cover_url": a.CoverURL,
			"like_count": a.LikeCount, "audience_count": a.AudienceCount, "read_count": a.ReadCount,
			"collect_count": a.CollectCount, "published_at": a.PublishedAt,
			"paid": l.svcCtx.Articles.IsArticleBoosted(a.ID),
		})
	}
	return map[string]interface{}{"list": items, "total": total}, nil
}

func (l *ArticleLogic) PublicDetail(id, userID uint64) (map[string]interface{}, error) {
	a, err := l.svcCtx.Articles.GetPublished(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	_ = l.svcCtx.Articles.RecordRead(id, userID)
	a, _ = l.svcCtx.Articles.GetPublished(id)
	liked, favorited := false, false
	if userID > 0 {
		liked, favorited = l.svcCtx.Articles.EngagementStatus(userID, id)
	}
	imgs, _ := l.svcCtx.Articles.ListImages(id)
	return map[string]interface{}{
		"article": a, "images": imgs,
		"liked": liked, "favorited": favorited,
		"paid": l.svcCtx.Articles.IsArticleBoosted(id),
	}, nil
}

func (l *ArticleLogic) LikeArticle(userID, articleID uint64, like bool) error {
	// 取消点赞允许文章已下架；新增点赞需文章仍在线
	if like {
		if _, err := l.svcCtx.Articles.GetPublished(articleID); err != nil {
			return errors.New("文章不存在")
		}
	}
	return l.svcCtx.Articles.ToggleLike(userID, articleID, like)
}

func (l *ArticleLogic) FavoriteArticle(userID, articleID uint64, fav bool) error {
	if fav {
		if _, err := l.svcCtx.Articles.GetPublished(articleID); err != nil {
			return errors.New("文章不存在")
		}
	}
	return l.svcCtx.Articles.ToggleFavorite(userID, articleID, fav)
}

func (l *ArticleLogic) EngagementStatus(userID, articleID uint64) (liked, favorited bool) {
	return l.svcCtx.Articles.EngagementStatus(userID, articleID)
}

func (l *ArticleLogic) ListMyFavorites(userID uint64, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.ListUserFavorites(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *ArticleLogic) ListMyLikes(userID uint64, page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.ListUserLikes(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

type CreateCommentReq struct {
	Content  string `json:"content"`
	ParentID uint64 `json:"parent_id"`
}

func (l *ArticleLogic) fillCommentUsers(list []model.CommunityArticleComment) {
	ids := make([]uint64, 0, len(list)*2)
	seen := map[uint64]struct{}{}
	add := func(id uint64) {
		if id == 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	for _, c := range list {
		add(c.UserID)
		add(c.ReplyToUserID)
		for _, ch := range c.Children {
			add(ch.UserID)
			add(ch.ReplyToUserID)
		}
	}
	m := l.svcCtx.Articles.MapUserBriefs(ids)
	nick := func(id uint64) string {
		u, ok := m[id]
		if !ok {
			return "用户"
		}
		if u.Nickname != "" {
			return u.Nickname
		}
		if len(u.Mobile) >= 4 {
			return "用户" + u.Mobile[len(u.Mobile)-4:]
		}
		return "用户"
	}
	for i := range list {
		list[i].UserNickname = nick(list[i].UserID)
		if list[i].ReplyToUserID > 0 {
			list[i].ReplyToNickname = nick(list[i].ReplyToUserID)
		}
		for j := range list[i].Children {
			list[i].Children[j].UserNickname = nick(list[i].Children[j].UserID)
			if list[i].Children[j].ReplyToUserID > 0 {
				list[i].Children[j].ReplyToNickname = nick(list[i].Children[j].ReplyToUserID)
			}
		}
	}
}

func (l *ArticleLogic) PublicListComments(articleID uint64, page, pageSize int) (map[string]interface{}, error) {
	if _, err := l.svcCtx.Articles.GetPublished(articleID); err != nil {
		return nil, errors.New("文章不存在")
	}
	roots, total, err := l.svcCtx.Articles.ListPublicCommentRoots(articleID, page, pageSize)
	if err != nil {
		return nil, err
	}
	rootIDs := make([]uint64, 0, len(roots))
	for _, c := range roots {
		rootIDs = append(rootIDs, c.ID)
	}
	children, _ := l.svcCtx.Articles.ListPublicCommentChildren(articleID, rootIDs)
	byRoot := map[uint64][]model.CommunityArticleComment{}
	for _, ch := range children {
		rid := ch.RootID
		if rid == 0 {
			rid = ch.ParentID
		}
		byRoot[rid] = append(byRoot[rid], ch)
	}
	for i := range roots {
		roots[i].Children = byRoot[roots[i].ID]
	}
	l.fillCommentUsers(roots)
	return map[string]interface{}{"list": roots, "total": total}, nil
}

func (l *ArticleLogic) CreatePublicComment(userID, articleID uint64, req CreateCommentReq) (*model.CommunityArticleComment, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("请输入评论内容")
	}
	if len([]rune(content)) > 500 {
		return nil, errors.New("评论最多 500 字")
	}
	a, err := l.svcCtx.Articles.GetPublished(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	c := &model.CommunityArticleComment{
		ArticleID: articleID,
		ShopID:    a.ShopID,
		UserID:    userID,
		Content:   content,
		Status:    model.CommentVisible,
	}
	var notifyUID uint64
	if req.ParentID > 0 {
		parent, err := l.svcCtx.Articles.GetComment(req.ParentID)
		if err != nil || parent.ArticleID != articleID {
			return nil, errors.New("回复的评论不存在")
		}
		c.ParentID = parent.ID
		if parent.RootID > 0 {
			c.RootID = parent.RootID
		} else {
			c.RootID = parent.ID
		}
		c.ReplyToUserID = parent.UserID
		notifyUID = parent.UserID
	}
	if err := l.svcCtx.Articles.CreateComment(c); err != nil {
		return nil, err
	}
	tmp := []model.CommunityArticleComment{*c}
	l.fillCommentUsers(tmp)
	*c = tmp[0]
	if notifyUID > 0 && notifyUID != userID && l.svcCtx.UserHTTP != nil {
		extra, _ := json.Marshal(map[string]interface{}{"comment_id": c.ID})
		preview := content
		if rs := []rune(preview); len(rs) > 40 {
			preview = string(rs[:40]) + "…"
		}
		_ = l.svcCtx.UserHTTP.Notify(l.ctx, userhttp.NotifyReq{
			UserID: notifyUID, Title: "收到新回复",
			Content:  fmt.Sprintf("%s 回复了你：%s", c.UserNickname, preview),
			MsgType:  "system",
			LinkType: "article",
			LinkID:   articleID,
			Extra:    string(extra),
		})
	}
	return c, nil
}

func (l *ArticleLogic) ListEmojisAdmin(page, pageSize int) (map[string]interface{}, error) {
	list, total, err := l.svcCtx.Articles.ListEmojisAdmin(page, pageSize)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (l *ArticleLogic) ListEmojisPublic() ([]model.CommunityCommentEmoji, error) {
	return l.svcCtx.Articles.ListEmojisPublic()
}

func (l *ArticleLogic) CreateEmoji(name, imageURL string, sort int, status int8) (*model.CommunityCommentEmoji, error) {
	name = strings.TrimSpace(name)
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return nil, errors.New("请上传表情图片")
	}
	if name == "" {
		name = "表情"
	}
	if status != 0 && status != 1 {
		status = 1
	}
	e := &model.CommunityCommentEmoji{Name: name, ImageURL: imageURL, Sort: sort, Status: status}
	if err := l.svcCtx.Articles.CreateEmoji(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (l *ArticleLogic) UpdateEmoji(id uint64, name, imageURL string, sort *int, status *int8) error {
	if _, err := l.svcCtx.Articles.GetEmoji(id); err != nil {
		return errors.New("表情不存在")
	}
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = strings.TrimSpace(name)
	}
	if imageURL != "" {
		updates["image_url"] = strings.TrimSpace(imageURL)
	}
	if sort != nil {
		updates["sort"] = *sort
	}
	if status != nil {
		updates["status"] = *status
	}
	if len(updates) == 0 {
		return nil
	}
	return l.svcCtx.Articles.UpdateEmoji(id, updates)
}

func (l *ArticleLogic) DeleteEmoji(id uint64) error {
	return l.svcCtx.Articles.DeleteEmoji(id)
}
