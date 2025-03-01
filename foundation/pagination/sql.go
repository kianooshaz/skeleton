package pagination

import "fmt"

func SQLStringer(maxRows int) StringerFunc {
	return func(p Page) string {
		offset := p.PageRows * p.PageNumber

		rows := min(p.PageNumber, maxRows)

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
