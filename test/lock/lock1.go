package main

import (
	"github.com/hashicorp/consul/api"
	"personal/jeremyxu/sandbox/test/lock/filelock"
	"os"
	"time"
	"fmt"
	"os/signal"
	"syscall"
	"personal/jeremyxu/sandbox/test/lock/election/client"
	"personal/jeremyxu/sandbox/test/lock/election"
)

func main() {
	config := api.DefaultConfig()                                  // Create a new api client config
	consulclient, _ := api.NewClient(config)                       // Create a Consul api client

	stopCh :=  make(chan struct{})

	leaderElection := &election.LeaderElection{
		StopElection:  stopCh,                        // The channel for stopping the election
		LeaderKey:     "service/my-service/leader", // The leadership key to create/aquire
		WatchWaitTime: 3,                                     // Time in seconds to check for leadership
		Client: &client.ConsulClient{Client:consulclient},     // The injected Consul api client
	}

	go leaderElection.ElectLeader()

	os.Create("/tmp/my-service.lock")

	flockScrambler := &filelock.FilelockScrambler{
		FileLocker: filelock.New("/tmp/my-service.lock"),
		StopCh: stopCh,
	}

	go flockScrambler.Scramble()

	go printLeaderStatus1(leaderElection, flockScrambler)

	errs := make(chan error)

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		err := fmt.Errorf("%s", <- ch)
		shutdown1(stopCh)
		errs <- err
	}()

	<- errs
}


func shutdown1(stopCh chan struct{}){
	defer close(stopCh)
	fmt.Printf("shutdown1...\n")
}

func printLeaderStatus1(le *election.LeaderElection, fs *filelock.FilelockScrambler) {
	for {
		if le.IsLeader() && fs.Locked {
			fmt.Println("lock1: I'm leader!")
		}
		time.Sleep(time.Second * 3)
	}
}
