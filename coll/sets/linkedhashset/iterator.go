package linkedhashset

import "github.com/youngpto/funs_tool/coll/list/doublylinkedlist"

type Iterator[T comparable] struct {
	iterator doublylinkedlist.Iterator[T]
}

func (set *Set[T]) Iterator() Iterator[T] {
	return Iterator[T]{iterator: set.ordering.Iterator()}
}

func (iterator *Iterator[T]) Next() bool {
	return iterator.iterator.Next()
}

func (iterator *Iterator[T]) Prev() bool {
	return iterator.iterator.Prev()
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

func (iterator *Iterator[T]) End() {
	iterator.iterator.End()
}

func (iterator *Iterator[T]) First() bool {
	return iterator.iterator.First()
}

func (iterator *Iterator[T]) Last() bool {
	return iterator.iterator.Last()
}

func (iterator *Iterator[T]) NextTo(f func(index int, value T) bool) bool {
	for iterator.Next() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}

func (iterator *Iterator[T]) PrevTo(f func(index int, value T) bool) bool {
	for iterator.Prev() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}
