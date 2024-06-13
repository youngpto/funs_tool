package coll

// Iterator 迭代接口
type Iterator[K, V any] interface {
	Next() bool
	Value() V
	Index() K
	Begin()
	First() bool
	NextTo(func(key K, value V) bool) bool
}

// ReverseIter 反向迭代接口
type ReverseIter[K, V any] interface {
	Iterator[K, V]

	Prev() bool
	End()
	Last() bool
	PrevTo(func(key K, value V) bool) bool
}

type RangeFunc[K, V any] func(key K, value V) (shouldContinue bool)

func Range[K, V any](iterator Iterator[K, V], rangeFunc RangeFunc[K, V]) {
	for iterator.Begin(); iterator.Next(); {
		if !rangeFunc(iterator.Index(), iterator.Value()) {
			return
		}
	}
}

func ReverseRange[K, V any](iterator ReverseIter[K, V], rangeFunc RangeFunc[K, V]) {
	for iterator.End(); iterator.Prev(); {
		if !rangeFunc(iterator.Index(), iterator.Value()) {
			return
		}
	}
}
