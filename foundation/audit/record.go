package audit

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID           int
	Action       Action
	CreatedAt    time.Time
	Data         json.RawMessage
	OriginIP     string
	ResourceID   int
	ResourceType string
	UserID       int
}
