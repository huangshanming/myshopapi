package dto

// ErrorResp 错误响应（HTTP 200，code 非 200）
type ErrorResp struct {
	Code int    `json:"code" example:"400"`
	Msg  string `json:"msg" example:"参数错误"`
	Data any    `json:"data"`
}
