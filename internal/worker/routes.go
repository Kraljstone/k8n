package worker

func (s *Server) routes() {
	s.router.HandleFunc("POST /tasks/run", s.handleTaskRun())
}
