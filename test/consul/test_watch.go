package main

import (
	"github.com/hashicorp/consul/watch"
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"fmt"
	"time"
)

func main() {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(`{"type":"key", "key":"test"}`), &params); err != nil {
		panic(err)
	}
	plan, err := watch.Parse(params)
	if err != nil {
		panic(err)
	}
	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			return // ignore
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			return // ignore
		}
		fmt.Printf("%+v\n", string(v.Value))
	}

	go func() {
		if err := plan.Run("127.0.0.1:8500"); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Minute * 5)

	plan.Stop()
}
