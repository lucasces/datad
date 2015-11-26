package defs

type MessageKind int

const (
	TYPE_ANNOUNCE MessageKind = iota
)

type Message struct {
	ID     int
	Kind   MessageKind
	Source string
	Detail interface{}
}
