package linkedhashset

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll/list/doublylinkedlist"
	"strings"
)

type Set[T comparable] struct {
	table    map[T]struct{}
	ordering *doublylinkedlist.List[T]
}

func New[T comparable](values ...T) *Set[T] {
	set := &Set[T]{
		table:    make(map[T]struct{}),
		ordering: doublylinkedlist.New[T](),
	}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func (set *Set[T]) String() string {
	str := "LinkedHashSet\n"
	var items []string
	it := set.Iterator()
	for it.Next() {
		items = append(items, fmt.Sprintf("%v", it.Value()))
	}
	str += strings.Join(items, ", ")
	return str
}

func (set *Set[T]) IsEmpty() bool {
	return set.Len() == 0
}

func (set *Set[T]) Len() int {
	return set.ordering.Len()
}

func (set *Set[T]) Clear() {
	set.table = make(map[T]struct{})
	set.ordering.Clear()
}

func (set *Set[T]) Values() []T {
	values := make([]T, set.Len())
	it := set.Iterator()
	for it.Next() {
		values[it.Index()] = it.Value()
	}
	return values
}

func (set *Set[T]) Contains(value T) bool {
	_, contains := set.table[value]
	return contains
}

func (set *Set[T]) Add(elements ...T) {
	for _, item := range elements {
		if _, contains := set.table[item]; !contains {
			set.table[item] = itemExists
			set.ordering.Append(item)
		}
	}
}

func (set *Set[T]) Remove(elements ...T) {
	for _, item := range elements {
		if _, contains := set.table[item]; contains {
			delete(set.table, item)
			index := set.ordering.IndexOf(item)
			set.ordering.Remove(index)
		}
	}
}

func (set *Set[T]) Intersection(another *Set[T]) *Set[T] {
	result := New[T]()

	// Iterate over smaller set (optimization)
	if set.Len() <= another.Len() {
		for item := range set.table {
			if _, contains := another.table[item]; contains {
				result.Add(item)
			}
		}
	} else {
		for item := range another.table {
			if _, contains := set.table[item]; contains {
				result.Add(item)
			}
		}
	}

	return result
}

func (set *Set[T]) Union(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.table {
		result.Add(item)
	}
	for item := range another.table {
		result.Add(item)
	}

	return result
}

func (set *Set[T]) Difference(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.table {
		if _, contains := another.table[item]; !contains {
			result.Add(item)
		}
	}

	return result
}

var itemExists = struct{}{}
