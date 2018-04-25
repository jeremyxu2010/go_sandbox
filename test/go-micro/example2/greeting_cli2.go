package main

import (
	"context"
	"fmt"

	"personal/jeremyxu/sandbox/test/go-micro/proto/greeter"
	"github.com/micro/go-plugins/client/http"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
)

func main() {
	// create a new service
	service := micro.NewService(
		micro.Client(http.NewClient(
			client.Registry(consul.NewRegistry(
				registry.Addrs("127.0.0.1:8500"),
			)),
			client.Selector(selector.NewSelector(
				selector.SetStrategy(selector.RoundRobin),
			)),
		)),
	)

	// parse command line flags
	service.Init()

	cl := service.Client()

	req := cl.NewRequest("sandbox.test.go-micro.greeter", "/hello", &greeter.Request{
		Name: "John",
	})
	rsp := &greeter.Response{}

	err := cl.Call(context.Background(), req, rsp)


	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
