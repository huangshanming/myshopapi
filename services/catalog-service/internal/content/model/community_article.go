package model

import "mymall/common"

const (
	ArticleAuditPending  = "pending"
	ArticleAuditApproved = "approved"
	ArticleAuditRejected = "rejected"
)

const (
	ArticleDraft     = "draft"
	ArticleScheduled = "scheduled"
	ArticlePublished = "published"
	ArticleOffline   = "offline"
	ArticleDeleted   = "deleted"
)

const (
	CommentVisible = "visible"
	CommentHidden  = "hidden"
	CommentDeleted = "deleted"
)

type CommunityArticleCategory struct {
	ID        uint64           `db:"id" json:"id"`
	ParentID  uint64           `db:"parent_id" json:"parent_id"`
	Name      string           `db:"name" json:"name"`
	Sort      int              `db:"sort" json:"sort"`
	Status    int8             `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (CommunityArticleCategory) TableName() string { return "community_article_category" }

type CommunityArticle struct {
	ID                uint64            `db:"id" json:"id"`
	ShopID            uint64            `db:"shop_id" json:"shop_id"`
	AuthorUserID      uint64            `db:"author_user_id" json:"author_user_id"`
	CategoryID        uint64            `db:"category_id" json:"category_id"`
	Title             string            `db:"title" json:"title"`
	CoverURL          string            `db:"cover_url" json:"cover_url"`
	Content           string            `db:"content" json:"content"`
	AuditStatus       string            `db:"audit_status" json:"audit_status"`
	RejectReason      string            `db:"reject_reason" json:"reject_reason"`
	Status            string            `db:"status" json:"status"`
	SchedulePublishAt *common.LocalTime `db:"schedule_publish_at" json:"schedule_publish_at,omitempty"`
	IsTop             int8              `db:"is_top" json:"is_top"`
	ViewCount         uint64            `db:"view_count" json:"view_count"`
	LikeCount         uint64            `db:"like_count" json:"like_count"`
	AudienceCount     uint64            `db:"audience_count" json:"audience_count"`
	ReadCount         uint64            `db:"read_count" json:"read_count"`
	CollectCount      uint64            `db:"collect_count" json:"collect_count"`
	CommentCount      uint64            `db:"comment_count" json:"comment_count"`
	PublishedAt       *common.LocalTime `db:"published_at" json:"published_at,omitempty"`
	DeletedAt         *common.LocalTime `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedBy         uint64            `db:"created_by" json:"created_by"`
	CreatedAt         common.LocalTime  `db:"created_at" json:"created_at"`
	UpdatedAt         common.LocalTime  `db:"updated_at" json:"updated_at"`
}

func (CommunityArticle) TableName() string { return "community_article" }

type CommunityArticleImg struct {
	ID        uint64           `db:"id" json:"id"`
	ArticleID uint64           `db:"article_id" json:"article_id"`
	ShopID    uint64           `db:"shop_id" json:"shop_id"`
	URL       string           `db:"url" json:"url"`
	Sort      int              `db:"sort" json:"sort"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
}

func (CommunityArticleImg) TableName() string { return "community_article_img" }

type CommunityArticleComment struct {
	ID            uint64           `db:"id" json:"id"`
	ArticleID     uint64           `db:"article_id" json:"article_id"`
	ShopID        uint64           `db:"shop_id" json:"shop_id"`
	UserID        uint64           `db:"user_id" json:"user_id"`
	ParentID      uint64           `db:"parent_id" json:"parent_id"`
	RootID        uint64           `db:"root_id" json:"root_id"`
	ReplyToUserID uint64           `db:"reply_to_user_id" json:"reply_to_user_id"`
	Content       string           `db:"content" json:"content"`
	Status        string           `db:"status" json:"status"`
	CreatedAt     common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `db:"updated_at" json:"updated_at"`

	UserNickname     string `db:"-" json:"user_nickname,omitempty"`
	ReplyToNickname  string `db:"-" json:"reply_to_nickname,omitempty"`
	Children         []CommunityArticleComment `db:"-" json:"children,omitempty"`
}

func (CommunityArticleComment) TableName() string { return "community_article_comment" }

// CommunityCommentEmoji 评论表情包（超管配置）
type CommunityCommentEmoji struct {
	ID        uint64           `db:"id" json:"id"`
	Name      string           `db:"name" json:"name"`
	ImageURL  string           `db:"image_url" json:"image_url"`
	Sort      int              `db:"sort" json:"sort"`
	Status    int8             `db:"status" json:"status"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
	UpdatedAt common.LocalTime `db:"updated_at" json:"updated_at"`
}

func (CommunityCommentEmoji) TableName() string { return "community_comment_emojis" }

type ArticleLike struct {
	ID        uint64           `db:"id" json:"id"`
	UserID    uint64           `db:"user_id" json:"user_id"`
	ArticleID uint64           `db:"article_id" json:"article_id"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
}

func (ArticleLike) TableName() string { return "article_likes" }

type ArticleFavorite struct {
	ID        uint64           `db:"id" json:"id"`
	UserID    uint64           `db:"user_id" json:"user_id"`
	ArticleID uint64           `db:"article_id" json:"article_id"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
}

func (ArticleFavorite) TableName() string { return "article_favorites" }

type ArticleAudience struct {
	ID        uint64           `db:"id" json:"id"`
	UserID    uint64           `db:"user_id" json:"user_id"`
	ArticleID uint64           `db:"article_id" json:"article_id"`
	CreatedAt common.LocalTime `db:"created_at" json:"created_at"`
}

func (ArticleAudience) TableName() string { return "article_audiences" }
