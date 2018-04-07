package main

import (
	"encoding/json"
	"bytes"
	"fmt"
	"os"
)

func main() {
	var b []byte
	buf := bytes.NewBuffer(b)
	encoder := json.NewEncoder(buf)
	encoder.Encode("vvdfdfd")

	decoder := json.NewDecoder(buf)
	var result string
	decoder.Decode(&result)
	fmt.Fprintf(os.Stderr, "%s", result)
}
