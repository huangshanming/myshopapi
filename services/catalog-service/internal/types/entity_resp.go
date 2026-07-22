package types

// Named entity responses — aliases of DataResp keep unwrap MarshalJSON / legacy wire.

type ArticleResp = DataResp
type ArticleStatsResp = DataResp
type AttrTemplateResp = DataResp
type BannerResp = DataResp
type CategoryResp = DataResp
type CommentResp = DataResp
type EmojiResp = DataResp
type ImportResultResp = DataResp
type ProductJobResp = DataResp
type ProductResp = DataResp
type RoleMenuIdsResp = DataResp
type ShopRoleResp = DataResp
type StockWarningsResp = DataResp
type TagResp = DataResp

// SalesRankResp matches mall-uni: { total, list: [...] } (not nested under data).
type SalesRankResp = PageListResp


type ShopAuthMeResp struct {
	Perms    interface{} `json:"perms"`
	Menus    interface{} `json:"menus"`
	MenuTree interface{} `json:"menu_tree"`
	IsOwner  bool        `json:"is_owner"`
}

type BindStaffResp struct {
	UserId uint64 `json:"user_id"`
	Msg    string `json:"msg"`
}
