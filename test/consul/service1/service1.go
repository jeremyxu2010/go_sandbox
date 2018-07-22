package main

import (
	"net/http"
	"log"
	"io"
)

func TestServer(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://127.0.0.1:38082/test2")
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		w.Write([]byte("make request failed\n"))
		return
	}
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/test1", TestServer)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}