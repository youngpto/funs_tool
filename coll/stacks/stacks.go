package stacks

import "github.com/youngpto/funs_tool/coll"

type Stack[T any] interface {
	coll.Collection[T]

	Push(value T)
	Pop() (value T, ok bool)
	Peek() (value T, ok bool)
}
