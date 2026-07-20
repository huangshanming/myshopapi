package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type PointsProductAdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsProductLogic
}

func NewPointsProductAdminHandler(svcCtx *svc.ServiceContext) *PointsProductAdminHandler {
	return &PointsProductAdminHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsProductLogic(context.Background(), svcCtx),
	}
}

func (h *PointsProductAdminHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	list, total, err := h.logic.List(page, pageSize, r.URL.Query().Get("status"), r.URL.Query().Get("keyword"))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *PointsProductAdminHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	p, err := h.logic.Get(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *PointsProductAdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req logic.PointsProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.Create(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *PointsProductAdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	var req logic.PointsProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.Update(id, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *PointsProductAdminHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.SetStatus(id, req.Status)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *PointsProductAdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "商品ID无效"))
		return
	}
	if err := h.logic.Delete(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *PointsProductAdminHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "上传失败"))
		return
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "请选择文件"))
		return
	}
	defer file.Close()
	buf, err := io.ReadAll(io.LimitReader(file, 5<<20+1))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "读取文件失败"))
		return
	}
	url, err := h.logic.SaveUpload(hdr.Filename, buf)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]string{"url": url})
}
