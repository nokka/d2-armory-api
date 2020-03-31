package httpserver

import (
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server is the HTTP server listener.
type Server struct {
	encoder          *encoder
	listener         net.Listener
	addr             string
	characterService characterService
}

// Open will open a tcp listener to serve http requests.
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = ln

	// Create an http server.
	server := http.Server{
		Handler: s.Handler(),
	}

	log.Println("starting HTTP server on:", s.addr)

	return server.Serve(s.listener)
}

// Handler will setup a router that implements the http.Handler interface.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	// Middleware for logging requests.
	r.Use(middleware.Logger)

	r.Route("/health", newHealthHandler().Routes)
	r.Route("/api/v1/characters", newCharacterHandler(s.encoder, s.characterService).Routes)

	// Deprecated handler, supported for consumers who rely on it.
	r.Route("/retrieving/v1/character", newCharacterHandler(s.encoder, s.characterService).Routes)

	return r
}

// NewServer returns a new server with all dependencies.
func NewServer(addr string, characterService characterService) *Server {
	return &Server{
		addr:             addr,
		encoder:          newEncoder(),
		characterService: characterService,
	}
}
