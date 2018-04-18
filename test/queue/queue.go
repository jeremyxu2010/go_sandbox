package main

import (
	"time"
	"github.com/pkg/errors"
	"fmt"
	"sync"
)

type Q struct {
	c chan interface{}
}

func (q *Q) Push(item interface{}, timeout time.Duration) (err error) {
	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()
	select {
	case q.c <- item:
		return nil
	case <-time.After(timeout):
		return errors.New("push timeouted")
	}
	return errors.New("never come here")
}

func (q *Q) Poll(timeout time.Duration) (item interface{}, err error) {
	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()
	select {
	case item = <- q.c:
		return item, nil
	case <- time.After(timeout):
		return nil, errors.New("Poll timeouted")
	}
	return nil, errors.New("never come here")
}

func NewQ(capacity int)*Q{
	q := Q{
		c: make(chan interface{}, capacity),
	}
	return &q
}

func main() {
	q := NewQ(100)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	go func() {
		for i := 0; i<1000; i++ {
			err := q.Push(i, time.Second*1)
			if err != nil {
				fmt.Printf("%+v", err)
			}
			time.Sleep(time.Second * 1)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i<1000; i++ {
			v, err := q.Poll(time.Second*1)
			if err != nil {
				fmt.Printf("%+v", err)
			}
			if v2, ok := v.(int); ok {
				fmt.Printf("%d\n", v2)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
