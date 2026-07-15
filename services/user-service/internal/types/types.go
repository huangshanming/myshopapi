package types

// LoginReq 登录请求
type LoginReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	ShopID   uint64 `json:"shop_id"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
