package main

import (
	"fmt"
	"os"
	"github.com/pkg/errors"
)

func main() {
	//err := errors.New("xxxx")
	//fmt.Printf("%+v\n", err)

	_, err := os.OpenFile("/tmp/test.txt",os.O_RDONLY, 0)
	if err != nil {
		err := errors.WithStack(err)
		fmt.Printf("%+v\n", err)
	}
}
