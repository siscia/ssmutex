package ssmutex_test

import (
	"testing"

	"github.com/siscia/ssmutex"
)

func Test_LockAndUnlock(t *testing.T) {
	state := 0
	s := ssmutex.NewContainer(&state)

	i, key := s.Lock()
	(*i)++
	key.Unlock()

	s.WithLock(func(i *int) {
		if *i != 1 {
			t.Fatalf("expected 1 got %d", *i)
		}
	})
}

func Test_WithLock(t *testing.T) {
	state := 0
	s := ssmutex.NewContainer(&state)

	s.WithLock(func(i *int) {
		*i++
	})
}

// func Test_NotUnlockingRaisePanic(t *testing.T) {

// 	defer func() {
// 		recover()
// 	}()

// 	state := 0
// 	s := ssmutex.NewMutex(&state)

// 	v, _ := s.Lock()
// 	*v++
// 	runtime.GC()

// 	// it should panic
// }
