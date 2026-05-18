package manager

func (s *Server) routes() {
	s.router.HandleFunc("GET /health", handleHealthCheck())
	s.router.HandleFunc("POST /tasks", s.handleTaskCreate())
}
