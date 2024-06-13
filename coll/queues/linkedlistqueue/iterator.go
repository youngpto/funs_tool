package linkedlistqueue

type Iterator[T comparable] struct {
	queue *Queue[T]
	index int
}

func (queue *Queue[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{queue: queue, index: -1}
}

func (iterator *Iterator[T]) Next() bool {
	if iterator.index < iterator.queue.Len() {
		iterator.index++
	}
	return iterator.queue.withinRange(iterator.index)
}

func (iterator *Iterator[T]) Value() T {
	value, _ := iterator.queue.list.Get(iterator.index)
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
