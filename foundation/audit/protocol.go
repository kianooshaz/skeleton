package audit

type Protocol interface {
	Record(record Record)
}
