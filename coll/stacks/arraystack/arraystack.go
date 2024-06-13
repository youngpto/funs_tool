package arraystack

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll/list/arraylist"
	"strings"
)

type Stack[T comparable] struct {
	list *arraylist.List[T]
}

func New[T comparable]() *Stack[T] {
	return &Stack[T]{list: arraylist.New[T]()}
}

func (stack *Stack[T]) String() string {
	str := "ArrayStack\n"
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

func (stack *Stack[T]) Push(value T) {
	stack.list.Add(value)
}

func (stack *Stack[T]) Pop() (value T, ok bool) {
	pop, ok := stack.list.Get(stack.list.Len() - 1)
	if ok {
		stack.list.Remove(stack.list.Len() - 1)
	}
	return pop, ok
}

func (stack *Stack[T]) Peek() (value T, ok bool) {
	return stack.list.Get(stack.list.Len() - 1)
}

func (stack *Stack[T]) withinRange(index int) bool {
	return index >= 0 && index < stack.list.Len()
}
