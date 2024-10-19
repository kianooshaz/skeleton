package protocol

import "encoding/hex"

// ID represents a 16-byte identifier, typically used for unique entity representation.
type ID [16]byte

// String returns the hexadecimal string representation of the ID.
// It converts the 16-byte ID into a readable hex-encoded format.
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}
