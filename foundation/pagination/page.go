// Package pagination provides utilities for managing pagination logic.
package pagination

import (
	"errors"
	"fmt"
	"strconv"
)

// Page holds pagination details: page number and rows per page.
type Page struct {
	Number int
	Rows   int
}

// Error definitions for invalid pagination inputs.
var (
	ErrInvalidPage  = errors.New("page must be a positive integer")
	ErrInvalidRows  = errors.New("rows must be a positive integer")
	ErrRowsTooLarge = errors.New("rows exceed maximum allowed")
	ErrPageParse    = errors.New("invalid page value, must be a number")
	ErrRowsParse    = errors.New("invalid rows value, must be a number")
)

// NewPage parses strings into a Page struct with validation.
func NewPage(pageStr, rowsStr string, maxRows int) (Page, error) {
	page, err := parseWithDefault(pageStr, 1)
	if err != nil {
		return Page{}, ErrPageParse
	}

	rows, err := parseWithDefault(rowsStr, 10)
	if err != nil {
		return Page{}, ErrRowsParse
	}

	if page <= 0 {
		return Page{}, ErrInvalidPage
	}
	if rows <= 0 {
		return Page{}, ErrInvalidRows
	}
	if rows > maxRows {
		return Page{}, ErrRowsTooLarge
	}

	return Page{Number: page, Rows: rows}, nil
}

// MustNewPage panics if NewPage returns an error.
func MustNewPage(pageStr, rowsStr string, maxRows int) Page {
	p, err := NewPage(pageStr, rowsStr, maxRows)
	if err != nil {
		panic(err)
	}
	return p
}

// Helper function to parse string to int with a fallback default.
func parseWithDefault(value string, defaultVal int) (int, error) {
	if value == "" {
		return defaultVal, nil
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// String returns a human-readable representation of Page.
func (p Page) String() string {
	return fmt.Sprintf("page: %d rows: %d", p.Number, p.Rows)
}

// Offset calculates the pagination offset.
func (p Page) Offset() int {
	return (p.Number - 1) * p.Rows
}
