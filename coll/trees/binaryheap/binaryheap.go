package binaryheap

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll"
	"github.com/youngpto/funs_tool/coll/list/arraylist"
	"strings"
)

type Heap[T comparable] struct {
	list       *arraylist.List[T]
	Comparator coll.Comparator[T]
}

func New[T coll.Ordered]() *Heap[T] {
	return &Heap[T]{list: arraylist.New[T](), Comparator: coll.Cmp[T]}
}

func NewWith[T comparable](comparator coll.Comparator[T]) *Heap[T] {
	return &Heap[T]{list: arraylist.New[T](), Comparator: comparator}
}

func (heap *Heap[T]) Init() {
	size := heap.list.Len()/2 + 1
	for i := size; i >= 0; i-- {
		heap.bubbleDownIndex(i)
	}
}

func (heap *Heap[T]) Push(values ...T) {
	if len(values) == 1 {
		heap.list.Add(values[0])
		heap.bubbleUp()
	} else {
		for _, value := range values {
			heap.list.Add(value)
		}
		heap.Init()
	}
}

func (heap *Heap[T]) Pop() (value T, ok bool) {
	value, ok = heap.list.Get(0)
	if !ok {
		return
	}
	lastIndex := heap.list.Len() - 1
	heap.list.Swap(0, lastIndex)
	heap.list.Remove(lastIndex)
	heap.bubbleDown()
	return
}

func (heap *Heap[T]) Peek() (value T, ok bool) {
	return heap.list.Get(0)
}

func (heap *Heap[T]) String() string {
	str := "BinaryHeap\n"
	var values []string
	for it := heap.Iterator(); it.Next(); {
		values = append(values, fmt.Sprintf("%v", it.Value()))
	}
	str += strings.Join(values, ", ")
	return str
}

func (heap *Heap[T]) IsEmpty() bool {
	return heap.list.IsEmpty()
}

func (heap *Heap[T]) Len() int {
	return heap.list.Len()
}

func (heap *Heap[T]) Clear() {
	heap.list.Clear()
}

func (heap *Heap[T]) Values() []T {
	values := make([]T, heap.list.Len(), heap.list.Len())
	for it := heap.Iterator(); it.Next(); {
		values[it.Index()] = it.Value()
	}
	return values
}

func (heap *Heap[T]) bubbleDown() {
	heap.bubbleDownIndex(0)
}

func (heap *Heap[T]) bubbleDownIndex(index int) {
	size := heap.list.Len()
	for leftIndex := index<<1 + 1; leftIndex < size; leftIndex = index<<1 + 1 {
		rightIndex := index<<1 + 2
		smallerIndex := leftIndex
		leftValue, _ := heap.list.Get(leftIndex)
		rightValue, _ := heap.list.Get(rightIndex)
		if rightIndex < size && heap.Comparator(leftValue, rightValue) > 0 {
			smallerIndex = rightIndex
		}
		indexValue, _ := heap.list.Get(index)
		smallerValue, _ := heap.list.Get(smallerIndex)
		if heap.Comparator(indexValue, smallerValue) > 0 {
			heap.list.Swap(index, smallerIndex)
		} else {
			break
		}
		index = smallerIndex
	}
}

func (heap *Heap[T]) bubbleUp() {
	index := heap.list.Len() - 1
	for parentIndex := (index - 1) >> 1; index > 0; parentIndex = (index - 1) >> 1 {
		indexValue, _ := heap.list.Get(index)
		parentValue, _ := heap.list.Get(parentIndex)
		if heap.Comparator(parentValue, indexValue) <= 0 {
			break
		}
		heap.list.Swap(index, parentIndex)
		index = parentIndex
	}
}

func (heap *Heap[T]) withinRange(index int) bool {
	return index >= 0 && index < heap.list.Len()
}
