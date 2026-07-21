package merchant

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"
)

func (h *ArticleHandler) shopUser(ctx context.Context) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(ctx)
	userID, _ = middleware.GetUserID(ctx)
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ArticleHandler) requirePerm(ctx context.Context, code string) error {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, code) {
		return xerr.New(http.StatusForbidden, "无权限: "+code)
	}
	return nil
}

func (h *ArticleHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	data, err := h.logic.List(ctx, repository.ArticleListFilter{
		ShopID: shopID, Title: in.QueryGet("title"),
		AuditStatus: in.QueryGet("audit_status"),
		Status:      in.QueryGet("status"),
		Page:        page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *ArticleHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
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

func (h *ArticleHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "article:add"); err != nil {
		return nil, err
	}
	shopID, uid, _ := h.shopUser(ctx)
	var req types.ArticleSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	a, err := h.logic.MerchantCreate(ctx, shopID, uid, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return a, nil
}

func (h *ArticleHandler) Update(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "article:edit"); err != nil {
		return nil, err
	}
	shopID, _, _ := h.shopUser(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.MerchantUpdate(ctx, shopID, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Delete(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "article:delete"); err != nil {
		return nil, err
	}
	shopID, _, _ := h.shopUser(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.MerchantDelete(ctx, shopID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CategoryList(ctx context.Context, in appinput.CallInput) (any, error) {
	if _, _, ok := h.shopUser(ctx); !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	tree, err := h.logic.CategoryTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return tree, nil
}

func (h *ArticleHandler) CommentList(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	articleID, _ := strconv.ParseUint(in.QueryGet("article_id"), 10, 64)
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
	if err := h.requirePerm(ctx, "article:edit"); err != nil {
		return nil, err
	}
	shopID, _, _ := h.shopUser(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.ArticleCommentPatchReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.PatchComment(ctx, id, shopID, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) CommentDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	if err := h.requirePerm(ctx, "article:edit"); err != nil {
		return nil, err
	}
	shopID, _, _ := h.shopUser(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteComment(ctx, id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ArticleHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "article:edit") &&
		!h.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "article:add") {
		return nil, xerr.New(http.StatusForbidden, "无权限: article:edit")
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
