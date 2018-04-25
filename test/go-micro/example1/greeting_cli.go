package main

import (
	"context"
	"fmt"

	"personal/jeremyxu/sandbox/test/go-micro/proto/greeter"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
)

func main() {
	// create a new service
	service := micro.NewService(
		micro.Registry(consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))),
		micro.Selector(selector.NewSelector(selector.Option(selector.SetStrategy(selector.RoundRobin)))),
	)

	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := greeter.NewGreeterService("sandbox.test.go-micro.greeter", service.Client())

	// Make request
	rsp, err := cl.Hello(context.Background(), &greeter.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
