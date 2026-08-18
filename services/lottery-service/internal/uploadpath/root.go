package uploadpath

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
	root string
)

// Root 返回可写上传根目录（与 catalog/user 共用仓库 uploads，便于 /uploads 统一托管）。
func Root() string {
	once.Do(func() {
		if v := os.Getenv("UPLOAD_ROOT"); v != "" {
			_ = os.MkdirAll(v, 0o755)
			root = v
			return
		}
		candidates := []string{
			filepath.Join("..", "..", "uploads"),
			"uploads",
		}
		for _, c := range candidates {
			if err := os.MkdirAll(c, 0o755); err != nil {
				continue
			}
			probe := filepath.Join(c, ".writetest")
			if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
				continue
			}
			_ = os.Remove(probe)
			root = c
			return
		}
		root = "uploads"
	})
	return root
}

func Abs(elem ...string) string {
	parts := append([]string{Root()}, elem...)
	return filepath.Join(parts...)
}
