package main

import (
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"os"
	"fmt"
)

func echoServer(ws *websocket.Conn) {
	websocket.Message.Receive()
	defer ws.Close()
	io.Copy(ws, ws)
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
