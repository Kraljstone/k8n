package manager

import (
	"fmt"

	"github.com/Kraljstone/k8n/internal/model"
)

// Scheduler handles assigning tasks to worker nodes.
type Scheduler struct{}

// NewScheduler creates a new Scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// Assign selects a worker node to run the given task.
// Returns the assigned worker, or nil with an error if no suitable worker is available.
func (s *Scheduler) Assign(task model.Task, workers []model.Node) (*model.Node, error) {
	for _, w := range workers {
		if w.State == model.NodeStateHealthy {
			return &w, nil
		}
	}
	return nil, fmt.Errorf("no healthy workers available")
}
