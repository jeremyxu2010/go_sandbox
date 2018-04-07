package main

import (
	"context"
	"fmt"

	hello "github.com/micro/examples/greeter/srv/proto/hello"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/selector"
)

func main() {
	// create a new service
	service := micro.NewService(
		micro.Selector(selector.NewSelector(selector.Option(selector.SetStrategy(selector.RoundRobin)))),
	)

	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := hello.NewSayClient("go.micro.srv.greeter", service.Client())

	// Make request
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Msg)
}
