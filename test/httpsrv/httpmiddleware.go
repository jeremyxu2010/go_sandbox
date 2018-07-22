package main

import (
	"net/http"
	"io"
	"path/filepath"
	"fmt"
)

func ServeHTTP(writer http.ResponseWriter, req *http.Request){
	io.WriteString(writer, "xxxx")
}

type Middleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

type handler struct {
	middlewares []Middleware
	processFunc http.HandlerFunc
}

func (h *handler)ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	targetProcessFunc := h.processFunc
	for i:=len(h.middlewares); i>0; i-- {
		m := h.middlewares[i-1]
		targetProcessFunc = m(targetProcessFunc)
	}
	targetProcessFunc(w, req)
}

func (h *handler)Use(m Middleware){
	h.middlewares = append(h.middlewares, m)
}

func DemoMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, req *http.Request){

		fmt.Println("do something in middleware")
		//do something in middleware

		handlerFunc(w, req)
	}
}

func newHandler()http.Handler {
	h := handler{
		processFunc: func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("testxxx"))
		},
	}
	h.Use(DemoMiddleware)
	return &h
}

func main() {
	path, _ := filepath.Abs(".")
	http.Handle("/", http.FileServer(http.Dir(path)))

	h := newHandler()
	http.Handle("/test", h)
	http.ListenAndServe(":2345", nil)
}
