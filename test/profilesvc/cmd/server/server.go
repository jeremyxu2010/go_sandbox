package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	profilesvcapi "personal/jeremyxu/sandbox/test/profilesvc/api"
	profilesvcimpl "personal/jeremyxu/sandbox/test/profilesvc/impl"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"personal/jeremyxu/sandbox/test/profilesvc/utils"
	"github.com/satori/go.uuid"
	"strconv"
	"errors"
)



const CONSUL_ADDR = "127.0.0.1:8500"
const SERVICE_NAME = "profilesvc"
const ENV_NAME = "prod"
const SERVICE_PORT = 8888

var (
	SrvNS = uuid.Must(uuid.FromString("7ee2a51a-4002-11e8-989a-b8aeed7d9c97"))
)

func main() {
	httpAddr := ":" + strconv.Itoa(SERVICE_PORT)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestamp)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	stdClient, err := consulapi.NewClient(&consulapi.Config{
		Address: CONSUL_ADDR,
	})
	if err != nil {
		panic(err)
	}

	client := consul.NewClient(stdClient)

	advertiseAddress := utils.GetNoLoopbackAddr()
	if advertiseAddress == "" {
		panic(errors.New("can not get not loopback interface address"))
	}

	srvUUID := uuid.NewV3(SrvNS, advertiseAddress  + ":" + strconv.Itoa(SERVICE_PORT))
	srvID := fmt.Sprintf(SERVICE_NAME + "_%s", srvUUID.String())
	// Produce a fake service registration.
	r := &consulapi.AgentServiceRegistration{
		ID:                srvID,
		Name:              SERVICE_NAME,
		Tags:              []string{ENV_NAME},
		Port:              SERVICE_PORT,
		Address:           advertiseAddress,
		EnableTagOverride: false,
		Check: &consulapi.AgentServiceCheck{
			TCP: advertiseAddress  + ":" + strconv.Itoa(SERVICE_PORT),
			Timeout:                        "2s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	// Build a registrar for r.
	registrar := consul.NewRegistrar(client, r, log.With(logger, "component", "registrar"))

	var s profilesvcapi.Service
	{
		s = profilesvcimpl.NewInmemService()
		s = profilesvcimpl.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = profilesvcimpl.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err := fmt.Errorf("%s", <-c)
		registrar.Deregister()
		errs <- err
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", httpAddr)
		registrar.Register()
		errs <- http.ListenAndServe(httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
