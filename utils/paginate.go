package utils

import (
	"mymall/pkg/pagination"

	"gorm.io/gorm"
)

func Paginate[T any](db *gorm.DB, pageReq *pagination.PageReq) (*pagination.PageRes[T], error) {
	return pagination.Paginate[T](db, pageReq)
}
