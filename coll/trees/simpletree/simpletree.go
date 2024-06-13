package simpletree

import (
	"github.com/youngpto/funs_tool/algorithm"
	"github.com/youngpto/funs_tool/coll_utils"
	"runtime"
)

type Node[T any] struct {
	value    T
	parent   *Node[T]
	children []*Node[T]
	level    int // 节点层级
}

func New[T any](value T) *Node[T] {
	return AddChild(nil, value)
}

func AddChild[T any](parent *Node[T], value T) *Node[T] {
	child := &Node[T]{
		value:    value,
		parent:   parent,
		children: make([]*Node[T], 0),
	}
	if parent != nil {
		parent.children = append(parent.children, child)
		child.level = parent.level + 1
	}
	return child
}

func (node *Node[T]) AddChildren(values ...T) {
	for _, value := range values {
		AddChild(node, value)
	}
}

type LoopFunc[T any] func(fu func(pop *Node[T]) (pusher []*Node[T]))

func (node *Node[T]) BFS(fu func(pop *Node[T]) (pusher []*Node[T])) {
	algorithm.BFS(node, fu)
}

func (node *Node[T]) DFS(fu func(pop *Node[T]) (pusher []*Node[T])) {
	algorithm.DFS(node, fu)
}

func (node *Node[T]) Value() T {
	return node.value
}

func (node *Node[T]) Children() []*Node[T] {
	dst := make([]*Node[T], len(node.children), len(node.children))
	copy(dst, node.children)
	return dst
}

type FitterFunc[T any] func(node *Node[T]) bool

func WithLeafs[T any]() FitterFunc[T] {
	return func(node *Node[T]) bool {
		return len(node.children) == 0
	}
}

func WithLeLevel[T any](level int) FitterFunc[T] {
	return func(node *Node[T]) bool {
		return node.level >= 0 && node.level <= level
	}
}

func WithEqLevel[T any](level int) FitterFunc[T] {
	return func(node *Node[T]) bool {
		return node.level == level
	}
}

func WithGeLevel[T any](level int) FitterFunc[T] {
	return func(node *Node[T]) bool {
		return node.level >= level
	}
}

func WithEqValue[T comparable](value T) FitterFunc[T] {
	return func(node *Node[T]) bool {
		return node.value == value
	}
}

func (node *Node[T]) Fitter(lfu LoopFunc[T], fus ...FitterFunc[T]) []*Node[T] {
	values := make([]*Node[T], 0)
	lfu(func(pop *Node[T]) (pusher []*Node[T]) {
		pusher = pop.children
		for _, fu := range fus {
			if !fu(pop) {
				return
			}
		}
		values = append(values, pop)
		return

	})
	return values
}

func (node *Node[T]) Leafs() []*Node[T] {
	return node.Fitter(node.BFS, WithLeafs[T]())
}

func (node *Node[T]) ToSlice() []*Node[T] {
	return node.Fitter(node.BFS)
}

func (node *Node[T]) Parents() []*Node[T] {
	values := make([]*Node[T], 0)
	for ptr := node; ptr != nil; ptr = ptr.parent {
		values = append(values, ptr)
	}
	return coll_utils.Reverse(values)
}

func (node *Node[T]) rebuild() {
	var diffLevel int
	if node.parent == nil {
		diffLevel = -node.level
	} else {
		diffLevel = node.parent.level - node.level + 1
	}
	node.BFS(func(pop *Node[T]) (pusher []*Node[T]) {
		pop.level = pop.level + diffLevel
		return pop.children
	})
}

func (node *Node[T]) PopChild(child *Node[T]) (*Node[T], bool) {
	for idx, n := range node.children {
		if n == child {
			node.children = append(node.children[:idx], node.children[idx+1:]...)
			n.rebuild()
			return n, true
		}
	}
	return child, false
}

func (node *Node[T]) RemoveParent() *Node[T] {
	parent := node.parent
	if parent == nil {
		return nil
	}
	parent.PopChild(node)
	return parent
}

func (node *Node[T]) Move(npt *Node[T]) {
	if node.parent == npt {
		return
	}
	node.RemoveParent()
	if npt != nil {
		node.parent = npt
		node.rebuild()
	}
}

func (node *Node[T]) Clear() {
	node.parent = nil
	node.children = nil
	runtime.GC()
}
