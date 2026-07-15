package apidoc

import "net/http"

// MountSwagger 已从 gin-swagger 迁出；OpenAPI 请用 docs/scalar + scripts/serve-docs.sh。
// 保留空实现避免旧调用编译失败；可选挂一个说明页。
func MountSwagger(mux interface{}) {}

func SwaggerInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Swagger UI 已下线，请访问 docs：bash scripts/serve-docs.sh → http://localhost:9099/scalar/index.html\n"))
	}
}
