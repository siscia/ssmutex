// ssmutex is a simple implementation for more ergonomic, and simpler to not misuse Locks in Go.
// The main idea behind the package is making the Lock a **container** for the value that is being locked.
// Making the lock a container allow to type check function boundaries.
//
// While working with shared state, it is common to have functions that take as input some shared state.
// It is important to know when the lock should be acquired and by whom.
// Should the lock be acquired before calling the function by the caller, or the function itself should
// acquire the lock? Using locks as containers allow to immediately answer this question, by looking only to the
// function declaration.
// If the argument of the function is a lock-container, the function itself need to acquire the lock (to get access
// to the shared state.)
// If the argument of the function is directly the shared state, it means that the lock was already acquired.
//
// Using locks as container allows forcing the correct acquisition of locks.
// However, it does nothing to force the correct release of locks.
// Unfortunately I was not able to figure out a completely satisfactory answer to the problem, however the
// current interface provides a reasonable solution.
// The Lock() methods returns the object being locked AND a struct (a key) to use to unlock the mutex.
// Only by using the key it is possible to unlock the mutex.
// In order to obtain the value a second time, you are force to use the key once.
// Each key can be only used once.
// If a key is leaked (and the GC run) the software panics.
// Enforcing the usage of the key could be done at LINT time, as we already do with context.Context and the
// cancel function.
//
// The package does not prohibit the leak of the shared state.
// However, it should be immediately evident.
package ssmutex

import (
	"runtime"
	"sync"
)

type Container[T any] struct {
	v T
	l sync.Mutex
}

type Unlocker interface {
	Unlock()
}

type key struct {
	l      *sync.Mutex
	once   sync.Once
	runned bool
}

func (k *key) Unlock() {
	k.once.Do(func() {
		k.l.Unlock()
		k.runned = true
	})
}

// Create a new container mutex holding the value to protect
func NewContainer[T any](v T) *Container[T] {
	return &Container[T]{
		v: v,
		l: sync.Mutex{},
	}
}

// Lock the container and return the shared state together with the key to unlock it again.
// If the key is leaked without being used, it forces a panic.
func (m *Container[T]) Lock() (T, Unlocker) {
	m.l.Lock()
	k := &key{l: &m.l}
	runtime.SetFinalizer(k, func(fk *key) {
		if !fk.runned {
			panic("ssmutex: unlocked mutex leaked")
		}
	})
	return m.v, k
}

// Convenient wrapper for locking, execute a function, and unlocking.
func (m *Container[T]) WithLock(f func(T)) {
	v, key := m.Lock()
	f(v)
	key.Unlock()
}
