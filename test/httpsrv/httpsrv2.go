package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

// see https://github.com/gorilla/mux
func main() {
	r := mux.NewRouter()
	r.Methods("OPTIONS").PathPrefix("/somepath").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test\n"))
	})
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler)

	//// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":8080", r))
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin"}), handlers.AllowedMethods([]string{"GET", "HEAD", "POST"}))(r)))
}
