package main

import (
	"golang.org/x/net/websocket"
	"os"
	"fmt"
	"io"
	"net"
)

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}
	config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s", "127.0.0.1:12345", "/"), "http://127.0.0.1")
	conn, err := websocket.NewClient(config, client)
	if err != nil {
		panic(err)
	}
	io.Copy(conn, os.Stdin)
}
