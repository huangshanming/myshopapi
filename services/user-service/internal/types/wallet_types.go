package types

type WalletAdjustReq struct {
	Field  string  `json:"field"`
	Amount float64 `json:"amount"`
	Remark string  `json:"remark"`
}

type WalletOrderOpReq struct {
	UserID  uint64  `json:"user_id"`
	Amount  float64 `json:"amount"`
	OrderID uint64  `json:"order_id"`
	OrderNo string  `json:"order_no"`
}
