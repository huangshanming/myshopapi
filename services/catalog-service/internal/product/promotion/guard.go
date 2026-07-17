package promotion

// Guard 活动商品拦截（暂无促销域，默认放行；后续可查活动表）
type Guard interface {
	CanOffSale(productID uint64) bool
	CanDelete(productID uint64) bool
}

type noopGuard struct{}

func (noopGuard) CanOffSale(uint64) bool { return true }
func (noopGuard) CanDelete(uint64) bool  { return true }

func NewNoop() Guard { return noopGuard{} }
