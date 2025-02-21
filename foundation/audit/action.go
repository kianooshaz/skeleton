package audit

type Action string

const (
	Insert Action = "insert"
	Update Action = "update"
	Delete Action = "delete"
	List   Action = "list"
	Get    Action = "get"
)
