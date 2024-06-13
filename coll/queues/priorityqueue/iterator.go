package priorityqueue

import "github.com/youngpto/funs_tool/coll/trees/binaryheap"

type Iterator[T comparable] struct {
	iterator *binaryheap.Iterator[T]
}

func (queue *Queue[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{iterator: queue.heap.Iterator()}
}

func (iterator *Iterator[T]) Next() bool {
	return iterator.iterator.Next()
}

func (iterator *Iterator[T]) Value() T {
	return iterator.iterator.Value()
}

func (iterator *Iterator[T]) Index() int {
	return iterator.iterator.Index()
}

func (iterator *Iterator[T]) Begin() {
	iterator.iterator.Begin()
}

func (iterator *Iterator[T]) First() bool {
	return iterator.iterator.First()
}

func (iterator *Iterator[T]) NextTo(f func(key int, value T) bool) bool {
	return iterator.iterator.NextTo(f)
}

func (iterator *Iterator[T]) Prev() bool {
	return iterator.iterator.Prev()
}

func (iterator *Iterator[T]) End() {
	iterator.iterator.End()
}

func (iterator *Iterator[T]) Last() bool {
	return iterator.iterator.Last()
}

func (iterator *Iterator[T]) PrevTo(f func(key int, value T) bool) bool {
	return iterator.iterator.PrevTo(f)
}
