package public

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"
)

func (h *ArticleHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	home := in.QueryGet("home") == "1"
	data, err := h.logic.PublicList(ctx, page, pageSize, home)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 && in.Request != nil {
		if raw := in.Request.Header.Get(middleware.GatewayUserIDHeader); raw != "" {
			userID, _ = strconv.ParseUint(raw, 10, 64)
		}
	}
	data, err := h.logic.PublicDetail(ctx, id, userID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) Like(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := h.logic.LikeArticle(ctx, userID, id, true); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Unlike(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := h.logic.LikeArticle(ctx, userID, id, false); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Favorite(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := h.logic.FavoriteArticle(ctx, userID, id, true); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Unfavorite(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := h.logic.FavoriteArticle(ctx, userID, id, false); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Status(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	liked, favorited := h.logic.EngagementStatus(ctx, userID, id)
	return map[string]bool{"liked": liked, "favorited": favorited}, nil
}

func (h *ArticleHandler) ListMyFavorites(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := h.logic.ListMyFavorites(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) ListMyLikes(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := h.logic.ListMyLikes(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) ListComments(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := h.logic.PublicListComments(ctx, id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) CreateComment(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	var req logic.CreateCommentReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := h.logic.CreatePublicComment(ctx, userID, id, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *ArticleHandler) CreateMine(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	var req struct {
		CategoryID uint64   `json:"category_id"`
		Title      string   `json:"title"`
		CoverURL   string   `json:"cover_url"`
		Content    string   `json:"content"`
		ImageURLs  []string `json:"image_urls"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	a, err := h.logic.UserCreate(ctx, userID, types.ArticleSaveReq{
		CategoryID: req.CategoryID, Title: req.Title, CoverURL: req.CoverURL,
		Content: req.Content, ImageURLs: req.ImageURLs,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return a, nil
}

func (h *ArticleHandler) ListMine(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	data, err := h.logic.UserListMine(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) DetailMine(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	data, err := h.logic.UserGetMine(ctx, userID, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) UpdateMine(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	var req struct {
		CategoryID uint64   `json:"category_id"`
		Title      string   `json:"title"`
		CoverURL   string   `json:"cover_url"`
		Content    string   `json:"content"`
		ImageURLs  []string `json:"image_urls"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UserUpdate(ctx, userID, id, types.ArticleSaveReq{
		CategoryID: req.CategoryID, Title: req.Title, CoverURL: req.CoverURL,
		Content: req.Content, ImageURLs: req.ImageURLs,
	}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) DeleteMine(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := h.logic.UserDelete(ctx, userID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) UploadMine(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	_ = userID
	if err := in.Request.ParseMultipartForm(6 << 20); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "上传失败")
	}
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "请选择文件")
	}
	defer file.Close()
	buf, err := io.ReadAll(io.LimitReader(file, 5<<20+1))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取文件失败")
	}
	url, err := h.logic.SaveUpload(0, hdr.Filename, buf)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}

func (h *ArticleHandler) ListEmojis(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListEmojisPublic(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}

func (h *ArticleHandler) ListBanners(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.PublicBanners()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}
