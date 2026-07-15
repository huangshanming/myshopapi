package httpserver

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func PathParam(r *http.Request, key string) string {
	if v := r.PathValue(key); v != "" {
		return v
	}
	if vars := pathvar.Vars(r); vars != nil {
		return vars[key]
	}
	return ""
}
