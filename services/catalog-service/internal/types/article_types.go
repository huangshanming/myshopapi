package types

type ArticleSaveReq struct {
	ShopID            uint64   `json:"shop_id"`
	CategoryID        uint64   `json:"category_id"`
	Title             string   `json:"title"`
	CoverURL          string   `json:"cover_url"`
	Content           string   `json:"content"`
	SchedulePublishAt string   `json:"schedule_publish_at"` // 2006-01-02 15:04:05 或空
	IsTop             *int8    `json:"is_top"`              // 仅管理员
	ImageURLs         []string `json:"image_urls"`
	Submit            bool     `json:"submit"` // 商家：true=提交审核
}

type ArticleAuditReq struct {
	Pass         bool   `json:"pass"`
	RejectReason string `json:"reject_reason"`
}

type ArticleBatchAuditReq struct {
	IDs          []uint64 `json:"ids"`
	Pass         bool     `json:"pass"`
	RejectReason string   `json:"reject_reason"`
}

type ArticleTopReq struct {
	IsTop int8 `json:"is_top"`
}

type ArticleCategorySaveReq struct {
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
	Status   *int8  `json:"status"`
}

type ArticleCommentPatchReq struct {
	Status string `json:"status"` // visible / hidden / deleted
}
