package middleware

import (
	"net/http"

	pkgmw "mymall/pkg/middleware"
)

type RequestIDMiddleware struct{}

func NewRequestIDMiddleware() *RequestIDMiddleware { return &RequestIDMiddleware{} }

func (m *RequestIDMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequestID()(next)
}
