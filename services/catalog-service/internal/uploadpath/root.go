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

// Root 返回可写的上传根目录。
// 优先 ./uploads；若被 root 占用不可写，则回退到仓库根 ../../uploads（从 services/catalog-service 启动时）。
func Root() string {
	once.Do(func() {
		if v := os.Getenv("UPLOAD_ROOT"); v != "" {
			_ = os.MkdirAll(v, 0o755)
			root = v
			return
		}
		candidates := []string{
			// 优先仓库根 uploads（services/*/uploads 常被 docker/root 占用不可写）
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

// Abs 拼接上传根下的相对路径
func Abs(elem ...string) string {
	parts := append([]string{Root()}, elem...)
	return filepath.Join(parts...)
}
