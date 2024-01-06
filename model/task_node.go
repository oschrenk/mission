package model

type Node[T any] interface {
	// Children returns an array of pointers to all children of this node.
	Children() []Node[T]

	Value() T
}

type Tree[T any] struct {
	roots []Node[T]
}

func (t *Tree[T]) Roots() []Node[T] {
	return t.roots
}
