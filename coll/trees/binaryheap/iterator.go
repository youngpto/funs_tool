package binaryheap

type Iterator[T comparable] struct {
	heap  *Heap[T]
	index int
}

func (heap *Heap[T]) Iterator() *Iterator[T] {
	return &Iterator[T]{heap: heap, index: -1}
}

func (iterator *Iterator[T]) Next() bool {
	if iterator.index < iterator.heap.Len() {
		iterator.index++
	}
	return iterator.heap.withinRange(iterator.index)
}

func (iterator *Iterator[T]) Value() T {
	start, end := evaluateRange(iterator.index)
	if end > iterator.heap.Len() {
		end = iterator.heap.Len()
	}
	tmpHeap := NewWith(iterator.heap.Comparator)
	for n := start; n < end; n++ {
		value, _ := iterator.heap.list.Get(n)
		tmpHeap.Push(value)
	}
	for n := 0; n < iterator.index-start; n++ {
		tmpHeap.Pop()
	}
	value, _ := tmpHeap.Pop()
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
	return iterator.heap.withinRange(iterator.index)
}

func (iterator *Iterator[T]) End() {
	iterator.index = iterator.heap.Len()
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

func numOfBits(n int) uint {
	var count uint
	for n != 0 {
		count++
		n >>= 1
	}
	return count
}

func evaluateRange(index int) (start int, end int) {
	bits := numOfBits(index+1) - 1
	start = 1<<bits - 1
	end = start + 1<<bits
	return
}
