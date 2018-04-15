package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
)

func main() {
	resp, err := http.Get("http://www.baidu.com")
	checkErr(err)
	fmt.Println(resp.Status)
	for key, value := range resp.Header {
		fmt.Printf("%s, %v\n", key, value)
	}

	io.Copy(os.Stdout, resp.Body)

	req, err := http.NewRequest("HEAD", "http://www.baidu.com", nil)

	resp, err = http.DefaultClient.Do(req)

	resp, err = (&http.Client{}).Do(req)
}
func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
