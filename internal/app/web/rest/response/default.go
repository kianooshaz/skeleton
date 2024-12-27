package response

import "github.com/kianooshaz/skeleton/foundation/code"

func New[T any](data T, meta *Meta) Default[T] {
	return Default[T]{
		Code: code.Success,
		Data: data,
		Meta: meta,
	}
}

type Default[T any] struct {
	Code code.Code `json:"code"`
	Data T         `json:"data"`
	Meta *Meta     `json:"meta,omitempty"`
}
