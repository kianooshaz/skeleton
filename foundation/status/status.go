package status

// Status represents a set of status flags, using a uint64 as the underlying type.
type Status uint

// Unset is the default value for Status, representing no status.
const Unset Status = 0

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
