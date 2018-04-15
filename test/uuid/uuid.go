package main

import (
	"github.com/satori/go.uuid"
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	// Creating UUID Version 4
	// panic on error
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	// or error handling
	u2:= uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u2)

	// Parsing UUID from string input
	u3 := uuid.Must(uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	fmt.Printf("Successfully parsed: %s", u3)

	u4 := genUUID()
	if u4 != uuid.Nil {
		fmt.Printf("%+v\n", u4)
	}
}


func genUUID() (result uuid.UUID){
	defer func() {
		if err := recover(); err != nil {
			err, ok := err.(error)
			if ok {
				fmt.Printf("%+v\n", errors.WithStack(err))
			}
			result = uuid.Nil
		}
	}()
	result = uuid.NewV4()
	//panic(errors.New("xxx"))
	return result
}