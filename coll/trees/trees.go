package trees

import "github.com/youngpto/funs_tool/coll"

type Tree[V any] interface {
	coll.Collection[V]
}
