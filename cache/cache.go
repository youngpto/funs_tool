package utils

import (
	"container/list"
	"github.com/youngpto/funs_tool/sync_utils"
)

type PopCallback[K comparable, V any] func(key K, val V)

type LRUCache[K comparable, V any] struct {
	capacity int
	stack    *list.List
	items    map[K]*list.Element
	lock     sync_utils.ReMutex
	popCb    PopCallback[K, V]
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewLRUCache[K comparable, V any](capacity int, popCb PopCallback[K, V]) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		stack:    list.New(),
		items:    make(map[K]*list.Element),
		popCb:    popCb,
	}
}

func (lru *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if elem, ok := lru.items[key]; ok {
		lru.stack.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, ok
	}
	return
}

func (lru *LRUCache[K, V]) MultipleGet(keys []K) (values []V, miss []K) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	values = make([]V, 0, len(keys))
	miss = make([]K, 0, len(keys))
	for _, key := range keys {
		if elem, ok := lru.items[key]; ok {
			lru.stack.MoveToFront(elem)
			values = append(values, elem.Value.(*entry[K, V]).value)
		} else {
			miss = append(miss, key)
		}
	}
	return
}

func (lru *LRUCache[K, V]) Peek(key K) (value V) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if elem, ok := lru.items[key]; ok {
		return elem.Value.(*entry[K, V]).value
	}
	return
}

func (lru *LRUCache[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	for cursor := lru.stack.Front(); cursor != nil; cursor = cursor.Next() {
		elem := cursor.Value.(*entry[K, V])
		if !f(elem.key, elem.value) {
			break
		}
	}
}

func (lru *LRUCache[K, V]) Contains(key K) bool {
	_, ok := lru.Get(key)
	return ok
}

func (lru *LRUCache[K, V]) Set(key K, value V) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if elem, ok := lru.items[key]; ok {
		lru.stack.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	elem := &entry[K, V]{
		key:   key,
		value: value,
	}
	lru.items[key] = lru.stack.PushFront(elem)
	if lru.capacity < 0 {
		return
	}
	if lru.stack.Len() > lru.capacity {
		lru.removeElement(lru.stack.Back())
	}
}

func (lru *LRUCache[K, V]) Remove(key K) bool {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if elem, ok := lru.items[key]; ok {
		lru.removeElement(elem)
		return true
	}
	return false
}

func (lru *LRUCache[K, V]) removeElement(element *list.Element) {
	lru.stack.Remove(element)
	elem := element.Value.(*entry[K, V])
	delete(lru.items, elem.key)

	if lru.popCb != nil {
		lru.popCb(elem.key, elem.value)
	}
}

func (lru *LRUCache[K, V]) ReCapacity(capacity int) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if capacity < 0 {
		lru.capacity = -1
		return
	}
	if capacity > lru.capacity && lru.capacity >= 0 {
		lru.capacity = capacity
		return
	}

	lru.capacity = capacity
	step := lru.stack.Len() - capacity
	for i := 0; i < step; i++ {
		lru.removeElement(lru.stack.Back())
	}
}

func (lru *LRUCache[K, V]) Size() int {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	return lru.stack.Len()
}

func (lru *LRUCache[K, V]) Capacity() int {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	return lru.capacity
}

func (lru *LRUCache[K, V]) FlushAll() {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if lru.popCb != nil {
		for key, value := range lru.items {
			lru.popCb(key, value.Value.(*entry[K, V]).value)
		}
	}
	lru.items = make(map[K]*list.Element)
	lru.stack.Init()
}
