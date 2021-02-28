package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Setup our http server we want to test on.
	srv := NewServer(":80", nil, nil, nil, true)

	// Setup a new test recorder.
	recorder := httptest.NewRecorder()

	// Create a new request to the specific handler.
	req := httptest.NewRequest("GET", "/health", nil)

	// Add the recorder and the request to our server, to serve it.
	srv.Handler().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("want status 200, got = %d", recorder.Code)
	}
}
