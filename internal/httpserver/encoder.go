package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nokka/d2-armory-api/internal/domain"
)

type encoder struct{}

func newEncoder() *encoder {
	return &encoder{}
}

// Response will determine Content-Type and encode the response properly.
func (e *encoder) Response(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}

// errorResponse will encapsulate errors to be transferred over http.
type errorResponse struct {
	Error string `json:"error"`
}

// encodeError will determine status code and content sent over the API.
func (e *encoder) Error(w http.ResponseWriter, err error) {
	resp := errorResponse{
		Error: err.Error(),
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	if errors.Is(err, domain.ErrTemporary) {
		w.Header().Set("x-temporary", "true")
	}

	switch errors.Unwrap(err) {
	case domain.ErrRequest:
		w.WriteHeader(http.StatusBadRequest)
	case domain.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case domain.ErrUnavailable:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(resp)
}
