package xerr

import (
	"fmt"
	"net/http"
)

// 业务错误码（与 HTTP 状态大致对齐，便于 httpx 映射）
const (
	CodeOK            = 0
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeConflict      = 409
	CodeTooMany       = 429
	CodeServerError   = 500
)

// CodeError 业务错误，供 httpx.ErrorCtx + SetErrorHandlerCtx 识别。
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

func New(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func Newf(code int, format string, args ...interface{}) *CodeError {
	return &CodeError{Code: code, Msg: fmt.Sprintf(format, args...)}
}

func BadRequest(msg string) *CodeError   { return New(CodeBadRequest, msg) }
func Unauthorized(msg string) *CodeError { return New(CodeUnauthorized, msg) }
func Forbidden(msg string) *CodeError    { return New(CodeForbidden, msg) }
func NotFound(msg string) *CodeError     { return New(CodeNotFound, msg) }
func Server(msg string) *CodeError       { return New(CodeServerError, msg) }

// HTTPStatus 将业务码映射为 HTTP 状态；未知码默认 400。
func HTTPStatus(code int) int {
	switch code {
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeTooMany:
		return http.StatusTooManyRequests
	case CodeServerError:
		return http.StatusInternalServerError
	case CodeBadRequest:
		return http.StatusBadRequest
	default:
		if code >= 400 && code < 600 {
			return code
		}
		return http.StatusBadRequest
	}
}
