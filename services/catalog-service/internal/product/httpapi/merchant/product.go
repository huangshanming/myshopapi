package merchant

import (
	"encoding/json"
	"io"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *ProductHandler) shopUser(r *http.Request) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(r.Context())
	userID, _ = middleware.GetUserID(r.Context())
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ProductHandler) requirePerm(w http.ResponseWriter, r *http.Request, code string) bool {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return false
	}
	// JWT 店主首次访问时自动建店主角色
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, code) {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "无权限: "+code))
		return false
	}
	return true
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	catID, _ := strconv.ParseUint(r.URL.Query().Get("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: r.URL.Query().Get("name"), ProductNo: r.URL.Query().Get("product_no"),
		CategoryID: catID, Status: r.URL.Query().Get("status"), ProductType: r.URL.Query().Get("product_type"),
		StockWarnOnly: r.URL.Query().Get("stock_warn") == "1",
		Page:          page, PageSize: pageSize, OrderBy: r.URL.Query().Get("order_by"),
		Recycle: r.URL.Query().Get("recycle") == "1",
	}
	data, err := h.logic.List(f)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ProductHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "product:add") && !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "product:edit") {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "无权限: product:add"))
		return
	}
	var req types.MerchantProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.Save(shopID, uid, 0, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:edit") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.MerchantProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.Save(shopID, uid, id, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *ProductHandler) Copy(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	p, err := h.logic.Copy(shopID, uid, id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *ProductHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:status") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var body types.SetStatusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.SetStatus(shopID, uid, id, body.Status); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) Batch(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:batch") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.BatchProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	job, err := h.logic.Batch(shopID, uid, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, job)
}

func (h *ProductHandler) JobStatus(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	job, err := h.logic.Job(shopID, id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "任务不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, job)
}

func (h *ProductHandler) RecycleRestore(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.RecycleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.Restore(shopID, uid, req.ProductIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) RecycleDelete(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.RecycleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.PermanentDelete(shopID, uid, req.ProductIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.StockAdjustReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	req.SkuID = id
	if err := h.logic.AdjustStock(shopID, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) BatchStock(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.BatchStockReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.BatchStock(shopID, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) StockWarnings(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.StockWarnings(shopID, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ProductHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "上传失败"))
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

func (h *ProductHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ScheduleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.CreateSchedule(shopID, uid, id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) CancelSchedule(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.CancelSchedule(id, shopID)
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) OpLogs(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	pid, _ := strconv.ParseUint(r.URL.Query().Get("product_id"), 10, 64)
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.OpLogs(shopID, pid, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *ProductHandler) Export(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	url, err := h.logic.ExportCSV(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]string{"url": url})
}

func (h *ProductHandler) Import(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	_ = r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "缺少文件"))
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	res, err := h.logic.ImportCSV(shopID, uid, string(data))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, res)
}

func (h *ProductHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	list, err := h.svcCtx.ProductAdmin.ListTags(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *ProductHandler) SaveTag(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.TagReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	tag := &model.ProductTag{ID: id, ShopID: shopID, Name: req.Name, Color: req.Color, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveTag(tag); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, tag)
}

func (h *ProductHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteTag(id, shopID)
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *ProductHandler) ListAttrTemplates(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	list, err := h.svcCtx.ProductAdmin.ListAttrTemplates(shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *ProductHandler) SaveAttrTemplate(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	var req types.AttrTemplateReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	t := &model.ProductAttrTemplate{ID: id, ShopID: shopID, Name: req.Name, AttrsJSON: req.AttrsJSON, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveAttrTemplate(t); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, t)
}

func (h *ProductHandler) DeleteAttrTemplate(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteAttrTemplate(id, shopID)
	httpx.OkJsonCtx(r.Context(), w, nil)
}
