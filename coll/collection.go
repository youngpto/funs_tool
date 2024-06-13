package coll

import (
	"fmt"
)

// Collection 容器接口
type Collection[T any] interface {
	fmt.Stringer

	IsEmpty() bool
	Len() int
	Clear()
	Values() []T
}

type Contains[T any] interface {
	Contains(value T) bool
}

func Any[T any](contains Contains[T], values ...T) bool {
	if len(values) == 0 {
		return true
	}
	for _, value := range values {
		if contains.Contains(value) {
			return true
		}
	}
	return false
}

func All[T any](contains Contains[T], values ...T) bool {
	if len(values) == 0 {
		return true
	}
	for _, value := range values {
		if !contains.Contains(value) {
			return false
		}
	}
	return true
}

func IsEmpty[T any](collection Collection[T]) bool {
	return collection.IsEmpty()
}

func Len[T any](collection Collection[T]) int {
	return collection.Len()
}

func Clear[T any](collection Collection[T]) {
	collection.Clear()
}

func Values[T any](collection Collection[T]) []T {
	return collection.Values()
}

func ReverseValues[T any](collection Collection[T]) []T {
	values := collection.Values()
	if len(values) == 0 {
		return values
	}
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
	return values
}
