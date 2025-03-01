package protocol

import (
	"encoding/json"
	"time"
)

type Audit interface {
	Record(record Record)
}

type Record struct {
	RequestID    string
	Action       Action
	CreatedAt    time.Time
	Data         json.RawMessage
	OriginIP     string
	ResourceID   int
	ResourceType string
	UserID       int
}

type Action string

const (
	Insert Action = "insert"
	Update Action = "update"
	Delete Action = "delete"
	List   Action = "list"
	Get    Action = "get"
)
