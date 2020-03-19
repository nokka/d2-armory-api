package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// healthHandler is used to do health probes to verify that the service is up and running.
type healthHandler struct{}

func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Routes(router chi.Router) {
	router.Get("/", h.get)
}

func (h *healthHandler) get(w http.ResponseWriter, r *http.Request) {
	ret := map[string]interface{}{
		"status": "OK",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	_ = json.NewEncoder(w).Encode(ret)
}
