package main

import (
	"encoding/base64"
	"bytes"
	"os"
	"fmt"
)

func main() {
	orig := []byte{1, 3, 5, 7}

	buf := bytes.NewBuffer([]byte{})
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	_, err := encoder.Write(orig)
	checkError(err)
	encoder.Close()

	fmt.Println(buf)

	decoder := base64.NewDecoder(base64.StdEncoding, buf)
	var result = make([]byte, 256)
	n, err2 := decoder.Read(result)
	checkError(err2)
	for _, ch := range result[0:n] {
		fmt.Printf("%x\n", ch)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(orig))
	result, err2 = base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(orig))
	fmt.Println(result)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
