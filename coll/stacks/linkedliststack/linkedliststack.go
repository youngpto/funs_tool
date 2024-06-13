package linkedliststack

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll/list/singlylinkedlist"
	"strings"
)

type Stack[T comparable] struct {
	list *singlylinkedlist.List[T]
}

func New[T comparable]() *Stack[T] {
	return &Stack[T]{list: singlylinkedlist.New[T]()}
}

func (stack *Stack[T]) String() string {
	str := "LinkedListStack\n"
	var values []string
	for _, value := range stack.list.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (stack *Stack[T]) IsEmpty() bool {
	return stack.list.IsEmpty()
}

func (stack *Stack[T]) Len() int {
	return stack.list.Len()
}

func (stack *Stack[T]) Clear() {
	stack.list.Clear()
}

func (stack *Stack[T]) Values() []T {
	return stack.list.Values()
}

func (stack *Stack[T]) Push(v T) {
	stack.list.Insert(0, v)
}

func (stack *Stack[T]) Pop() (T, bool) {
	value, ok := stack.list.Get(0)
	if ok {
		stack.list.Remove(0)
	}
	return value, ok
}

func (stack *Stack[T]) Peek() (T, bool) {
	return stack.list.Get(0)
}

func (stack *Stack[T]) withinRange(index int) bool {
	return index >= 0 && index < stack.list.Len()
}
