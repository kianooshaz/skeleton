package order

type ParserFunc func(string) (OrderBy, error)

func Parse(parser ParserFunc, orderBy string) (OrderBy, error) {
	return parser(orderBy)
}
