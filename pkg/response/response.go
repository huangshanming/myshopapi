package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func writeJSON(w http.ResponseWriter, httpStatus int, body Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, data interface{}, msg string) {
	if msg == "" {
		msg = "success"
	}
	writeJSON(w, http.StatusOK, Response{Code: 200, Msg: msg, Data: data})
}

// Error 业务错误：HTTP 仍 200，用 body.code 表达（与现前端约定一致）
func Error(w http.ResponseWriter, msg string, code int) {
	if code == 0 {
		code = 400
	}
	writeJSON(w, http.StatusOK, Response{Code: code, Msg: msg, Data: nil})
}

func AbortJSON(w http.ResponseWriter, httpStatus int, code int, msg string) {
	writeJSON(w, httpStatus, Response{Code: code, Msg: msg, Data: nil})
}
