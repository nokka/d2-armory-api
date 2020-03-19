package httpserver

import (
	"net/http"

	"github.com/go-chi/chi"
)

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

	char, err := h.characterService.Parse(name)
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
