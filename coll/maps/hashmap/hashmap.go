package hashmap

import (
	"fmt"
	"github.com/youngpto/funs_tool/coll"
)

type Map[K comparable, V any] struct {
	m map[K]V
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{m: make(map[K]V)}
}

func (m *Map[K, V]) Put(key K, value V) {
	m.m[key] = value
}

func (m *Map[K, V]) Get(key K) (value V, found bool) {
	value, found = m.m[key]
	return
}

func (m *Map[K, V]) Contains(key K) (found bool) {
	_, found = m.m[key]
	return
}

func (m *Map[K, V]) SetDefault(key K, value V) V {
	v, found := m.Get(key)
	if !found {
		m.Put(key, value)
		v = value
	}
	return v
}

func (m *Map[K, V]) Remove(key K) {
	delete(m.m, key)
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) Len() int {
	return len(m.m)
}

func (m *Map[K, V]) Keys() []K {
	keys := make([]K, m.Len())
	count := 0
	for key := range m.m {
		keys[count] = key
		count++
	}
	return keys
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, m.Len())
	count := 0
	for _, value := range m.m {
		values[count] = value
		count++
	}
	return values
}

func (m *Map[K, V]) Clear() {
	m.m = make(map[K]V)
}

func (m *Map[K, V]) String() string {
	str := "HashMap\n"
	str += fmt.Sprintf("%v", m.m)
	return str
}

func (m *Map[K, V]) Range(rangeFunc coll.RangeFunc[K, V]) {
	for k, v := range m.m {
		if !rangeFunc(k, v) {
			return
		}
	}
}
