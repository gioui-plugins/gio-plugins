package internal

import (
	"crypto/rand"
	"sync"
	"sync/atomic"
)

var (
	key       = [32]byte{}
	handles   = sync.Map{} // map[Handle]interface{}
	handleIdx uintptr      // atomic
)

func init() {
	rand.Read(key[:])
}

type Handle uintptr

// NewHandle returns a handle for a given value.
func NewHandle(v any) Handle {
	h := atomic.AddUintptr(&handleIdx, 1)
	if h == 0 {
		panic("ran out of handle space")
	}

	handles.Store(h, v)
	return Handle(h)
}

// Value returns the associated Go value for a valid handle.
//
// The method panics if the handle is invalid.
func (h Handle) Value() any {
	v, ok := handles.Load(uintptr(h))
	if !ok {
		panic("misuse of an invalid Handle")
	}
	return v
}

// IsValid returns true if the handle is valid.
func (h Handle) IsValid() bool {
	_, ok := handles.Load(uintptr(h))
	return ok
}

// Delete invalidates a handle.
func (h Handle) Delete() {
	if _, ok := handles.LoadAndDelete(uintptr(h)); !ok {
		panic("misuse of an invalid Handle")
	}
}
