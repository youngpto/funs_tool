package doublylinkedlist

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll"
	"sort"
	"strings"
)

type List[T comparable] struct {
	first *element[T]
	last  *element[T]
	size  int
}

type element[T comparable] struct {
	value T
	prev  *element[T]
	next  *element[T]
}

func New[T comparable](values ...T) *List[T] {
	list := &List[T]{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (list *List[T]) Append(values ...T) {
	list.Add(values...)
}

func (list *List[T]) Prepend(values ...T) {
	for v := len(values) - 1; v >= 0; v-- {
		newElement := &element[T]{value: values[v], next: list.first}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.first.prev = newElement
			list.first = newElement
		}
		list.size++
	}
}

func (list *List[T]) Get(index int) (T, bool) {
	if !list.withinRange(index) {
		var t T
		return t, false
	}

	// determine traveral direction, last to first or first to last
	if list.size-index < index {
		element := list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
		return element.value, true
	}
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
	}
	return element.value, true
}

func (list *List[T]) Remove(index int) {
	if !list.withinRange(index) {
		return
	}

	if list.size == 1 {
		list.Clear()
		return
	}

	var element *element[T]
	// determine traversal direction, last to first or first to last
	if list.size-index < index {
		element = list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
	} else {
		element = list.first
		for e := 0; e != index; e, element = e+1, element.next {
		}
	}

	if element == list.first {
		list.first = element.next
	}
	if element == list.last {
		list.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	element = nil

	list.size--
}

func (list *List[T]) Add(values ...T) {
	for _, value := range values {
		newElement := &element[T]{value: value, prev: list.last}
		if list.size == 0 {
			list.first = newElement
			list.last = newElement
		} else {
			list.last.next = newElement
			list.last = newElement
		}
		list.size++
	}
}

func (list *List[T]) Contains(value T) bool {
	if list.size == 0 {
		return false
	}
	for element := list.first; element != nil; element = element.next {
		if element.value == value {
			return true
		}
	}
	return false
}

func (list *List[T]) Sort(comparator coll.Comparator[T]) {
	if list.size < 2 {
		return
	}

	values := list.Values()
	sort.Slice(values, func(i, j int) bool {
		return comparator(values[i], values[j]) < 0
	})

	list.Clear()

	list.Add(values...)

}

func (list *List[T]) Swap(index1, index2 int) {
	if list.withinRange(index1) && list.withinRange(index2) && index1 != index2 {
		var element1, element2 *element[T]
		for e, currentElement := 0, list.first; element1 == nil || element2 == nil; e, currentElement = e+1, currentElement.next {
			switch e {
			case index1:
				element1 = currentElement
			case index2:
				element2 = currentElement
			}
		}
		element1.value, element2.value = element2.value, element1.value
	}
}

func (list *List[T]) Insert(index int, values ...T) {
	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.Add(values...)
		}
		return
	}

	var beforeElement *element[T]
	var foundElement *element[T]
	// determine traversal direction, last to first or first to last
	if list.size-index < index {
		foundElement = list.last
		beforeElement = list.last.prev
		for e := list.size - 1; e != index; e, foundElement = e-1, foundElement.prev {
			beforeElement = beforeElement.prev
		}
	} else {
		foundElement = list.first
		for e := 0; e != index; e, foundElement = e+1, foundElement.next {
			beforeElement = foundElement
		}
	}

	if foundElement == list.first {
		oldNextElement := list.first
		for i, value := range values {
			newElement := &element[T]{value: value}
			if i == 0 {
				list.first = newElement
			} else {
				newElement.prev = beforeElement
				beforeElement.next = newElement
			}
			beforeElement = newElement
		}
		oldNextElement.prev = beforeElement
		beforeElement.next = oldNextElement
	} else {
		oldNextElement := beforeElement.next
		for _, value := range values {
			newElement := &element[T]{value: value}
			newElement.prev = beforeElement
			beforeElement.next = newElement
			beforeElement = newElement
		}
		oldNextElement.prev = beforeElement
		beforeElement.next = oldNextElement
	}

	list.size += len(values)
}

func (list *List[T]) Set(index int, value T) {
	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.Add(value)
		}
		return
	}

	var foundElement *element[T]
	// determine traversal direction, last to first or first to last
	if list.size-index < index {
		foundElement = list.last
		for e := list.size - 1; e != index; {
			fmt.Println("Set last", index, value, foundElement, foundElement.prev)
			e, foundElement = e-1, foundElement.prev
		}
	} else {
		foundElement = list.first
		for e := 0; e != index; {
			e, foundElement = e+1, foundElement.next
		}
	}

	foundElement.value = value
}

func (list *List[T]) String() string {
	str := "DoublyLinkedList\n"
	var values []string
	for element := list.first; element != nil; element = element.next {
		values = append(values, fmt.Sprintf("%v", element.value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (list *List[T]) IsEmpty() bool {
	return list.size == 0
}

func (list *List[T]) Len() int {
	return list.size
}

func (list *List[T]) Clear() {
	list.size = 0
	list.first = nil
	list.last = nil
}

func (list *List[T]) Values() []T {
	values := make([]T, list.size, list.size)
	for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
		values[e] = element.value
	}
	return values
}

func (list *List[T]) IndexOf(value T) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.Values() {
		if element == value {
			return index
		}
	}
	return -1
}

func (list *List[T]) withinRange(index int) bool {
	return index >= 0 && index < list.size
}
