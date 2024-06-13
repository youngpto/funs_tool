package deque

import "github.com/youngpto/funs_tool/coll/queues"

type Deque[T comparable] interface {
	queues.Queue[T]

	PushFront(v T)
	PopBack() (T, bool)
	PeekBack() (T, bool)
}
