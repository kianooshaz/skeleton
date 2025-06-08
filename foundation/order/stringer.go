package order

// StringerFunc defines the function signature for constructing string representations
// of By instances. This allows customization for different storage backends.
type StringerFunc func(OrderBy) string

// String generates the string representation of a By instance using the provided StringerFunc.
func (b OrderBy) String(stringer StringerFunc) string {
	return stringer(b)
}
