package userstatussrv

import "github.com/kianooshaz/skeleton/foundation/status"

// Status defines a set of possible flags to represent an entity's state.
// Each constant represents a specific state, and multiple states can be combined
// using bitwise operations for more complex scenarios.
const (
	// Registered indicates the user is registered.
	Registered status.Status = 1

	// Pending indicates the user is awaiting further action.
	Pending status.Status = 1 << iota

	// Inactive indicates the user is inactive.
	Inactive

	// Locked indicates the user is locked and cannot be accessed.
	Locked

	// Blocked indicates the user is blocked due to policy violations.
	Blocked

	// Suspended indicates the user is temporarily disabled.
	Suspended

	// Hidden indicates the user is hidden from regular users.
	Hidden

	// UnderReview indicates the user is under active review.
	UnderReview

	// Flagged indicates the user has been marked for special attention.
	Flagged

	Verified

	// ManuallyAdded indicates the user's status was manually set by an admin.
	ManuallyAdded

	// ManuallyVerified indicates the user's status was manually verified.
	ManuallyVerified
)
