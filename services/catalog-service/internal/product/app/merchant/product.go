package merchant

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/product/model"
	"mymall/services/catalog-service/internal/product/repository"
	"mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"
)

func (h *ProductHandler) shopUser(ctx context.Context) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(ctx)
	userID, _ = middleware.GetUserID(ctx)
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ProductHandler) requirePerm(ctx context.Context, code string) error {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	// JWT 店主首次访问时自动建店主角色
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, code) {
		return xerr.New(http.StatusForbidden, "无权限: "+code)
	}
	return nil
}

func (h *ProductHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	catID, _ := strconv.ParseUint(in.QueryGet("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: in.QueryGet("name"), ProductNo: in.QueryGet("product_no"),
		CategoryID: catID, Status: in.QueryGet("status"), ProductType: in.QueryGet("product_type"),
		StockWarnOnly: in.QueryGet("stock_warn") == "1",
		Page:          page, PageSize: pageSize, OrderBy: in.QueryGet("order_by"),
		Recycle: in.QueryGet("recycle") == "1",
	}
	data, err := h.logic.List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ProductHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	data, err := h.logic.Detail(ctx, id, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return data, nil
}

func (h *ProductHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "product:add") && !h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "product:edit") {
		return nil, xerr.New(http.StatusForbidden, "无权限: product:add")
	}
	var req types.MerchantProductSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.Save(ctx, shopID, uid, 0, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *ProductHandler) Update(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "product:edit"); err != nil {
		return nil, err
	}
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.MerchantProductSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.Save(ctx, shopID, uid, id, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *ProductHandler) Copy(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	p, err := h.logic.Copy(ctx, shopID, uid, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *ProductHandler) SetStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "product:status"); err != nil {
		return nil, err
	}
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body types.SetStatusReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SetStatus(ctx, shopID, uid, id, body.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) Batch(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "product:batch"); err != nil {
		return nil, err
	}
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.BatchProductReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	job, err := h.logic.Batch(ctx, shopID, uid, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return job, nil
}

func (h *ProductHandler) JobStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	job, err := h.logic.Job(ctx, shopID, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "任务不存在")
	}
	return job, nil
}

func (h *ProductHandler) RecycleRestore(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.RecycleReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.Restore(shopID, uid, req.ProductIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) RecycleDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.RecycleReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.PermanentDelete(ctx, shopID, uid, req.ProductIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) AdjustStock(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.StockAdjustReq
	_ = appinput.BindBody(in, &req)
	req.SkuID = id
	if err := h.logic.AdjustStock(ctx, shopID, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) BatchStock(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.BatchStockReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.BatchStock(ctx, shopID, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) StockWarnings(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	data, err := h.logic.StockWarnings(ctx, shopID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ProductHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if err := in.Request.ParseMultipartForm(6 << 20); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "上传失败")
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
	url, err := h.logic.SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}

func (h *ProductHandler) Schedule(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ScheduleReq
	_ = appinput.BindBody(in, &req)
	if err := h.logic.CreateSchedule(ctx, shopID, uid, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ProductHandler) CancelSchedule(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.CancelSchedule(ctx, id, shopID)
	return nil, nil
}

func (h *ProductHandler) OpLogs(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	pid, _ := strconv.ParseUint(in.QueryGet("product_id"), 10, 64)
	page, pageSize := in.Page()
	data, err := h.logic.OpLogs(ctx, shopID, pid, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ProductHandler) Export(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	url, err := h.logic.ExportCSV(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]string{"url": url}, nil
}

func (h *ProductHandler) Import(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	_ = in.Request.ParseMultipartForm(10 << 20)
	file, _, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	res, err := h.logic.ImportCSV(shopID, uid, string(data))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return res, nil
}

func (h *ProductHandler) ListTags(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	list, err := h.svcCtx.ProductAdmin.ListTags(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *ProductHandler) SaveTag(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.TagReq
	_ = appinput.BindBody(in, &req)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	tag := &model.ProductTag{ID: id, ShopID: shopID, Name: req.Name, Color: req.Color, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveTag(ctx, tag); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return tag, nil
}

func (h *ProductHandler) DeleteTag(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteTag(ctx, id, shopID)
	return nil, nil
}

func (h *ProductHandler) ListAttrTemplates(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	list, err := h.svcCtx.ProductAdmin.ListAttrTemplates(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *ProductHandler) SaveAttrTemplate(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var req types.AttrTemplateReq
	_ = appinput.BindBody(in, &req)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	t := &model.ProductAttrTemplate{ID: id, ShopID: shopID, Name: req.Name, AttrsJSON: req.AttrsJSON, Status: 1}
	if err := h.svcCtx.ProductAdmin.SaveAttrTemplate(ctx, t); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return t, nil
}

func (h *ProductHandler) DeleteAttrTemplate(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	_ = h.svcCtx.ProductAdmin.DeleteAttrTemplate(ctx, id, shopID)
	return nil, nil
}
