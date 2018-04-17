package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"fmt"
	"io"
)

func echoServer(ws *websocket.Conn) {
	defer ws.Close()
	io.Copy(os.Stdout, ws)
}
func main() {
	http.Handle("/", websocket.Handler(echoServer))
	err := http.ListenAndServe(":12345", nil)
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
