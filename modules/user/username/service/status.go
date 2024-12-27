package service

import "github.com/kianooshaz/skeleton/foundation/status"

// Status defines a set of possible flags to represent an entity's state.
// Each constant represents a specific state, and multiple states can be combined
// using bitwise operations for more complex scenarios.
const (
	// Reserved indicates the entity's status is reserved for special cases or future actions.
	Reserved status.Status = 1

	// Hidden indicates the entity is hidden from regular users.
	Hidden status.Status = 1 << iota
)
