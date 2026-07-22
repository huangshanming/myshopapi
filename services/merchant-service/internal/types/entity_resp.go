package types

// Named entity responses — aliases of DataResp keep unwrap MarshalJSON / legacy wire.

type CouponGrantResp = DataResp
type CouponResp = DataResp
type CouponStatsResp = DataResp
type HomeSlotsResp = DataResp
type HomepageOrderResp = DataResp
type MatchCouponsResp = DataResp
type SeckillConsumeResp = DataResp
type SeckillCurrentResp = DataResp
type SeckillEntryResp = DataResp
type SeckillRuleResp = DataResp
type SeckillSessionsResp = DataResp
type ShopApplicationResp = DataResp
type ShopListResp = DataResp
type ShopResp = DataResp
type SlotPackageResp = DataResp
type ThemeOrderResp = DataResp
type ThemePackageResp = DataResp
type UserCouponResp = DataResp
type WalletResp = DataResp


type GrantedCountResp struct {
	Granted int64 `json:"granted"`
}
