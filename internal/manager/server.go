package manager

import (
	"net/http"

	"github.com/Kraljstone/k8n/internal/store"
)

// Server is the HTTP handler for the manager control plane.
type Server struct {
	router    *http.ServeMux
	store     *store.MemoryStore
	scheduler *Scheduler
}

// NewServer creates a new manager server with its dependencies wired up.
func NewServer() http.Handler {
	s := &Server{
		router:    http.NewServeMux(),
		store:     store.NewMemoryStore(),
		scheduler: NewScheduler(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
