package logic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/common"
	"mymall/services/catalog-service/internal/model"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
	"mymall/services/catalog-service/internal/uploadpath"
)

type ArticleLogic struct {
	svcCtx *svc.ServiceContext
}

func NewArticleLogic(svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{svcCtx: svcCtx}
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

// AdminCreate 管理员创建：跳过审核
func (l *ArticleLogic) AdminCreate(operatorID uint64, req types.ArticleSaveReq) (*model.CommunityArticle, error) {
	if req.ShopID == 0 {
		return nil, errors.New("请选择商家")
	}
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
		"title":               strings.TrimSpace(req.Title),
		"cover_url":           req.CoverURL,
		"content":             req.Content,
		"category_id":         req.CategoryID,
		"schedule_publish_at": scheduleAt,
	}
	if req.ShopID > 0 {
		updates["shop_id"] = req.ShopID
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
	shopID := a.ShopID
	if req.ShopID > 0 {
		shopID = req.ShopID
	}
	if req.ImageURLs != nil {
		_ = l.svcCtx.Articles.ReplaceImages(id, shopID, req.ImageURLs)
	}
	return nil
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
		if strings.TrimSpace(req.RejectReason) == "" {
			return errors.New("请填写驳回理由")
		}
		return l.svcCtx.Articles.Update(id, map[string]interface{}{
			"audit_status":  model.ArticleAuditRejected,
			"reject_reason": strings.TrimSpace(req.RejectReason),
			"status":        model.ArticleDraft,
			"published_at":  nil,
		})
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

func (l *ArticleLogic) Offline(id uint64) error {
	a, err := l.svcCtx.Articles.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if a.Status == model.ArticleDeleted {
		return errors.New("已删除")
	}
	return l.svcCtx.Articles.Update(id, map[string]interface{}{"status": model.ArticleOffline})
}

func (l *ArticleLogic) SoftDelete(id uint64) error {
	return l.svcCtx.Articles.SoftDelete(id)
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
	if shopID == 0 {
		return "", errors.New("缺少店铺")
	}
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	dir := uploadpath.Abs("articles", fmt.Sprintf("%d", shopID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/articles/" + fmt.Sprintf("%d", shopID) + "/" + name, nil
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
