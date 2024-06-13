package queues

import "github.com/youngpto/funs_tool/coll"

type Queue[T comparable] interface {
	coll.Collection[T]

	Push(v T)
	Pop() (T, bool)
	Peek() (T, bool)
}
