package arraylist

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll"
	"sort"
	"strings"
)

type List[T comparable] struct {
	elements []T
	size     int
}

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

func New[T comparable](values ...T) *List[T] {
	list := &List[T]{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (l *List[T]) Get(index int) (T, bool) {
	if !l.withinRange(index) {
		var t T
		return t, false
	}

	return l.elements[index], true
}

func (l *List[T]) Remove(index int) {
	if !l.withinRange(index) {
		return
	}

	l.elements = append(l.elements[:index], l.elements[index+1:l.size]...)
	l.size--

	l.shrink()
}

func (l *List[T]) Add(values ...T) {
	l.growBy(len(values))
	for _, value := range values {
		l.elements[l.size] = value
		l.size++
	}
}

func (l *List[T]) Contains(value T) bool {
	for index := 0; index < l.size; index++ {
		if l.elements[index] == value {
			return true
		}
	}
	return false
}

func (l *List[T]) Sort(comparator coll.Comparator[T]) {
	if len(l.elements) < 2 {
		return
	}

	sort.Slice(l.elements[:l.size], func(i, j int) bool {
		return comparator(l.elements[i], l.elements[j]) < 0
	})
}

func (l *List[T]) Swap(index1, index2 int) {
	if l.withinRange(index1) && l.withinRange(index2) {
		l.elements[index1], l.elements[index2] = l.elements[index2], l.elements[index1]
	}
}

func (l *List[T]) Insert(index int, values ...T) {
	if !l.withinRange(index) {
		// Append
		if index == l.size {
			l.Add(values...)
		}
		return
	}

	length := len(values)
	l.growBy(length)
	l.size += length
	copy(l.elements[index+length:], l.elements[index:l.size-length])
	copy(l.elements[index:], values)
}

func (l *List[T]) Set(index int, value T) {
	if !l.withinRange(index) {
		// Append
		if index == l.size {
			l.Add(value)
		}
		return
	}

	l.elements[index] = value
}

func (l *List[T]) String() string {
	str := "ArrayList\n"
	values := make([]string, 0, l.size)
	for _, value := range l.elements[:l.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

func (l *List[T]) IsEmpty() bool {
	return l.size == 0
}

func (l *List[T]) Len() int {
	return l.size
}

func (l *List[T]) Clear() {
	l.size = 0
	l.elements = []T{}
}

func (l *List[T]) Values() []T {
	newElements := make([]T, l.size, l.size)
	copy(newElements, l.elements[:l.size])
	return newElements
}

func (l *List[T]) withinRange(index int) bool {
	return index >= 0 && index < l.size
}

func (l *List[T]) resize(cap int) {
	newElements := make([]T, cap, cap)
	copy(newElements, l.elements)
	l.elements = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (l *List[T]) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(l.elements)
	if l.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		l.resize(newCapacity)
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (l *List[T]) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(l.elements)
	if l.size <= int(float32(currentCapacity)*shrinkFactor) {
		l.resize(l.size)
	}
}
