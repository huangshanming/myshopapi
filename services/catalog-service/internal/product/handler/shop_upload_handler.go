package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mymall/services/catalog-service/internal/uploadpath"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

// ShopUploadHandler 平台后台店铺图片上传（存本地 /uploads/shops）
type ShopUploadHandler struct{}

func NewShopUploadHandler() *ShopUploadHandler {
	return &ShopUploadHandler{}
}

func (h *ShopUploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
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
	url, err := saveShopUpload(shopID, hdr.Filename, data)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]string{"url": url})
}

func saveShopUpload(shopID uint64, filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	owner := "platform"
	if shopID > 0 {
		owner = fmt.Sprintf("%d", shopID)
	}
	dir := uploadpath.Abs("shops", owner)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/shops/" + owner + "/" + name, nil
}
