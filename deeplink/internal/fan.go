package internal

import "sync"

// Fan is a fan-out channel.
type Fan[T any] struct {
	mutex     sync.Mutex
	listeners []chan T
}

// Send sends a value to all listeners.
func (f *Fan[T]) Send(v T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.listeners) == 0 {
		return
	}

	for _, c := range f.listeners {
		select {
		case c <- v:
		}
	}
}

// Add adds a listener to the fan-out channel.
func (f *Fan[T]) Add(n int) chan T {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	c := make(chan T, n)
	f.listeners = append(f.listeners, c)
	return c
}

// Close removes a listener from the fan-out channel.
func (f *Fan[T]) Close() {
	for _, c := range f.listeners {
		close(c)
	}
}
