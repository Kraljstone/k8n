package worker

import (
	"net/http"
)

// Server is the HTTP handler for the worker node.
type Server struct {
	router *http.ServeMux
}

// NewServer creates a new worker server with routes configured.
func NewServer() http.Handler {
	s := &Server{
		router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
