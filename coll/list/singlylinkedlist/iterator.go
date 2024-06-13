package singlylinkedlist

type Iterator[T comparable] struct {
	list    *List[T]
	index   int
	element *element[T]
}

func (list *List[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{list: list, index: -1, element: nil}
}

func (iterator *Iterator[T]) Next() bool {
	if iterator.index < iterator.list.size {
		iterator.index++
	}
	if !iterator.list.withinRange(iterator.index) {
		iterator.element = nil
		return false
	}
	if iterator.index == 0 {
		iterator.element = iterator.list.first
	} else {
		iterator.element = iterator.element.next
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
