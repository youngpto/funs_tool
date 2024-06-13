package sync_utils

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type ReMutex struct {
	mu        sync.Mutex
	recursion int32
	owner     int64
}

func getGid() int64 {
	buf := [64]byte{}
	s := buf[:runtime.Stack(buf[:], false)]
	s = s[len("goroutine "):]
	s = s[:bytes.IndexAny(s, " ")]
	gid, _ := strconv.ParseInt(string(s), 10, 64)

	return gid
}

func (r *ReMutex) Lock() {
	if atomic.LoadInt64(&r.owner) == getGid() {
		r.recursion++
		return
	}

	r.mu.Lock()
	atomic.StoreInt64(&r.owner, getGid())
	r.recursion = 1
	return
}

func (r *ReMutex) Unlock() {
	if atomic.LoadInt64(&r.owner) != getGid() {
		panic(fmt.Sprintf("invalid onwer(%d) to unlock onwer(%d)`s lock", getGid(), r.owner))
	}

	r.recursion--
	if r.recursion != 0 {
		return
	}
	atomic.StoreInt64(&r.owner, -1)
	r.mu.Unlock()
}
