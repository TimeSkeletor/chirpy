package main

import "net/http"

const (
	platform = "dev"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != platform {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
