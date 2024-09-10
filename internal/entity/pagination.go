package entity

import (
	"fmt"
)

type Pagination struct {
	Limit  int32
	Offset int32
}

func NewPagination(limit *int32, offset *int32) Pagination {
	p := Pagination{
		Limit:  5,
		Offset: 0,
	}
	if limit != nil {
		p.Limit = *limit
	}
	if offset != nil {
		p.Offset = *offset
	}
	return p
}

func (p Pagination) ToSQL() string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", p.Limit, p.Offset)
}
