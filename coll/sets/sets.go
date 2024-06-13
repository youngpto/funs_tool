package sets

import "github.com/youngpto/funs_tool/coll"

type Set[T comparable] interface {
	coll.Collection[T]
	coll.Contains[T]

	Add(elements ...T)
	Remove(elements ...T)
}
