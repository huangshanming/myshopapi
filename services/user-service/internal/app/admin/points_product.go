package admin

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/services/user-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/pkg/xerr"
)

func (h *PointsProductHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	list, total, err := h.logic.List(ctx, page, pageSize, in.QueryGet("status"), in.QueryGet("keyword"))
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *PointsProductHandler) Detail(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	p, err := h.logic.Get(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *PointsProductHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.PointsProductSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.Create(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *PointsProductHandler) Update(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	var req biz.PointsProductSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.Update(ctx, id, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *PointsProductHandler) SetStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.SetStatus(ctx, id, req.Status)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *PointsProductHandler) Delete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "商品ID无效")
	}
	if err := h.logic.Delete(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *PointsProductHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

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
	url, err := h.logic.SaveUpload(ctx, hdr.Filename, buf)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
}
