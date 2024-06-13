package coll

import "sort"

// Comparator 定义比较算法类型
type Comparator[T any] func(x, y T) int

type ComparatorSort[T any] interface {
	Sorted[T]

	Compare(x, y T) int
}

type Sorted[T any] interface {
	Sort(comparator Comparator[T])
}

func GetSortedValueFunc[T any](collection Collection[T], comparator Comparator[T]) []T {
	values := collection.Values()
	if len(values) < 2 {
		return values
	}
	sort.Slice(values, func(i, j int) bool {
		return comparator(values[i], values[j]) < 0
	})
	return values
}

func SortFunc[T any](sorted Sorted[T], comparator Comparator[T]) {
	sorted.Sort(comparator)
}

func Sort[T any](imp ComparatorSort[T]) {
	imp.Sort(imp.Compare)
}

// Ordered 代表所有可比大小排序的类型
type Ordered interface {
	Integer | Float | ~string
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type DefaultComparatorImp[T Ordered] struct{}

func (d *DefaultComparatorImp[T]) Compare(x, y T) int {
	return Cmp[T](x, y)
}

func Cmp[T Ordered](x, y T) int {
	if x == y {
		return 0
	}
	if x < y {
		return -1
	}
	return 1
}
