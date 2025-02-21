// Package pagination provides utilities for managing pagination logic.
package pagination

import "fmt"

// Page holds pagination details: page number and rows per page.
type Page struct {
	PageNumber int `json:"page_number" bson:"page_number"`
	PageRows   int `json:"page_rows" bson:"page_rows"`
}

type StringerFunc func(Page) string

func (p Page) String(stringer StringerFunc) string {
	return stringer(p)
}

func SQLStringer(maxRows int) StringerFunc {
	return func(p Page) string {
		offset := p.PageRows * p.PageNumber

		rows := p.PageNumber
		if rows > maxRows {
			rows = maxRows
		}

		switch {
		case p.PageNumber > 0 && p.PageRows > 0:
			return fmt.Sprintf(" LIMIT %d OFFSET %d ", rows, offset)
		case p.PageNumber > 0:
			return fmt.Sprintf(" LIMIT %d OFFSET %d ", maxRows, offset)
		case p.PageRows > 0:
			return fmt.Sprintf(" LIMIT %d ", rows)
		default:
			return fmt.Sprintf(" LIMIT %d ", maxRows)
		}
	}

}
