package manager

import (
	"net/http"

	"github.com/Kraljstone/k8n/internal/types"
	"github.com/google/uuid"
)

// 2. THE CONSTRUCTOR: NewServer now initializes the struct, allocates memory
// for the maps, configures the routes, and returns it.
func NewServer() http.Handler {
	server := &Server{
		router:  http.NewServeMux(),
		tasks:   make(map[uuid.UUID]types.Task),
		workers: make(map[uuid.UUID]types.Node),
	}

	// Connect our routes
	server.routes()

	return server
}

// 3. THE UNIFORM: This method makes our custom 'server' struct satisfy
// Go's standard http.Handler interface so main.go can read it.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
