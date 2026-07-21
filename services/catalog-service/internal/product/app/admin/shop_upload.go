package admin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mymall/pkg/appinput"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/uploadpath"
)

func (h *ShopUploadHandler) Upload(ctx context.Context, in appinput.CallInput) (any, error) {
	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := saveShopUpload(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]string{"url": url}, nil
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
