package list

import (
	"github.com/youngpto/funs_tool/coll"
)

type List[T comparable] interface {
	coll.Collection[T]
	coll.Contains[T]

	Get(index int) (T, bool)
	Remove(index int)
	Add(values ...T)
	Sort(comparator coll.Comparator[T])
	Swap(index1, index2 int)
	Insert(index int, values ...T)
	Set(index int, value T)
}
