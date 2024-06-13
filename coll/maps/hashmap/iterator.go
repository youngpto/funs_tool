package hashmap

type Iterator[K comparable, V any] struct {
	hashMap *Map[K, V]
	keys    []K
	index   int
}

func (m *Map[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{
		hashMap: m,
		keys:    m.Keys(),
		index:   -1,
	}
}

func (iterator *Iterator[K, V]) Next() bool {
	if iterator.index < len(iterator.keys) {
		iterator.index++
	}
	if iterator.index < 0 || iterator.index >= len(iterator.keys) {
		return false
	}
	for !iterator.hashMap.Contains(iterator.keys[iterator.index]) {
		if iterator.index == len(iterator.keys)-1 {
			return false
		}
		iterator.index++
	}
	return true
}

func (iterator *Iterator[K, V]) Value() V {
	value, _ := iterator.hashMap.Get(iterator.keys[iterator.index])
	return value
}

func (iterator *Iterator[K, V]) Index() K {
	return iterator.keys[iterator.index]
}

func (iterator *Iterator[K, V]) Begin() {
	iterator.index = -1
}

func (iterator *Iterator[K, V]) First() bool {
	iterator.Begin()
	return iterator.Next()
}

func (iterator *Iterator[K, V]) NextTo(f func(key K, value V) bool) bool {
	for iterator.Next() {
		index, value := iterator.Index(), iterator.Value()
		if f(index, value) {
			return true
		}
	}
	return false
}
