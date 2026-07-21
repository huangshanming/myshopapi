package upload

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mymall/services/catalog-service/internal/uploadpath"
)

// SaveShopImage stores a shop/platform image under uploads/shops.
func SaveShopImage(shopID uint64, filename string, data []byte) (string, error) {
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
