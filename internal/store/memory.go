package store

import (
	"sync"

	"github.com/Kraljstone/k8n/internal/model"
	"github.com/google/uuid"
)

// MemoryStore provides an in-memory, concurrency-safe store for tasks and workers.
type MemoryStore struct {
	mu      sync.RWMutex
	tasks   map[uuid.UUID]model.Task
	workers map[uuid.UUID]model.Node
}

// NewMemoryStore creates a new empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:   make(map[uuid.UUID]model.Task),
		workers: make(map[uuid.UUID]model.Node),
	}
}

// ---------- Task operations ----------

// CreateTask inserts a new task into the store.
func (s *MemoryStore) CreateTask(t model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

// GetTask retrieves a task by ID. Returns false if not found.
func (s *MemoryStore) GetTask(id uuid.UUID) (model.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

// UpdateTask replaces an existing task in the store.
func (s *MemoryStore) UpdateTask(t model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

// ListTasks returns a snapshot of all tasks.
func (s *MemoryStore) ListTasks() []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

// ---------- Worker operations ----------

// AddWorker registers a new worker node.
func (s *MemoryStore) AddWorker(w model.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers[w.ID] = w
}

// GetWorker retrieves a worker by ID. Returns false if not found.
func (s *MemoryStore) GetWorker(id uuid.UUID) (model.Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.workers[id]
	return w, ok
}

// UpdateWorker updates an existing worker's state.
func (s *MemoryStore) UpdateWorker(w model.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers[w.ID] = w
}

// ListWorkers returns a snapshot of all registered workers.
func (s *MemoryStore) ListWorkers() []model.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Node, 0, len(s.workers))
	for _, w := range s.workers {
		result = append(result, w)
	}
	return result
}

// NextWorker returns the first healthy worker from the store.
// Returns nil if no healthy workers are available.
func (s *MemoryStore) NextWorker() *model.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.workers {
		if w.State == model.NodeStateHealthy {
			clone := w
			return &clone
		}
	}
	return nil
}
