package auditproto

import "github.com/google/uuid"

type RecordID uuid.UUID

func (r RecordID) String() string {
	return uuid.UUID(r).String()
}
