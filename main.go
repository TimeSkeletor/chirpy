package main

import (
	"net/http"
)


type Server struct {}


func main() {
	mux := http.NewServeMux()

	server := http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

	server.ListenAndServe()
}