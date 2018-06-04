package main

import (
	"fmt"
	"bytes"
)

func main() {
	b := []byte{239, 187, 191, 97}
	b = bytes.TrimPrefix(b, []byte{239, 187, 191})
	str := string(b)
	fmt.Println(str)
}