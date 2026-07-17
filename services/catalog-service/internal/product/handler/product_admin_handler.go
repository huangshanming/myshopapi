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
	"mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"mymall/services/catalog-service/internal/svc"
)

// ProductAdminHandler 商家商品中台
type ProductAdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ProductAdminLogic
}

func NewProductAdminHandler(svcCtx *svc.ServiceContext) *ProductAdminHandler {
	return &ProductAdminHandler{svcCtx: svcCtx, logic: logic.NewProductAdminLogic(svcCtx)}
}

func (h *ProductAdminHandler) shopUser(r *http.Request) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(r.Context())
	userID, _ = middleware.GetUserID(r.Context())
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ProductAdminHandler) requirePerm(w http.ResponseWriter, r *http.Request, code string) bool {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return false
	}
	// JWT 店主首次访问时自动建店主角色
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, code) {
		response.Error(w, "无权限: "+code, http.StatusForbidden)
		return false
	}
	return true
}

func (h *ProductAdminHandler) List(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
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
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ProductAdminHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductAdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	if middleware.GetUserRole(r.Context()) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "product:add") && !h.svcCtx.ShopRBAC.HasPerm(shopID, uid, "product:edit") {
		response.Error(w, "无权限: product:add", http.StatusForbidden)
		return
	}
	var req types.MerchantProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	p, err := h.logic.Save(shopID, uid, 0, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, p, "创建成功")
}

func (h *ProductAdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:edit") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.MerchantProductSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	p, err := h.logic.Save(shopID, uid, id, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, p, "更新成功")
}

func (h *ProductAdminHandler) Copy(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	p, err := h.logic.Copy(shopID, uid, id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, p, "已复制为草稿")
}

func (h *ProductAdminHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:status") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var body types.SetStatusReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.SetStatus(shopID, uid, id, body.Status); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ProductAdminHandler) Batch(w http.ResponseWriter, r *http.Request) {
	if !h.requirePerm(w, r, "product:batch") {
		return
	}
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.BatchProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	job, err := h.logic.Batch(shopID, uid, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, job, "任务已提交")
}

func (h *ProductAdminHandler) JobStatus(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	job, err := h.logic.Job(shopID, id)
	if err != nil {
		response.Error(w, "任务不存在", http.StatusNotFound)
		return
	}
	response.Success(w, job, "ok")
}

func (h *ProductAdminHandler) RecycleRestore(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.RecycleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.Restore(shopID, uid, req.ProductIDs); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已恢复")
}

func (h *ProductAdminHandler) RecycleDelete(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.RecycleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.PermanentDelete(shopID, uid, req.ProductIDs); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已永久删除")
}

func (h *ProductAdminHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.StockAdjustReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	req.SkuID = id
	if err := h.logic.AdjustStock(shopID, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ProductAdminHandler) BatchStock(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.BatchStockReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.BatchStock(shopID, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *ProductAdminHandler) StockWarnings(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.StockWarnings(shopID, page, pageSize)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ProductAdminHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		response.Error(w, "上传失败", http.StatusBadRequest)
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
	response.Success(w, map[string]string{"url": url}, "上传成功")
}

func (h *ProductAdminHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req types.ScheduleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.CreateSchedule(shopID, uid, id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已设置")
}

func (h *ProductAdminHandler) CancelSchedule(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.CancelSchedule(id, shopID)
	response.Success(w, nil, "已取消")
}

func (h *ProductAdminHandler) OpLogs(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	pid, _ := strconv.ParseUint(r.URL.Query().Get("product_id"), 10, 64)
	page, pageSize := middleware.ParsePage(r)
	data, err := h.logic.OpLogs(shopID, pid, page, pageSize)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *ProductAdminHandler) Export(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	url, err := h.logic.ExportCSV(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, map[string]string{"url": url}, "ok")
}

func (h *ProductAdminHandler) Import(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	_ = r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		response.Error(w, "缺少文件", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	res, err := h.logic.ImportCSV(shopID, uid, string(data))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, res, "导入完成")
}

func (h *ProductAdminHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	list, err := h.svcCtx.ProductAdmin.ListTags(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *ProductAdminHandler) SaveTag(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.TagReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	tag := &model.ProductTag{ID: id, ShopID: shopID, Name: req.Name, Color: req.Color, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveTag(tag); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, tag, "ok")
}

func (h *ProductAdminHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteTag(id, shopID)
	response.Success(w, nil, "ok")
}

func (h *ProductAdminHandler) ListAttrTemplates(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	list, err := h.svcCtx.ProductAdmin.ListAttrTemplates(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *ProductAdminHandler) SaveAttrTemplate(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	var req types.AttrTemplateReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	t := &model.ProductAttrTemplate{ID: id, ShopID: shopID, Name: req.Name, AttrsJSON: req.AttrsJSON, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveAttrTemplate(t); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, t, "ok")
}

func (h *ProductAdminHandler) DeleteAttrTemplate(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteAttrTemplate(id, shopID)
	response.Success(w, nil, "ok")
}
