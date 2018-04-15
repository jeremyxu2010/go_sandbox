package main

import (
	"fmt"
	"os"
	"time"
	"github.com/pkg/errors"
)

func doSomething(done chan <- interface{}) {
	defer func() {
		done <- nil
	}()
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Fprintf(os.Stderr, "goroutine dead, error: %+v\n", err)
		}
	}()
	panic(errors.New("test error"))
	time.Sleep(time.Second * 5)
}

func main() {
	done := make(chan interface{})
	go doSomething(done)
	<- done
}
