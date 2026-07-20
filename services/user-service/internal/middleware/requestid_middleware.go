package middleware

import (
	pkgmw "mymall/pkg/middleware"
	"net/http"
)

type RequestIDMiddleware struct{}

func NewRequestIDMiddleware() *RequestIDMiddleware { return &RequestIDMiddleware{} }

func (m *RequestIDMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return pkgmw.RequestID()(next)
}
