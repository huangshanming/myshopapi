package xerr

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// CodeMsg httpx 错误响应体（标准形态）。
type CodeMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// RegisterErrorHandler 注册全局 httpx 错误处理：业务 CodeError → 对应 HTTP；其余 → 500。
func RegisterErrorHandler() {
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		var ce *CodeError
		if errors.As(err, &ce) && ce != nil {
			return HTTPStatus(ce.Code), CodeMsg{Code: ce.Code, Msg: ce.Msg}
		}
		logx.WithContext(ctx).Errorf("internal error: %+v", err)
		return http.StatusInternalServerError, CodeMsg{Code: CodeServerError, Msg: "服务器内部错误"}
	})
}
