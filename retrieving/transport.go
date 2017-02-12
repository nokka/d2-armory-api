package retrieving

import (
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// MakeHandler returns a handler for the character service.
func MakeHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	retrieveCharacterHandler := kithttp.NewServer(
		ctx,
		makeRetrieveCharacterEndpoint(s),
		decodeRetrieveCharacterRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/retrieving/v1/character", retrieveCharacterHandler).Methods("GET")

	return r
}

func decodeRetrieveCharacterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")

	return retrieveCharacterRequest{
		CharacterName: name,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic layer
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case ErrNonExistingCharacter:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
