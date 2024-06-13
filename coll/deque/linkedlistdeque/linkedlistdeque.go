package linkedlistdeque

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll/list/doublylinkedlist"
	"strings"
)

type Deque[T comparable] struct {
	list *doublylinkedlist.List[T]
}

func New[T comparable]() *Deque[T] {
	return &Deque[T]{
		list: doublylinkedlist.New[T](),
	}
}

func (deque *Deque[T]) String() string {
	str := "LinkedListDeque\n"
	var values []string
	for _, value := range deque.list.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (deque *Deque[T]) IsEmpty() bool {
	return deque.list.IsEmpty()
}

func (deque *Deque[T]) Len() int {
	return deque.list.Len()
}

func (deque *Deque[T]) Clear() {
	deque.list.Clear()
}

func (deque *Deque[T]) Values() []T {
	return deque.list.Values()
}

func (deque *Deque[T]) Push(v T) {
	deque.list.Add(v)
}

func (deque *Deque[T]) Pop() (T, bool) {
	value, ok := deque.list.Get(0)
	if ok {
		deque.list.Remove(0)
	}
	return value, ok
}

func (deque *Deque[T]) Peek() (T, bool) {
	return deque.list.Get(0)
}

func (deque *Deque[T]) PushFront(v T) {
	deque.list.Insert(0, v)
}

func (deque *Deque[T]) PopBack() (T, bool) {
	pop, ok := deque.list.Get(deque.list.Len() - 1)
	if ok {
		deque.list.Remove(deque.list.Len() - 1)
	}
	return pop, ok
}

func (deque *Deque[T]) PeekBack() (T, bool) {
	return deque.list.Get(deque.list.Len() - 1)
}

func (deque *Deque[T]) withinRange(index int) bool {
	return index >= 0 && index < deque.list.Len()
}
