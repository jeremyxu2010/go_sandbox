package main

import (
	"log"
	"time"

	"net/http"

	"github.com/micro/go-micro"
	httpServer "github.com/micro/go-plugins/server/http"
	"personal/jeremyxu/sandbox/test/go-micro/proto/greeter"
	"context"
	"io/ioutil"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
)

type say struct{}

func (s *say) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

var s = &say{}

func hello(w http.ResponseWriter, r *http.Request) {
	cf, ok := DefaultHTTPCodecs[r.Header.Get("Content-Type")]
	if !ok {
		log.Printf("can not find codec\n")
	}
	req := &greeter.Request{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	err = cf.Unmarshal(b, req)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	rsp := &greeter.Response{}
	err = s.Hello(context.Background(), req, rsp)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	b, err = cf.Marshal(rsp)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	w.Write(b)
}

func main() {
	srv := httpServer.NewServer()

	r := mux.NewRouter()
	r.Methods("POST").Path("/hello").HandlerFunc(hello)

	hd := srv.NewHandler(r)

	srv.Handle(hd)

	service := micro.NewService(
		micro.Server(srv),
		micro.Name("sandbox.test.go-micro.greeter"),
		micro.Registry(consul.NewRegistry(
			registry.Addrs("127.0.0.1:8500"),
		)),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type jsonCodec struct{}

type protoCodec struct{}

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(b []byte, v interface{}) error
	String() string
}

var (
	DefaultHTTPCodecs = map[string]Codec{
		"application/json":         jsonCodec{},
		"application/proto":        protoCodec{},
		"application/protobuf":     protoCodec{},
		"application/octet-stream": protoCodec{},
	}

)

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (protoCodec) String() string {
	return "proto"
}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) String() string {
	return "json"
}
