package utils

import (
	"gorm.io/gorm"
	"mymall/common"
)

func Paginate[T any](db *gorm.DB, pageReq *common.PageReq) (*common.PageRes[T], error) {
	// 1. 分页参数默认值与合法性校验（全局统一规则，无需每个接口重复写）
	page := pageReq.Page
	if page < 1 {
		page = 1
	}
	pageSize := pageReq.PageSize
	switch {
	case pageSize < 1:
		pageSize = 10
	case pageSize > 100:
		pageSize = 100 // 全局限制最大页大小，防止性能问题
	}
	offset := (page - 1) * pageSize

	// 2. 声明变量接收数据
	var (
		list  []T
		total int64
	)

	// 3. 查询总记录数（自动忽略Offset/Limit，仅统计符合业务条件的总条数）
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 4. 执行分页查询（支持排序）
	dbQuery := db.Offset(offset).Limit(pageSize)
	if pageReq.OrderBy != "" {
		dbQuery = dbQuery.Order(pageReq.OrderBy) // 动态排序
	}
	if err := dbQuery.Find(&list).Error; err != nil {
		return nil, err
	}

	// 5. 计算总页数，封装通用分页结果
	totalPage := (int(total) + pageSize - 1) / pageSize
	return &common.PageRes[T]{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
		List:      list,
	}, nil
}
