package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncoderResponse(t *testing.T) {
	e := &encoder{}

	w := httptest.NewRecorder()

	e.Response(w, struct{ Foo int }{123})

	if got, want := w.Header().Get("Content-Type"), "application/json; charset=utf-8"; got != want {
		t.Fatalf(`w.Header().Get("Content-Type") = %q, want %q`, got, want)
	}

	var r struct{ Foo int }

	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := r.Foo, 123; got != want {
		t.Fatalf("r.Foo = %d, want %d", got, want)
	}
}

func TestEncoderError(t *testing.T) {
	e := &encoder{}

	t.Run("different status codes for known errors", func(t *testing.T) {
		for _, tt := range []struct {
			err  error
			want int
		}{
			{errors.New("something went terribly wrong"), http.StatusInternalServerError},
		} {
			tt := tt

			w := httptest.NewRecorder()

			e.Error(w, tt.err)

			if got := w.Code; got != tt.want {
				t.Fatalf("w.Code = %d, want %d", got, tt.want)
			}
		}
	})
}
