package main

import (
	"net/http"
	"io"
	"path/filepath"
)

func ServeHTTP(writer http.ResponseWriter, req *http.Request){
	io.WriteString(writer, "xxxx")
}

func main() {
	path, _ := filepath.Abs(".")
	http.Handle("/", http.FileServer(http.Dir(path)))

	http.HandleFunc("/test", ServeHTTP)
	http.ListenAndServe(":2345", nil)
}
