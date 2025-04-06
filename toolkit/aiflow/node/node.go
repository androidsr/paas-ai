package node

type Node interface {
	ID() string
	Execute(input map[string]interface{}, output map[string]interface{}, emitter chan string) bool
}
