package admin

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *ArticleHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	f := repository.ArticleListFilter{
		Title:       in.QueryGet("title"),
		AuditStatus: in.QueryGet("audit_status"),
		Status:      in.QueryGet("status"),
		Page:        page, PageSize: pageSize,
		Recycle: in.QueryGet("recycle") == "1" || (in.Request != nil && strings.Contains(in.Request.URL.Path, "/recycle")),
	}
	if s := in.QueryGet("shop_id"); s != "" {
		shopID, _ := strconv.ParseUint(s, 10, 64)
		f.ShopID = shopID
		f.FilterShop = true
	}
	if s := in.QueryGet("has_schedule"); s == "1" {
		v := true
		f.HasSchedule = &v
	} else if s == "0" {
		v := false
		f.HasSchedule = &v
	}
	if s := in.QueryGet("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := in.QueryGet("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	data, err := h.logic.List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	data, err := h.logic.Detail(ctx, id, 0)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	var req types.ArticleSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	a, err := h.logic.AdminCreate(ctx, uid, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return a, nil
}

func (h *ArticleHandler) Update(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AdminUpdate(ctx, id, uid, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Audit(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleAuditReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.Audit(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) BatchAudit(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.ArticleBatchAuditReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.BatchAudit(ctx, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Top(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleTopReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SetTop(ctx, id, req.IsTop); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Offline(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleRemarkReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.Offline(ctx, id, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) SoftDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleRemarkReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.SoftDelete(ctx, id, req.Remark); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) RecycleRestore(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		ID uint64 `json:"id"`
	}
	_ = appinput.BindBody(in, &body)
	if body.ID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少 id")
	}
	if err := h.logic.Restore(ctx, body.ID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) RecycleDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		ID uint64 `json:"id"`
	}
	_ = appinput.BindBody(in, &body)
	if body.ID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少 id")
	}
	if err := h.logic.PermanentDelete(ctx, body.ID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Stats(ctx context.Context, in appinput.CallInput) (any, error) {
	data, err := h.logic.Stats(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) CategoryList(ctx context.Context, in appinput.CallInput) (any, error) {
	tree, err := h.logic.CategoryTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return tree, nil
}

func (h *ArticleHandler) CategoryCreate(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.ArticleCategorySaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SaveCategory(ctx, 0, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CategoryUpdate(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleCategorySaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SaveCategory(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CategoryDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteCategory(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CommentList(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	articleID, _ := strconv.ParseUint(in.QueryGet("article_id"), 10, 64)
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	data, err := h.logic.ListComments(ctx, repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: in.QueryGet("status"),
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) CommentPatch(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleCommentPatchReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.PatchComment(ctx, id, 0, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CommentDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteComment(ctx, id, 0); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) EmojiList(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	data, err := h.logic.ListEmojisAdmin(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) EmojiCreate(ctx context.Context, in appinput.CallInput) (any, error) {
	var req struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Sort     int    `json:"sort"`
		Status   *int8  `json:"status"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	e, err := h.logic.CreateEmoji(ctx, req.Name, req.ImageURL, req.Sort, status)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return e, nil
}

func (h *ArticleHandler) EmojiUpdate(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
		Sort     *int   `json:"sort"`
		Status   *int8  `json:"status"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateEmoji(ctx, id, req.Name, req.ImageURL, req.Sort, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) EmojiDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteEmoji(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := h.logic.SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}

func (h *ArticleHandler) ListBanners(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := h.logic.AdminListBanners(page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) GetBanner(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	b, err := h.logic.AdminGetBanner(id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return b, nil
}

func (h *ArticleHandler) CreateBanner(ctx context.Context, in appinput.CallInput) (any, error) {
	var req logic.BannerSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	b, err := h.logic.AdminCreateBanner(req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return b, nil
}

func (h *ArticleHandler) UpdateBanner(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var req logic.BannerSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AdminUpdateBanner(id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) DeleteBanner(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	if err := h.logic.AdminDeleteBanner(id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) UploadBanner(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := h.logic.SaveBannerUpload(hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}
