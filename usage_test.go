package ssmutex_test

import (
	"fmt"

	"github.com/siscia/ssmutex"
)

type StructureWithComplexSharedState struct {
	sharedState *ssmutex.Container[*int]
}

func NewStructureWithComplexSharedState(n int) *StructureWithComplexSharedState {
	return &StructureWithComplexSharedState{
		sharedState: ssmutex.NewContainer(&n),
	}
}

func (s *StructureWithComplexSharedState) SimpleWayToAdd1() {
	s.sharedState.WithLock(func(i *int) {
		(*i)++
	})
}

func callerDoesNotAcquireTheLock(sharedState *ssmutex.Container[*int]) {
	i, key := sharedState.Lock()
	defer key.Unlock()
	*i++
}

func callerHasTheLock(i *int) {
	(*i)++
}

func (s *StructureWithComplexSharedState) complexWayToAdd2() {
	callerDoesNotAcquireTheLock(s.sharedState)
	i, key := s.sharedState.Lock()
	defer key.Unlock()
	callerHasTheLock(i)
}

func Example() {
	s := NewStructureWithComplexSharedState(0)
	s.SimpleWayToAdd1()
	s.complexWayToAdd2()

	s.sharedState.WithLock(func(i *int) {
		fmt.Printf("Complex Shared State Value = %d", *i)
	})
	// Output:
	// Complex Shared State Value = 3
}
