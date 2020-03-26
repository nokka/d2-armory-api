package httpserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nokka/d2-armory-api/internal/domain"
)

// characterService represents the functionality we need to perform our character requests.
type characterService interface {
	// Parse parses a character binary.
	Parse(ctx context.Context, name string) (*domain.Character, error)
}

// characterHandler is used to put parse characters.
type characterHandler struct {
	encoder          *encoder
	characterService characterService
}

func (h characterHandler) Routes(router chi.Router) {
	router.Get("/", h.parseCharacter)
}

func (h characterHandler) parseCharacter(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	// Pass the request context in order to make use of cancellation for lower level work.
	char, err := h.characterService.Parse(r.Context(), name)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.Response(w, char)
}

func newCharacterHandler(encoder *encoder, characterService characterService) *characterHandler {
	return &characterHandler{
		encoder:          encoder,
		characterService: characterService,
	}
}
