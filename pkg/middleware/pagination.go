package middleware

import (
	"net/http"

	"mymall/pkg/pagination"
)

func ParsePage(r *http.Request) (page, pageSize int) {
	q := r.URL.Query()
	req := pagination.PageReq{
		Page:     atoiDefault(q.Get("page"), 1),
		PageSize: atoiDefault(q.Get("page_size"), 10),
	}
	page, pageSize, _ = pagination.Normalize(&req)
	return page, pageSize
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return def
		}
		n = n*10 + int(c-'0')
	}
	if n == 0 {
		return def
	}
	return n
}
