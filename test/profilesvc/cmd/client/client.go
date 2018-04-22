// Package client provides a profilesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package main

import (
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	profilesvcapi "personal/jeremyxu/sandbox/test/profilesvc/api"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"os"
	"context"
	"fmt"
	"github.com/go-kit/kit/sd/lb"
	"github.com/fatih/structs"
	"reflect"
	consulapi  "github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd/consul"
)

const CONSUL_ADDR = "127.0.0.1:8500"
const SERVICE_NAME = "profilesvc"
const ENV_NAME = "prod"

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	srv, err := getRemoteProfileSvcAPIService(SERVICE_NAME, ENV_NAME, logger)
	if err != nil {
		panic(err)
	}
	err = srv.PostProfile(context.Background(), profilesvcapi.Profile{
		ID: "1",
		Name: "test1",
	})
	if err != nil {
		panic(err)
	}
	profile, err := srv.GetProfile(context.Background(), "1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", profile)
	profile.Name = "test2"
	err = srv.PutProfile(context.Background(), profile.ID, profile)
	if err != nil {
		panic(err)
	}
	profile, err = srv.GetProfile(context.Background(), "1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", profile)
	err = srv.DeleteProfile(context.Background(), "1")
	if err != nil {
		panic(err)
	}
}

// New returns a service that's load-balanced over instances of profilesvc found
// in the provided Consul server. The mechanism of looking up profilesvc
// instances in Consul is hard-coded into the client.
func getRemoteProfileSvcAPIService(serviceName string, envName string, logger log.Logger) (profilesvcapi.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: CONSUL_ADDR,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of profilesvc, we declare and enforce these
	// parameters for all of the profilesvc consumers.
	var (
		consulTags    = []string{envName}
		passingOnly   = true

	)

	sdclient  := consul.NewClient(apiclient)
	instancer := sd.Instancer(consul.NewInstancer(sdclient, logger, serviceName, consulTags, passingOnly))

	//instancer := sd.Instancer(sd.FixedInstancer{"192.168.1.201:8888"})

	var endpoints profilesvcapi.Endpoints

	//endpoints.PostProfileEndpoint = createClientEndpoint(instancer, profilesvcapi.MakePostProfileEndpoint, logger)
	//endpoints.GetProfileEndpoint = createClientEndpoint(instancer, profilesvcapi.MakeGetProfileEndpoint, logger)
	//endpoints.PutProfileEndpoint = createClientEndpoint(instancer, profilesvcapi.MakePutProfileEndpoint, logger)
	//endpoints.PatchProfileEndpoint = createClientEndpoint(instancer, profilesvcapi.MakePatchProfileEndpoint, logger)
	//endpoints.DeleteProfileEndpoint = createClientEndpoint(instancer, profilesvcapi.MakeDeleteProfileEndpoint, logger)
	//endpoints.GetAddressesEndpoint = createClientEndpoint(instancer, profilesvcapi.MakeGetAddressesEndpoint, logger)
	//endpoints.GetAddressEndpoint = createClientEndpoint(instancer, profilesvcapi.MakeGetAddressEndpoint, logger)
	//endpoints.PostAddressEndpoint = createClientEndpoint(instancer, profilesvcapi.MakePostAddressEndpoint, logger)
	//endpoints.DeleteAddressEndpoint = createClientEndpoint(instancer, profilesvcapi.MakeDeleteAddressEndpoint, logger)

	names := structs.Names(&endpoints)
	for _, name := range names {
		makerName := "Make" + name
		if maker, ok := profilesvcapi.EndpointMakerMap[makerName]; ok {
			e := createClientEndpoint(instancer, maker, logger)
			v := reflect.ValueOf(e)
			reflect.ValueOf(&endpoints).Elem().FieldByName(name).Set(v)
		}

	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(profilesvcapi.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := profilesvcapi.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}

func createClientEndpoint(instancer sd.Instancer, endpointMaker func(s profilesvcapi.Service) endpoint.Endpoint, logger log.Logger) endpoint.Endpoint{
	var (
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)
	factory := factoryFor(endpointMaker)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, retryTimeout, balancer)
	return retry
}


