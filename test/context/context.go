package main

import (
	"context"
	"time"
	"fmt"
)

func doSomething(ctx context.Context, done chan interface{}) {
	defer func() {
		done <- nil
	}()
	for {
		time.Sleep(time.Second * 3) // do something for a while
		select {
		case <-ctx.Done():
			fmt.Printf("goroutine is canceled, cancel reasion: %v\n", ctx.Err())
			return
		default:
		}
	}
}

func doSomething2(ctx context.Context, done chan interface{}) {
	defer func() {
		done <- nil
	}()
	if value, ok := ctx.Value("testKey").(string); ok {
		fmt.Printf("get testKey from context, value: %s\n", value)
	}
}

func main() {
	childCtx, cancel := context.WithCancel(context.Background())
	done := make(chan interface{})
	go doSomething(childCtx, done)
	time.Sleep(time.Second * 10) // cancel after 10 seconds
	cancel()
	<-done

	childCtx, _ = context.WithTimeout(context.Background(), time.Second*10) // auto cancel after 10 seconds
	done = make(chan interface{})
	go doSomething(childCtx, done)
	<-done

	childCtx = context.WithValue(context.Background(), "testKey", "testValue")
	done = make(chan interface{})
	go doSomething2(childCtx, done)
	<-done
}
