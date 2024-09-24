package protocol

import "encoding/hex"

type ID [16]byte

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}
