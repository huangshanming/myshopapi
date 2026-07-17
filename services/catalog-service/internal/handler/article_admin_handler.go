package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/logic"
	"mymall/services/catalog-service/internal/repository"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"
)

type ArticleAdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ArticleLogic
}

func NewArticleAdminHandler(svcCtx *svc.ServiceContext) *ArticleAdminHandler {
	return &ArticleAdminHandler{svcCtx: svcCtx, logic: logic.NewArticleLogic(svcCtx)}
}

func (h *ArticleAdminHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	f := repository.ArticleListFilter{
		ShopID: shopID, Title: r.URL.Query().Get("title"),
		AuditStatus: r.URL.Query().Get("audit_status"),
		Status:      r.URL.Query().Get("status"),
		Page:        page, PageSize: pageSize,
		Recycle: r.URL.Query().Get("recycle") == "1" || strings.Contains(r.URL.Path, "/recycle"),
	}
	if s := r.URL.Query().Get("has_schedule"); s == "1" {
		v := true
		f.HasSchedule = &v
	} else if s == "0" {
		v := false
		f.HasSchedule = &v
	}
	if s := r.URL.Query().Get("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := r.URL.Query().Get("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	data, err := h.logic.List(f)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleAdminHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	data, err := h.logic.Detail(id, 0)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleAdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	a, err := h.logic.AdminCreate(uid, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, a, "ok")
}

func (h *ArticleAdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.GetUserID(r.Context())
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.AdminUpdate(id, uid, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) Audit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleAuditReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.Audit(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) BatchAudit(w http.ResponseWriter, r *http.Request) {
	var req types.ArticleBatchAuditReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.BatchAudit(req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) Top(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleTopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SetTop(id, req.IsTop); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) Offline(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.Offline(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.SoftDelete(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) RecycleRestore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID uint64 `json:"id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.ID == 0 {
		response.Error(w, "缺少 id", http.StatusBadRequest)
		return
	}
	if err := h.logic.Restore(body.ID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) RecycleDelete(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID uint64 `json:"id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.ID == 0 {
		response.Error(w, "缺少 id", http.StatusBadRequest)
		return
	}
	if err := h.logic.PermanentDelete(body.ID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.Stats()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ArticleAdminHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	tree, err := h.logic.CategoryTree()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, tree, "ok")
}

func (h *ArticleAdminHandler) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	var req types.ArticleCategorySaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SaveCategory(0, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleCategorySaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SaveCategory(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteCategory(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) CommentList(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	articleID, _ := strconv.ParseUint(r.URL.Query().Get("article_id"), 10, 64)
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
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

func (h *ArticleAdminHandler) CommentPatch(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ArticleCommentPatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.PatchComment(id, 0, req.Status); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.DeleteComment(id, 0); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ArticleAdminHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	if shopID == 0 {
		response.Error(w, "请传 shop_id", http.StatusBadRequest)
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
