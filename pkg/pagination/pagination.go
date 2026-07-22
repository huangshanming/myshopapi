package pagination

// Page helpers shared by business services (sqlx). GORM Paginate was removed;
// inventory-sync does not use this package.

type PageReq struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	OrderBy  string `form:"order_by" json:"order_by"`
}

type PageRes[T any] struct {
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalPage int   `json:"total_page"`
	List      []T   `json:"list"`
}

func Normalize(req *PageReq) (page, pageSize, offset int) {
	page = req.Page
	if page < 1 {
		page = 1
	}
	pageSize = req.PageSize
	switch {
	case pageSize < 1:
		pageSize = 10
	case pageSize > 100:
		pageSize = 100
	}
	offset = (page - 1) * pageSize
	return page, pageSize, offset
}

func NewPageRes[T any](list []T, total int64, page, pageSize int) *PageRes[T] {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPage := (int(total) + pageSize - 1) / pageSize
	return &PageRes[T]{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
		List:      list,
	}
}
