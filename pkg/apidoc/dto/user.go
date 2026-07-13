package dto

// LoginReq 登录请求
type LoginReq struct {
	Mobile   string `json:"mobile" binding:"required,len=11" example:"13800138000"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	Mobile   string `json:"mobile" binding:"required,len=11" example:"13800138000"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// UserInfo 用户信息（不含密码）
type UserInfo struct {
	ID        uint64 `json:"id" example:"1"`
	Mobile    string `json:"mobile" example:"13800138000"`
	Nickname  string `json:"nickname" example:"用户8000"`
	Avatar    string `json:"avatar" example:""`
	Gender    int    `json:"gender" example:"0"`
	Status    int    `json:"status" example:"1"`
	CreatedAt string `json:"created_at" example:"2026-01-01 12:00:00"`
	UpdatedAt string `json:"updated_at" example:"2026-01-01 12:00:00"`
}

// LoginResp 登录响应 data
type LoginResp struct {
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserInfo `json:"user"`
}

// UserProfileResp 个人资料响应
type UserProfileResp struct {
	UserInfo
}
