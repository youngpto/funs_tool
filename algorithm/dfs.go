package algorithm

import (
	"github.com/youngpto/funs_tool/coll/queues/linkedlistqueue"
	"github.com/youngpto/funs_tool/coll/stacks/linkedliststack"
)

func DFS[T comparable](root T, fu func(pop T) []T) {
	stack := linkedliststack.New[T]()
	stack.Push(root)

	for !stack.IsEmpty() {
		pop, _ := stack.Pop()
		pushList := fu(pop)
		for _, t := range pushList {
			stack.Push(t)
		}
	}
}

func BFS[T comparable](root T, fu func(pop T) []T) {
	queue := linkedlistqueue.New[T]()
	queue.Push(root)

	for !queue.IsEmpty() {
		pop, _ := queue.Pop()
		pushList := fu(pop)
		for _, t := range pushList {
			queue.Push(t)
		}
	}
}
