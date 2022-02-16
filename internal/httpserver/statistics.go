package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nokka/d2-armory-api/internal/domain"
)

// statisticsService encapsulates the business logic around statistics.
type statisticsService interface {
	// Gets the character statistics.
	GetCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error)

	// Parse parses a character binary.
	Parse(ctx context.Context, stats []domain.StatisticsRequest) error

	// DeleteStats deletes all stats by the given character.
	DeleteStats(ctx context.Context, character string) error
}

// statisticsHandler is used to put parse statistics requests.
type statisticsHandler struct {
	encoder           *encoder
	statisticsService statisticsService
	credentials       map[string]string
}

func (h statisticsHandler) Routes(router chi.Router) {
	// Posting and deleting statistics requires authentication.
	router.With(middleware.BasicAuth("statistics", h.credentials)).Post("/", h.postStatistics)
	router.With(middleware.BasicAuth("statistics", h.credentials)).Delete("/{name}", h.deleteStatistics)

	// Get statistics by character.
	router.Get("/", h.getStatistics)
}

func (h statisticsHandler) postStatistics(w http.ResponseWriter, r *http.Request) {
	var stats []domain.StatisticsRequest
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		h.encoder.Error(w, err)
		return
	}

	// Pass the request context in order to make use of cancellation for lower level work.
	err := h.statisticsService.Parse(r.Context(), stats)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, map[string]string{"status": "accepted"}, http.StatusAccepted)
}
func (h statisticsHandler) deleteStatistics(w http.ResponseWriter, r *http.Request) {
	characterName := chi.URLParam(r, "name")

	// Pass the request context in order to make use of cancellation for lower level work.
	err := h.statisticsService.DeleteStats(r.Context(), characterName)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h statisticsHandler) getStatistics(w http.ResponseWriter, r *http.Request) {
	characterName := r.URL.Query().Get("character")

	//Pass the request context in order to make use of cancellation for lower level work.
	stats, err := h.statisticsService.GetCharacter(r.Context(), characterName)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.Response(w, stats)
}

func newStatisticsHandler(encoder *encoder, statisticsService statisticsService, credentials map[string]string) *statisticsHandler {
	return &statisticsHandler{
		encoder:           encoder,
		statisticsService: statisticsService,
		credentials:       credentials,
	}
}
