package pagination

import "math"

type Response[T any] struct {
	Page
	TotalRows int `json:"total_rows" bson:"total_rows"`
	TotalPage int `json:"total_page" bson:"total_page"`
	Data      []T `json:"data" bson:"data"`
}

func NewResponse[T any](page Page, totalRows int, rows []T) Response[T] {
	return Response[T]{
		Page:      page,
		TotalRows: totalRows,
		TotalPage: int(math.Ceil(float64(totalRows) / float64(page.PageRows))),
		Data:      rows,
	}
}
