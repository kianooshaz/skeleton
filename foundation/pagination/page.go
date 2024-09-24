package pagination

import (
	"errors"
	"fmt"
	"strconv"
)

type Page struct {
	number int
	rows   int
}

var (
	ErrPageConversion    = errors.New("invalid page value, must be a number")
	ErrRowsConversion    = errors.New("invalid rows value, must be a number")
	ErrPageValueTooSmall = errors.New("page value too small, must be larger than 0")
	ErrRowsValueTooSmall = errors.New("rows value too small, must be larger than 0")
	ErrRowsValueTooLarge = errors.New("rows value too large, must be less than the maximum allowed")
)

func Parse(page string, rowsPerPage string, maxRowsPerPage int) (Page, error) {
	number := 1
	if page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, ErrPageConversion
		}
	}

	rows := 10
	if rowsPerPage != "" {
		var err error
		rows, err = strconv.Atoi(rowsPerPage)
		if err != nil {
			return Page{}, ErrRowsConversion
		}
	}

	if number <= 0 {
		return Page{}, ErrPageValueTooSmall
	}

	if rows <= 0 {
		return Page{}, ErrRowsValueTooSmall
	}

	if rows > maxRowsPerPage {
		return Page{}, ErrRowsValueTooLarge
	}

	p := Page{
		number: number,
		rows:   rows,
	}

	return p, nil
}

func MustParse(page string, rowsPerPage string, maxRowsPerPage int) Page {
	pg, err := Parse(page, rowsPerPage, maxRowsPerPage)
	if err != nil {
		panic(err)
	}

	return pg
}

func (p Page) String() string {
	return fmt.Sprintf("page: %d rows: %d", p.number, p.rows)
}

func (p Page) Number() int {
	return p.number
}

// Offset returns the read offset.
func (p Page) Offset() int {
	return (p.number - 1) * p.rows
}

func (p Page) RowsPerPage() int {
	return p.rows
}
