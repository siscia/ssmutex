package ssmutex_test

import (
	"fmt"

	"github.com/siscia/ssmutex"
)

func Example_leak() {
	state := 0
	lock := ssmutex.NewContainer(&state)

	var leakedPtr *int

	i, key := lock.Lock()
	leakedPtr = i
	key.Unlock()

	(*leakedPtr)++
	fmt.Println("New shared state =", *leakedPtr)
	// Output:
	// New shared state = 1
}
