package arrayqueue

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll/list/arraylist"
	"strings"
)

type Queue[T comparable] struct {
	list *arraylist.List[T]
}

func New[T comparable]() *Queue[T] {
	return &Queue[T]{list: arraylist.New[T]()}
}

func (queue *Queue[T]) String() string {
	str := "ArrayQueue\n"
	var values []string
	for _, value := range queue.list.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (queue *Queue[T]) IsEmpty() bool {
	return queue.list.IsEmpty()
}

func (queue *Queue[T]) Len() int {
	return queue.list.Len()
}

func (queue *Queue[T]) Clear() {
	queue.list.Clear()
}

func (queue *Queue[T]) Values() []T {
	return queue.list.Values()
}

func (queue *Queue[T]) Push(v T) {
	queue.list.Add(v)
}

func (queue *Queue[T]) Pop() (T, bool) {
	pop, ok := queue.list.Get(0)
	if ok {
		queue.list.Remove(0)
	}
	return pop, ok
}

func (queue *Queue[T]) Peek() (T, bool) {
	return queue.list.Get(0)
}

func (queue *Queue[T]) withinRange(index int) bool {
	return index >= 0 && index < queue.list.Len()
}
