package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
)

type ArticleMerchantHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ArticleLogic
}

func NewArticleMerchantHandler(svcCtx *svc.ServiceContext) *ArticleMerchantHandler {
	return &ArticleMerchantHandler{svcCtx: svcCtx, logic: logic.NewArticleLogic(svcCtx)}
}

func (h *ArticleMerchantHandler) shopUser(r *http.Request) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(r.Context())
	userID, _ = middleware.GetUserID(r.Context())
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ArticleMerchantHandler) requirePerm(w http.ResponseWriter, r *http.Request, code string) bool {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return false
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, code) {
		response.Error(w, "无权限: "+code, http.StatusForbidden)
		return false
	}
	return true
}

func (h *ArticleMerchantHandler) List(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
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
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleMerchantHandler) Detail(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	data, err := h.logic.Detail(id, shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleMerchantHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:add") {
		return
	}
	shopID, uid, _ := h.shopUser(r)
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	a, err := h.logic.MerchantCreate(shopID, uid, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, a, "ok")
}

func (h *ArticleMerchantHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.MerchantUpdate(shopID, id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleMerchantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:delete") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.MerchantDelete(shopID, id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleMerchantHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.shopUser(r); !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	tree, err := h.logic.CategoryTree()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, tree, "ok")
}

func (h *ArticleMerchantHandler) CommentList(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	page, pageSize := middleware.ParsePage(r)
	articleID, _ := strconv.ParseUint(r.URL.Query().Get("article_id"), 10, 64)
	data, err := h.logic.ListComments(repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: r.URL.Query().Get("status"),
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleMerchantHandler) CommentPatch(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleCommentPatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.PatchComment(id, shopID, req.Status); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleMerchantHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "article:edit") {
		return
	}
	shopID, _, _ := h.shopUser(r)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteComment(id, shopID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleMerchantHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "article:edit") &&
		!h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "article:add") {
		response.Error(w, "无权限: article:edit", http.StatusForbidden)
		return
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		response.Error(w, "缺少文件", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, "读取失败", http.StatusBadRequest)
		return
	}
	url, err := h.logic.SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, map[string]string{"url": url}, "ok")
}
