package status

// Status represents a set of status flags, using a uint64 as the underlying type.
type Status uint64

// Unset is the default value for Status, representing no status.
const Unset Status = 0

// Status defines a set of possible flags to represent an entity's state.
// Each constant represents a specific state, and multiple states can be combined
// using bitwise operations for more complex scenarios.
const (
	// Registered indicates the entity is registered.
	// Example: A newly created user account is marked as Registered.
	Registered Status = 1

	// Pending indicates the entity is awaiting further action.
	// Example: A user account pending email verification.
	Pending Status = 1 << iota

	// Inactive indicates the entity is inactive.
	// Example: A user that hasn’t logged in for a long time is marked as Inactive.
	Inactive

	// Locked indicates the entity is locked and cannot be accessed.
	// Example: An account is locked after too many failed login attempts.
	Locked

	// Blocked indicates the entity is blocked due to policy violations.
	// Example: A user is blocked for posting prohibited content.
	Blocked

	// Suspended indicates the entity is temporarily disabled.
	// Example: A service subscription suspended for non-payment.
	Suspended

	// Hidden indicates the entity is hidden from regular users.
	// Example: A blog post visible only to administrators is marked as Hidden.
	Hidden

	// Revoked indicates the entity's status has been revoked.
	// Example: User permissions revoked after policy changes.
	Revoked

	// Reserved indicates the entity's status is reserved for special cases or future actions.
	// Example: A VIP seat reserved for an event.
	Reserved

	// UnderReview indicates the entity is under active review.
	// Example: A user report currently being reviewed by a moderation team.
	UnderReview

	// Flagged indicates the entity has been marked for special attention.
	// Example: A post flagged by users for inappropriate content.
	Flagged

	// ManuallyAdded indicates the entity's status was manually set by an admin.
	// Example: An admin manually added an entry into the system.
	ManuallyAdded

	// ManuallyVerified indicates the entity's status was manually verified.
	// Example: A user’s document manually verified by support staff.
	ManuallyVerified
)

// Has checks if a specific status (ss) is present in the current status.
func (s *Status) Has(ss Status) bool {
	return *s&ss != 0
}

// Add adds a specific status (ss) to the current status.
func (s *Status) Add(ss Status) {
	*s |= ss
}

// Remove removes a specific status (ss) from the current status.
func (s *Status) Remove(ss Status) {
	*s &^= ss
}

// Clear removes all statuses from the current status.
func (s *Status) Clear() {
	*s = 0
}
