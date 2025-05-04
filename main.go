package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiConfig := apiConfig{} 

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics", apiConfig.getMetrics)
	mux.HandleFunc("/reset", apiConfig.resetMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
        Handler: mux,
    }

	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	 cfg.fileServerHits.Swap(0)
}