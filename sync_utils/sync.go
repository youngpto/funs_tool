package sync_utils

import (
	"encoding/json"
	"fmt"
	"sync"
)

type RangeFunc[K, V any] func(key K, value V) (shouldContinue bool)

type RangeKeyFunc[K any] func(key K) (shouldContinue bool)

type RangeValueFunc[V any] func(value V) (shouldContinue bool)

type SyncMap[K, V any] struct {
	content sync.Map
}

func (syncMap *SyncMap[K, V]) Load(key K) (ret V, ok bool) {
	load, ok := syncMap.content.Load(key)
	if load != nil {
		return load.(V), ok
	}
	return
}

func (syncMap *SyncMap[K, V]) LoadBool(key K) bool {
	_, ok := syncMap.Load(key)
	return ok
}

func (syncMap *SyncMap[K, V]) LoadVal(key K) V {
	load, _ := syncMap.Load(key)
	return load
}

func (syncMap *SyncMap[K, V]) Store(key K, value V) {
	syncMap.content.Store(key, value)
}

func (syncMap *SyncMap[K, V]) LoadOrStore(key K, value V) (ret V, ok bool) {
	load, ok := syncMap.content.LoadOrStore(key, value)
	if load != nil {
		return load.(V), ok
	}
	return
}

func (syncMap *SyncMap[K, V]) LoadAndDelete(key K) (ret V, ok bool) {
	value, loaded := syncMap.content.LoadAndDelete(key)
	if value != nil {
		return value.(V), loaded
	}
	return
}

func (syncMap *SyncMap[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		syncMap.content.Delete(key)
	}
}

func (syncMap *SyncMap[K, V]) Range(rangeFunc RangeFunc[K, V]) {
	syncMap.content.Range(func(key, value any) bool {
		return rangeFunc(key.(K), value.(V))
	})
}

func (syncMap *SyncMap[K, V]) RangeKey(keyFunc RangeKeyFunc[K]) {
	syncMap.content.Range(func(key, value any) bool {
		return keyFunc(key.(K))
	})
}

func (syncMap *SyncMap[K, V]) RangeValue(valueFunc RangeValueFunc[V]) {
	syncMap.content.Range(func(key, value any) bool {
		return valueFunc(value.(V))
	})
}

func (syncMap *SyncMap[K, V]) Len() int {
	var i int
	syncMap.Range(func(key K, value V) (shouldContinue bool) {
		i++
		return true
	})
	return i
}

func (syncMap *SyncMap[K, V]) Clear() {
	syncMap.content = sync.Map{}
}

var SHARD_COUNT = 32

type Stringer interface {
	fmt.Stringer
	comparable
}

type ConcurrentMap[K comparable, V any] struct {
	shards   []*ConcurrentMapShared[K, V]
	sharding func(key K) uint32
}

type ConcurrentMapShared[K comparable, V any] struct {
	items        map[K]V
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

func create[K comparable, V any](sharding func(key K) uint32) ConcurrentMap[K, V] {
	m := ConcurrentMap[K, V]{
		sharding: sharding,
		shards:   make([]*ConcurrentMapShared[K, V], SHARD_COUNT),
	}
	for i := 0; i < SHARD_COUNT; i++ {
		m.shards[i] = &ConcurrentMapShared[K, V]{items: make(map[K]V)}
	}
	return m
}

func New[V any]() ConcurrentMap[string, V] {
	return create[string, V](fnv32)
}

func NewStringer[K Stringer, V any]() ConcurrentMap[K, V] {
	return create[K, V](strfnv32[K])
}

func NewWithCustomShardingFunction[K comparable, V any](sharding func(key K) uint32) ConcurrentMap[K, V] {
	return create[K, V](sharding)
}

func (m *ConcurrentMap[K, V]) GetShard(key K) *ConcurrentMapShared[K, V] {
	return m.shards[uint(m.sharding(key))%uint(SHARD_COUNT)]
}

func (m *ConcurrentMap[K, V]) MSet(data map[K]V) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCb[V any] func(exist bool, valueInMap V, newValue V) V

func (m *ConcurrentMap[K, V]) Upsert(key K, value V, cb UpsertCb[V]) (res V) {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

func (m *ConcurrentMap[K, V]) SetIfAbsent(key K, value V) bool {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m *ConcurrentMap[K, V]) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m *ConcurrentMap[K, V]) Has(key K) bool {
	shard := m.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (m *ConcurrentMap[K, V]) Remove(key K) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

type RemoveCb[K any, V any] func(key K, v V, exists bool) bool

func (m *ConcurrentMap[K, V]) RemoveCb(key K, cb RemoveCb[K, V]) bool {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return remove
}

func (m *ConcurrentMap[K, V]) Pop(key K) (v V, exists bool) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

func (m *ConcurrentMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

type Tuple[K comparable, V any] struct {
	Key K
	Val V
}

func (m *ConcurrentMap[K, V]) Iter() <-chan Tuple[K, V] {
	chans := snapshot(m)
	ch := make(chan Tuple[K, V])
	go fanIn(chans, ch)
	return ch
}

func (m *ConcurrentMap[K, V]) IterBuffered() <-chan Tuple[K, V] {
	chans := snapshot(m)
	total := 0
	for _, c := range chans {
		total += cap(c)
	}
	ch := make(chan Tuple[K, V], total)
	go fanIn(chans, ch)
	return ch
}

func (m *ConcurrentMap[K, V]) Clear() {
	for item := range m.IterBuffered() {
		m.Remove(item.Key)
	}
}

func snapshot[K comparable, V any](m *ConcurrentMap[K, V]) (chans []chan Tuple[K, V]) {
	if len(m.shards) == 0 {
		panic(`cmap.ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	chans = make([]chan Tuple[K, V], SHARD_COUNT)
	wg := sync.WaitGroup{}
	wg.Add(SHARD_COUNT)
	for index, shard := range m.shards {
		go func(index int, shard *ConcurrentMapShared[K, V]) {
			shard.RLock()
			chans[index] = make(chan Tuple[K, V], len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chans[index] <- Tuple[K, V]{key, val}
			}
			shard.RUnlock()
			close(chans[index])
		}(index, shard)
	}
	wg.Wait()
	return chans
}

func fanIn[K comparable, V any](chans []chan Tuple[K, V], out chan Tuple[K, V]) {
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan Tuple[K, V]) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

func (m *ConcurrentMap[K, V]) Items() map[K]V {
	tmp := make(map[K]V)

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

type IterCb[K comparable, V any] func(key K, v V)

func (m *ConcurrentMap[K, V]) IterCb(fn IterCb[K, V]) {
	for idx := range m.shards {
		shard := (m.shards)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

func (m *ConcurrentMap[K, V]) Keys() []K {
	count := m.Count()
	ch := make(chan K, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShared[K, V]) {
				// Foreach key, value pair.
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]K, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

func (m *ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	tmp := make(map[K]V)

	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

func strfnv32[K fmt.Stringer](key K) uint32 {
	return fnv32(key.String())
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m *ConcurrentMap[K, V]) UnmarshalJSON(b []byte) (err error) {
	tmp := make(map[K]V)

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	for key, val := range tmp {
		m.Set(key, val)
	}
	return nil
}
