package sync_utils

import "sync"

type SafeClose func()

type SafeObj[T any] struct {
	ReMutex
	data T
}

func NewSafeObj[T any](t T) *SafeObj[T] {
	return &SafeObj[T]{
		data: t,
	}
}

func (safe *SafeObj[T]) Get(f func(t T)) {
	WithGuard(safe, func() {
		f(safe.data)
	})
}

func (safe *SafeObj[T]) GetWithClose() (T, SafeClose) {
	safe.Lock()
	return safe.data, safe.Unlock
}

func (safe *SafeObj[T]) Set(t T) {
	WithGuard(safe, func() {
		safe.data = t
	})
}

func WithGuard(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}

func WithDefer(f, df func()) {
	defer df()
	f()
}
