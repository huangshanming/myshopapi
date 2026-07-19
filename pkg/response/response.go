package response

import (
	"encoding/json"
	"net/http"

	"mymall/pkg/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 旧信封（仅兼容文档/测试）；新契约成功体为业务 DTO 本身。
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func writeJSON(w http.ResponseWriter, httpStatus int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}

// Success 新契约：HTTP 200，body 即为业务 DTO（对齐 httpx.OkJson）。
func Success(w http.ResponseWriter, data interface{}, msg string) {
	_ = msg
	if data == nil {
		httpx.Ok(w)
		return
	}
	httpx.OkJson(w, data)
}

// Error 新契约：非 2xx + {code,msg}（对齐 httpx.Error）。
func Error(w http.ResponseWriter, msg string, code int) {
	if code == 0 {
		code = xerr.CodeBadRequest
	}
	status := xerr.HTTPStatus(code)
	writeJSON(w, status, xerr.CodeMsg{Code: code, Msg: msg})
}

// AbortJSON 鉴权等硬中断：指定 HTTP 状态 + {code,msg}。
func AbortJSON(w http.ResponseWriter, httpStatus int, code int, msg string) {
	if code == 0 {
		code = httpStatus
	}
	writeJSON(w, httpStatus, xerr.CodeMsg{Code: code, Msg: msg})
}
