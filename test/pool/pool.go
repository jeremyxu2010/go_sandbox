package main

import (
	"sync"
	"fmt"
)

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return 0
		},
	}

	pool.Put(2)
	if v, ok := pool.Get().(int); ok {
		fmt.Printf("get value %d\n", v)
	}
	if v, ok := pool.Get().(int); ok {
		fmt.Printf("get value %d\n", v)
	}
}
