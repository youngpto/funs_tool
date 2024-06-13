package rand_utils

import (
	"github.com/youngpto/funs_tool/times"
	"math/rand"
	"sync"
)

type Random interface {
	Seed(seed int64)
	Int63() int64
	Uint32() uint32
	Uint64() uint64
	Int31() int32
	Int() int
	Int63n(n int64) int64
	Int31n(n int32) int32
	Intn(n int) int
	Float64() float64
	Float32() float32
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
	Read(p []byte) (n int, err error)
	NormFloat64() float64
	ExpFloat64() float64
}

func New(seed int64, locked bool) Random {
	rng := rand.New(rand.NewSource(seed))
	if locked {
		return &lockedRandom{rng: rng}
	}
	return &normalRandom{rng: rng}
}

type normalRandom struct {
	rng *rand.Rand
}

func (nrd *normalRandom) Seed(seed int64) {
	nrd.rng.Seed(seed)
}

func (nrd *normalRandom) Int63() int64 {
	return nrd.rng.Int63()
}

func (nrd *normalRandom) Uint32() uint32 {
	return nrd.rng.Uint32()
}

func (nrd *normalRandom) Uint64() uint64 {
	return nrd.rng.Uint64()
}

func (nrd *normalRandom) Int31() int32 {
	return nrd.rng.Int31()
}

func (nrd *normalRandom) Int() int {
	return nrd.rng.Int()
}

func (nrd *normalRandom) Int63n(n int64) int64 {
	return nrd.rng.Int63n(n)
}

func (nrd *normalRandom) Int31n(n int32) int32 {
	return nrd.rng.Int31n(n)
}

func (nrd *normalRandom) Intn(n int) int {
	return nrd.rng.Intn(n)
}

func (nrd *normalRandom) Float64() float64 {
	return nrd.rng.Float64()
}

func (nrd *normalRandom) Float32() float32 {
	return nrd.rng.Float32()
}

func (nrd *normalRandom) Perm(n int) []int {
	return nrd.rng.Perm(n)
}

func (nrd *normalRandom) Shuffle(n int, swap func(i int, j int)) {
	nrd.rng.Shuffle(n, swap)
}

func (nrd *normalRandom) Read(p []byte) (n int, err error) {
	return nrd.rng.Read(p)
}

func (nrd *normalRandom) NormFloat64() float64 {
	return nrd.rng.NormFloat64()
}

func (nrd *normalRandom) ExpFloat64() float64 {
	return nrd.rng.ExpFloat64()
}

type lockedRandom struct {
	rng    *rand.Rand
	locker sync.Mutex
}

func (lrd *lockedRandom) Seed(seed int64) {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	lrd.rng.Seed(seed)
}

func (lrd *lockedRandom) Int63() int64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Int63()
}

func (lrd *lockedRandom) Uint32() uint32 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Uint32()
}

func (lrd *lockedRandom) Uint64() uint64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Uint64()
}

func (lrd *lockedRandom) Int31() int32 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Int31()
}

func (lrd *lockedRandom) Int() int {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Int()
}

func (lrd *lockedRandom) Int63n(n int64) int64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Int63n(n)
}

func (lrd *lockedRandom) Int31n(n int32) int32 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Int31n(n)
}

func (lrd *lockedRandom) Intn(n int) int {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Intn(n)
}

func (lrd *lockedRandom) Float64() float64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Float64()
}

func (lrd *lockedRandom) Float32() float32 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Float32()
}

func (lrd *lockedRandom) Perm(n int) []int {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Perm(n)
}

func (lrd *lockedRandom) Shuffle(n int, swap func(i int, j int)) {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	lrd.rng.Shuffle(n, swap)
}

func (lrd *lockedRandom) Read(p []byte) (n int, err error) {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.Read(p)
}

func (lrd *lockedRandom) NormFloat64() float64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.NormFloat64()
}

func (lrd *lockedRandom) ExpFloat64() float64 {
	lrd.locker.Lock()
	defer lrd.locker.Unlock()
	return lrd.rng.ExpFloat64()
}

var globalRand = New(times.Now().UnixNano(), true)

func Seed(seed int64) {
	globalRand.Seed(seed)
}

func Int63() int64 {
	return globalRand.Int63()
}

func Uint32() uint32 {
	return globalRand.Uint32()
}

func Uint64() uint64 {
	return globalRand.Uint64()
}

func Int31() int32 {
	return globalRand.Int31()
}

func Int() int {
	return globalRand.Int()
}

func Int63n(n int64) int64 {
	return globalRand.Int63n(n)
}

func Int31n(n int32) int32 {
	return globalRand.Int31n(n)
}

func Intn(n int) int {
	return globalRand.Intn(n)
}

func Float64() float64 {
	return globalRand.Float64()
}

func Float32() float32 {
	return globalRand.Float32()
}

func Perm(n int) []int {
	return globalRand.Perm(n)
}

func Shuffle(n int, swap func(i int, j int)) {
	globalRand.Shuffle(n, swap)
}

func Read(p []byte) (n int, err error) {
	return globalRand.Read(p)
}

func NormFloat64() float64 {
	return globalRand.NormFloat64()
}

func ExpFloat64() float64 {
	return globalRand.ExpFloat64()
}
