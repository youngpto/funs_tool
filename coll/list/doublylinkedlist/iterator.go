package doublylinkedlist

type Iterator[T comparable] struct {
	list    *List[T]
	index   int
	element *element[T]
}

func (list *List[T]) Iterator() Iterator[T] {
	return Iterator[T]{list: list, index: -1, element: nil}
}

func (iterator *Iterator[T]) Next() bool {
	if iterator.index < iterator.list.size {
		iterator.index++
	}
	if !iterator.list.withinRange(iterator.index) {
		iterator.element = nil
		return false
	}
	if iterator.index != 0 {
		iterator.element = iterator.element.next
	} else {
		iterator.element = iterator.list.first
	}
	return true
}

func (iterator *Iterator[T]) Value() T {
	return iterator.element.value
}

func (iterator *Iterator[T]) Index() int {
	return iterator.index
}

func (iterator *Iterator[T]) Begin() {
	iterator.index = -1
	iterator.element = nil
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
	if !iterator.list.withinRange(iterator.index) {
		iterator.element = nil
		return false
	}
	if iterator.index == iterator.list.size-1 {
		iterator.element = iterator.list.last
	} else {
		iterator.element = iterator.element.prev
	}
	return iterator.list.withinRange(iterator.index)
}

func (iterator *Iterator[T]) End() {
	iterator.index = iterator.list.size
	iterator.element = iterator.list.last
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
