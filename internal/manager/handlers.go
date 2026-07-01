package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kraljstone/k8n/internal/api"
	"github.com/Kraljstone/k8n/internal/model"
	"github.com/Kraljstone/k8n/internal/runner"
	"github.com/google/uuid"
)

func handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	}
}

func (s *Server) handleTaskCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid task payload"})
			return
		}

		task := model.Task{
			ID:     uuid.New(),
			Name:   req.Name,
			State:  model.StatePending,
			Image:  req.Image,
			Memory: req.Memory,
			Disk:   req.Disk,
		}

		s.store.CreateTask(task)
		fmt.Printf("[Manager] Stored Task: %s (ID: %s)\n", task.Name, task.ID)

		// Try to assign to a worker; fall back to local execution.
		workers := s.store.ListWorkers()
		assigned, err := s.scheduler.Assign(task, workers)
		if err == nil && assigned != nil {
			// TODO: Dispatch task to assigned worker via HTTP.
			fmt.Printf("[Manager] Assigned task %s to worker %s\n", task.Name, assigned.Name)
		} else {
			// Run locally.
			go func(t model.Task) {
				t.State = model.StateRunning
				s.store.UpdateTask(t)
				runner.ExecuteTask(context.Background(), &t)
				s.store.UpdateTask(t)
			}(task)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(toTaskResponse(task))
	}
}

func (s *Server) handleWorkerRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.RegisterWorkerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid registration payload"})
			return
		}

		id := uuid.New()
		node := model.Node{
			ID:       id,
			Name:     req.ID,
			Address:  req.Address,
			Role:     model.NodeRoleWorker,
			State:    model.NodeStateHealthy,
			LastPing: time.Now(),
		}
		s.store.AddWorker(node)

		fmt.Printf("[Manager] Worker registered: %s (ID: %s)\n", req.ID, id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "registered",
			"worker_id": id.String(),
		})
	}
}

// toTaskResponse converts a model.Task to an API response.
func toTaskResponse(t model.Task) api.TaskResponse {
	resp := api.TaskResponse{
		ID:     t.ID.String(),
		Name:   t.Name,
		State:  t.State.String(),
		Image:  t.Image,
		Memory: t.Memory,
		Disk:   t.Disk,
	}
	if !t.StartTime.IsZero() {
		resp.StartTime = t.StartTime.Format(time.RFC3339)
	}
	if !t.FinishTime.IsZero() {
		resp.FinishTime = t.FinishTime.Format(time.RFC3339)
	}
	return resp
}
