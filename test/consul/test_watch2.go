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
	if err := json.Unmarshal([]byte(`{"type":"keyprefix", "prefix":"test/"}`), &params); err != nil {
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
		kvs, ok := raw.(api.KVPairs)
		if !ok || len(kvs) == 0 {
			return
		}
		for _, kv := range kvs {
			fmt.Printf("key: %s, value: %+v\n", kv.Key, string(kv.Value))
		}
	}

	go func() {
		if err := plan.Run("127.0.0.1:8500"); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Minute * 5)

	plan.Stop()
}
