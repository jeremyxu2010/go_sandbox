package main

import (
	"os/signal"
	"os"
	"syscall"
	"fmt"
)

func main() {
	errs := make(chan error)
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		err := fmt.Errorf("%s", <- ch)
		shutdown()
		errs <- err
	}()

	fmt.Printf("exit: %v", <- errs)

}


func shutdown(){
	fmt.Printf("do something...\n")
}
