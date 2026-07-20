package merchant

import (
	"encoding/json"
	"io"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *ArticleHandler) shopUser(r *http.Request) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(r.Context())
	userID, _ = middleware.GetUserID(r.Context())
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ArticleHandler) requirePerm(w http.ResponseWriter, r *http.Request, code string) bool {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return false
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, code) {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "无权限: "+code))
		return false
	}
	return true
}

func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.List(repository.ArticleListFilter{
		ShopID: shopID, Title: r.URL.Query().Get("title"),
		AuditStatus: r.URL.Query().Get("audit_status"),
		Status:      r.URL.Query().Get("status"),
		Page:        page, PageSize: pageSize,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ArticleHandler) Detail(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	data, err := h.logic.Detail(id, shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:add") {
		return
	}
	shopID, uid, _ := h.shopUser(r)
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	a, err := h.logic.MerchantCreate(shopID, uid, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, a)
}

func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.MerchantUpdate(shopID, id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:delete") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.MerchantDelete(shopID, id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ArticleHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.shopUser(r); !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	tree, err := h.logic.CategoryTree()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, tree)
}

func (h *ArticleHandler) CommentList(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	articleID, _ := strconv.ParseUint(r.URL.Query().Get("article_id"), 10, 64)
	data, err := h.logic.ListComments(repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: r.URL.Query().Get("status"),
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ArticleHandler) CommentPatch(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleCommentPatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.PatchComment(id, shopID, req.Status); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ArticleHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteComment(id, shopID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ArticleHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "article:edit") &&
		!h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "article:add") {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "无权限: article:edit"))
		return
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "缺少文件"))
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "读取失败"))
		return
	}
	url, err := h.logic.SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]string{"url": url})
}
