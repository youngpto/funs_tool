package linkedlistdeque

type Iterator[T comparable] struct {
	deque *Deque[T]
	index int
}

func (deque *Deque[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{deque: deque, index: -1}
}

func (iterator *Iterator[T]) Next() bool {
	if iterator.index < iterator.deque.Len() {
		iterator.index++
	}
	return iterator.deque.withinRange(iterator.index)
}

func (iterator *Iterator[T]) Value() T {
	value, _ := iterator.deque.list.Get(iterator.index)
	return value
}

func (iterator *Iterator[T]) Index() int {
	return iterator.index
}

func (iterator *Iterator[T]) Begin() {
	iterator.index = -1
}

func (iterator *Iterator[T]) First() bool {
	iterator.Begin()
	return iterator.Next()
}

func (iterator *Iterator[T]) NextTo(f func(key int, value T) bool) bool {
	for iterator.Next() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}

func (iterator *Iterator[T]) Prev() bool {
	if iterator.index >= 0 {
		iterator.index--
	}
	return iterator.deque.withinRange(iterator.index)
}

func (iterator *Iterator[T]) End() {
	iterator.index = iterator.deque.Len()
}

func (iterator *Iterator[T]) Last() bool {
	iterator.End()
	return iterator.Prev()
}

func (iterator *Iterator[T]) PrevTo(f func(key int, value T) bool) bool {
	for iterator.Prev() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}
