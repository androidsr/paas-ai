package node

type Node interface {
	ID() string
	Execute(input map[string]any, output map[string]any, emitter chan string) error
}
