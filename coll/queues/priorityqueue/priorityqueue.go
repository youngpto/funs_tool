package priorityqueue

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll"
	"github.com/youngpto/funs_tool/coll/trees/binaryheap"
	"strings"
)

type Queue[T comparable] struct {
	heap       *binaryheap.Heap[T]
	Comparator coll.Comparator[T]
}

func New[T coll.Ordered]() *Queue[T] {
	return NewWith[T](coll.Cmp[T])
}

func NewWith[T comparable](comparator coll.Comparator[T]) *Queue[T] {
	return &Queue[T]{heap: binaryheap.NewWith(comparator), Comparator: comparator}
}

func (queue *Queue[T]) String() string {
	str := "PriorityQueue\n"
	values := make([]string, queue.heap.Len(), queue.heap.Len())
	for index, value := range queue.heap.Values() {
		values[index] = fmt.Sprintf("%v", value)
	}
	str += strings.Join(values, ", ")
	return str
}

func (queue *Queue[T]) IsEmpty() bool {
	return queue.heap.IsEmpty()
}

func (queue *Queue[T]) Len() int {
	return queue.heap.Len()
}

func (queue *Queue[T]) Clear() {
	queue.heap.Clear()
}

func (queue *Queue[T]) Values() []T {
	return queue.heap.Values()
}

func (queue *Queue[T]) PushBack(v T) {
	queue.heap.Push(v)
}

func (queue *Queue[T]) PopFront() (T, bool) {
	return queue.heap.Pop()
}

func (queue *Queue[T]) PeekFront() (T, bool) {
	return queue.heap.Peek()
}

func (queue *Queue[T]) Sort() {
	queue.heap.Init()
}
