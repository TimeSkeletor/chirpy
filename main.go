package main

import (
	"log"
	"net/http"
)


func main() {
	const filepathRoot = "."
	const port = "8080"
	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":" + port,
        Handler: mux,
    }

	log.Fatal(srv.ListenAndServe())
}
