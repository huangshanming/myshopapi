package apidoc

// Response 统一 API 响应包装（code=200 表示业务成功）
type Response struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"success"`
	Data any    `json:"data"`
}

// EmptyData 无 data 字段时的占位
type EmptyData struct{}
