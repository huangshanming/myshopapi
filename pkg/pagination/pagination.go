package pagination

import "gorm.io/gorm"

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

func Paginate[T any](db *gorm.DB, pageReq *PageReq) (*PageRes[T], error) {
	page, pageSize, offset := Normalize(pageReq)

	var (
		list  []T
		total int64
	)

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	dbQuery := db.Offset(offset).Limit(pageSize)
	if pageReq.OrderBy != "" {
		dbQuery = dbQuery.Order(pageReq.OrderBy)
	}
	if err := dbQuery.Find(&list).Error; err != nil {
		return nil, err
	}

	totalPage := (int(total) + pageSize - 1) / pageSize
	return &PageRes[T]{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
		List:      list,
	}, nil
}
