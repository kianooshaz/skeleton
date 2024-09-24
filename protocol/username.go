package protocol

import "github.com/kianooshaz/skeleton/protocol/status"

type Username interface {
	ID() ID
	Value() string
	Status() status.Status
}

type UsernameService interface {
}
