package maps

import "github.com/youngpto/funs_tool/coll"

// Map interface that all maps implement
type Map[K comparable, V any] interface {
	coll.Collection[V]

	Put(key K, value V)
	Get(key K) (value V, found bool)
	Contains(key K) (found bool)
	SetDefault(key K, value V) V
	Remove(key K)
	Keys() []K
}

// BidiMap interface that all bidirectional maps implement (extends the Map interface)
type BidiMap[K comparable, V comparable] interface {
	Map[K, V]

	GetKey(value V) (key K, found bool)
}
