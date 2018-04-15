package main

import (
	"github.com/avast/retry-go"
	"time"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"errors"
)

func main() {
	client, err := consulapi.NewClient(&consulapi.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		panic(err)
	}
	kv := client.KV()

	//// PUT a new KV pair
	//kp := &consulapi.KVPair{Key:"test", Value:[]byte("xxxx")}
	//_, err = kv.Put(kp, nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Lookup the pair
	//pair, _, err := kv.Get("test", nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("KV: %v", pair.Value)
	lastIndex := uint64(0)
	retry.Do(
		func() error {
			pair, qm, err := kv.Get("xxx/1112", &consulapi.QueryOptions{WaitIndex: lastIndex, WaitTime: time.Second * 1000})
			if err != nil {
				return err
			}
			if pair == nil {
				return errors.New("xxxx")
			}
			if qm.LastIndex > lastIndex {
				fmt.Printf("KV: %v", string(pair.Value))
				println(qm.LastIndex)
				lastIndex = qm.LastIndex
			}
			return nil
		}, retry.Attempts(100), retry.Delay(3), retry.Units(time.Second), retry.OnRetry(func(n uint, err error) {
			fmt.Printf("the %d retry...\n", n)
		}))
}
