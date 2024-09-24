package response

import (
	"math"

	"github.com/kianooshaz/skeleton/foundation/pagination"
)

func NewMeta(page pagination.Page, total, from, to int64, path string) *Meta {
	var lastPage float64
	if total != 0 {
		lastPage = math.Ceil(float64(total) / float64(page.RowsPerPage()))
	}

	return &Meta{
		CurrentPage: page.Number(),
		LastPage:    int(lastPage),
		PerPage:     page.RowsPerPage(),
		Total:       total,
		From:        from,
		To:          to,
		Path:        path,
	}
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`

	From int64 `json:"from"`
	To   int64 `json:"to"`

	Path string `json:"path"`
}
