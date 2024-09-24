package status

type Status int64

const Unset Status = 0

const (
	Registered Status = 1
	Pending    Status = 1 << iota
	Inactive
	Locked
	Blocked
	Suspended
	Hidden
	Revoked
	Deleted
	Reserved
	UnderReview
	Flagged
	ManuallyAdded
	ManuallyVerified
)

// Has checks if a specific status is present.
func (s *Status) Has(ss Status) bool {
	return *s&ss != 0
}

// Add adds a status to the current status.
func (s *Status) Add(ss Status) {
	*s |= ss
}

// Remove removes a status from the current status.
func (s *Status) Remove(ss Status) {
	*s &^= ss
}

// Clear removes all statuses from the current status.
func (s *Status) Clear() {
	*s = 0
}
