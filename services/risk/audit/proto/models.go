package auditproto

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID           RecordID        `json:"id"`
	RequestID    string          `json:"request_id"`
	Action       Action          `json:"action"`
	CreatedAt    time.Time       `json:"created_at"`
	Data         json.RawMessage `json:"data"`
	OriginIP     string          `json:"origin_ip"`
	ResourceID   int             `json:"resource_id"`
	ResourceType string          `json:"resource_type"`
	UserID       int             `json:"user_id"`
}

type Action string

const (
	Insert Action = "insert"
	Update Action = "update"
	Delete Action = "delete"
	List   Action = "list"
	Get    Action = "get"
)
