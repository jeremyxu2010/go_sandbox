package main

import (
	"log"
	"time"

	"personal/jeremyxu/sandbox/test/go-micro/proto/greeter"
	"github.com/micro/go-micro"
	"context"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

type say struct{}

func (s *say) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

var s = &say{}

func main() {
	service := micro.NewService(
		micro.Name("sandbox.test.go-micro.greeter"),
		micro.Registry(consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	greeter.RegisterGreeterHandler(service.Server(), s)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
