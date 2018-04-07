package main

import (
	"encoding/asn1"
	"fmt"
	"os"
	"time"
)

func main() {
	bytes, err := asn1.Marshal("abcd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshall error")
	}
	var result string
	_, err1 := asn1.Unmarshal(bytes, &result)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "unmarshall error")
	}

	fmt.Fprintf(os.Stderr, "the result is %s", result)

	t := time.Now()
	bytes2, _ := asn1.Marshal(&t)
	var result2 = new(time.Time)
	asn1.Unmarshal(bytes2, &result2)
	fmt.Fprintf(os.Stderr, "the result is %v", *result2)
}
