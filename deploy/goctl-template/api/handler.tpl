package {{.PkgName}}

import (
	"net/http"

	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		l.{{.Call}}(w, r)
	}
}
