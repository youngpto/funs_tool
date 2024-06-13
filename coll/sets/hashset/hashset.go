package hashset

import (
	"fmt"
	"strings"
)

type Set[T comparable] struct {
	items map[T]struct{}
}

func New[T comparable](values ...T) *Set[T] {
	set := &Set[T]{items: make(map[T]struct{})}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func (set *Set[T]) String() string {
	str := "HashSet\n"
	var items []string
	for k := range set.items {
		items = append(items, fmt.Sprintf("%v", k))
	}
	str += strings.Join(items, ", ")
	return str
}

func (set *Set[T]) IsEmpty() bool {
	return set.Len() == 0
}

func (set *Set[T]) Len() int {
	return len(set.items)
}

func (set *Set[T]) Clear() {
	set.items = make(map[T]struct{})
}

func (set *Set[T]) Values() []T {
	values := make([]T, set.Len())
	count := 0
	for item := range set.items {
		values[count] = item
		count++
	}
	return values
}

func (set *Set[T]) Contains(value T) bool {
	_, contains := set.items[value]
	return contains
}

func (set *Set[T]) Add(elements ...T) {
	for _, item := range elements {
		set.items[item] = itemExists
	}
}

func (set *Set[T]) Remove(elements ...T) {
	for _, item := range elements {
		delete(set.items, item)
	}
}

func (set *Set[T]) Intersection(another *Set[T]) *Set[T] {
	result := New[T]()

	// Iterate over smaller set (optimization)
	if set.Len() <= another.Len() {
		for item := range set.items {
			if _, contains := another.items[item]; contains {
				result.Add(item)
			}
		}
	} else {
		for item := range another.items {
			if _, contains := set.items[item]; contains {
				result.Add(item)
			}
		}
	}

	return result
}

func (set *Set[T]) Union(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.items {
		result.Add(item)
	}
	for item := range another.items {
		result.Add(item)
	}

	return result
}

func (set *Set[T]) Difference(another *Set[T]) *Set[T] {
	result := New[T]()

	for item := range set.items {
		if _, contains := another.items[item]; !contains {
			result.Add(item)
		}
	}

	return result
}

var itemExists = struct{}{}
