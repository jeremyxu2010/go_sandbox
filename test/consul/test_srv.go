package main

import (
	consulapi "github.com/hashicorp/consul/api"
	"fmt"
	"github.com/satori/go.uuid"
)

var (
	SrvNS = uuid.Must(uuid.FromString("7ee2a51a-4002-11e8-989a-b8aeed7d9c97"))
)

func main() {
	client, err := consulapi.NewClient(&consulapi.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		panic(err)
	}

	srvUUID := uuid.NewV3(SrvNS, "192.168.1.201:9527")
	registration := consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("demoSrv_%s", srvUUID.String()),
		Name:    "demoSrv",
		Address: "127.0.0.1",
		Port:    9527,
		Tags:    []string{"proj1", "cluster1", "dev"},
		Check: &consulapi.AgentServiceCheck{
			//Args:                           []string{"sh", "-c", "sleep 1 && exit 0"},
			HTTP:                           "http://192.168.1.201:9527/check",
			Timeout:                        "3s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "1m", //check失败后30秒删除本服务
		},
	}

	agent := client.Agent()

	err = agent.ServiceRegister(&registration)
	if err != nil {
		panic(err)
	}

	services, err := agent.Services()
	if err != nil {
		panic(err)
	}
	for srvId := range services {
		fmt.Printf("name: %s, agent: %v", srvId, services[srvId])
	}
}
