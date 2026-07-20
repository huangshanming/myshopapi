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
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ParentID  uint64           `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	Name      string           `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	Status    int8             `gorm:"column:status;default:1" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (CommunityArticleCategory) TableName() string { return "community_article_category" }

type CommunityArticle struct {
	ID                uint64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID            uint64            `gorm:"column:shop_id;not null;index" json:"shop_id"`
	AuthorUserID      uint64            `gorm:"column:author_user_id;not null;default:0;index" json:"author_user_id"`
	CategoryID        uint64            `gorm:"column:category_id;not null;default:0" json:"category_id"`
	Title             string            `gorm:"column:title;type:varchar(200);not null" json:"title"`
	CoverURL          string            `gorm:"column:cover_url;type:varchar(500)" json:"cover_url"`
	Content           string            `gorm:"column:content;type:mediumtext" json:"content"`
	AuditStatus       string            `gorm:"column:audit_status;type:enum('pending','approved','rejected');default:pending" json:"audit_status"`
	RejectReason      string            `gorm:"column:reject_reason;type:varchar(500)" json:"reject_reason"`
	Status            string            `gorm:"column:status;type:enum('draft','scheduled','published','offline','deleted');default:draft" json:"status"`
	SchedulePublishAt *common.LocalTime `gorm:"column:schedule_publish_at" json:"schedule_publish_at,omitempty"`
	IsTop             int8              `gorm:"column:is_top;default:0" json:"is_top"`
	ViewCount         uint64            `gorm:"column:view_count;default:0" json:"view_count"`
	LikeCount         uint64            `gorm:"column:like_count;default:0" json:"like_count"`
	AudienceCount     uint64            `gorm:"column:audience_count;default:0" json:"audience_count"`
	ReadCount         uint64            `gorm:"column:read_count;default:0" json:"read_count"`
	CollectCount      uint64            `gorm:"column:collect_count;default:0" json:"collect_count"`
	CommentCount      uint64            `gorm:"column:comment_count;default:0" json:"comment_count"`
	PublishedAt       *common.LocalTime `gorm:"column:published_at" json:"published_at,omitempty"`
	DeletedAt         *common.LocalTime `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	CreatedBy         uint64            `gorm:"column:created_by;not null;default:0" json:"created_by"`
	CreatedAt         common.LocalTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         common.LocalTime  `gorm:"column:updated_at" json:"updated_at"`
}

func (CommunityArticle) TableName() string { return "community_article" }

type CommunityArticleImg struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ArticleID uint64           `gorm:"column:article_id;not null;index" json:"article_id"`
	ShopID    uint64           `gorm:"column:shop_id;not null" json:"shop_id"`
	URL       string           `gorm:"column:url;type:varchar(500);not null" json:"url"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (CommunityArticleImg) TableName() string { return "community_article_img" }

type CommunityArticleComment struct {
	ID            uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ArticleID     uint64           `gorm:"column:article_id;not null;index" json:"article_id"`
	ShopID        uint64           `gorm:"column:shop_id;not null;index" json:"shop_id"`
	UserID        uint64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	ParentID      uint64           `gorm:"column:parent_id;not null;default:0" json:"parent_id"`
	RootID        uint64           `gorm:"column:root_id;not null;default:0;index" json:"root_id"`
	ReplyToUserID uint64           `gorm:"column:reply_to_user_id;not null;default:0" json:"reply_to_user_id"`
	Content       string           `gorm:"column:content;type:varchar(1000);not null" json:"content"`
	Status        string           `gorm:"column:status;type:enum('visible','hidden','deleted');default:visible" json:"status"`
	CreatedAt     common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     common.LocalTime `gorm:"column:updated_at" json:"updated_at"`

	UserNickname     string `gorm:"-" json:"user_nickname,omitempty"`
	ReplyToNickname  string `gorm:"-" json:"reply_to_nickname,omitempty"`
	Children         []CommunityArticleComment `gorm:"-" json:"children,omitempty"`
}

func (CommunityArticleComment) TableName() string { return "community_article_comment" }

// CommunityCommentEmoji 评论表情包（超管配置）
type CommunityCommentEmoji struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string           `gorm:"column:name;type:varchar(64)" json:"name"`
	ImageURL  string           `gorm:"column:image_url;type:varchar(500)" json:"image_url"`
	Sort      int              `gorm:"column:sort" json:"sort"`
	Status    int8             `gorm:"column:status;default:1" json:"status"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt common.LocalTime `gorm:"column:updated_at" json:"updated_at"`
}

func (CommunityCommentEmoji) TableName() string { return "community_comment_emojis" }

type ArticleLike struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64           `gorm:"column:user_id;not null;uniqueIndex:uk_user_article" json:"user_id"`
	ArticleID uint64           `gorm:"column:article_id;not null;uniqueIndex:uk_user_article;index" json:"article_id"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ArticleLike) TableName() string { return "article_likes" }

type ArticleFavorite struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64           `gorm:"column:user_id;not null;uniqueIndex:uk_user_article" json:"user_id"`
	ArticleID uint64           `gorm:"column:article_id;not null;uniqueIndex:uk_user_article;index" json:"article_id"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ArticleFavorite) TableName() string { return "article_favorites" }

type ArticleAudience struct {
	ID        uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint64           `gorm:"column:user_id;not null;uniqueIndex:uk_user_article" json:"user_id"`
	ArticleID uint64           `gorm:"column:article_id;not null;uniqueIndex:uk_user_article;index" json:"article_id"`
	CreatedAt common.LocalTime `gorm:"column:created_at" json:"created_at"`
}

func (ArticleAudience) TableName() string { return "article_audiences" }
