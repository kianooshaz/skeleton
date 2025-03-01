// Package pagination provides utilities for managing pagination logic.
package pagination

// Page holds pagination details: page number and rows per page.
type Page struct {
	PageNumber int `json:"page_number" bson:"page_number"`
	PageRows   int `json:"page_rows" bson:"page_rows"`
}

type StringerFunc func(Page) string

func (p Page) String(stringer StringerFunc) string {
	return stringer(p)
}
