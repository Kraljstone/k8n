package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kraljstone/k8n/internal/runner"
	"github.com/Kraljstone/k8n/internal/types"
	"github.com/google/uuid"
)

func handleHealthCheck() http.HandlerFunc {
	// Any setup code here runs ONCE when the server boots up.

	return func(w http.ResponseWriter, r *http.Request) {
		// This code runs on EVERY network request.
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "healthy"}`)
	}
}

func (s *Server) handleTaskCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task types.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "invalid task payload"}`)
			return
		}

		task.ID = uuid.New()
		task.State = types.Pending
		task.StartTime = time.Now()

		// Save to our server memory map
		s.mu.Lock()
		s.tasks[task.ID] = task
		s.mu.Unlock()

		fmt.Printf("[Manager] Stored Task: %s (ID: %s)\n", task.Name, task.ID)

		// 2. THE LINK: Run the task in the background using 'go'
		// We use context.Background() here so the task keeps running
		// even after this specific HTTP request finishes and closes.
		go func(t types.Task) {
			err := runner.ExecuteTask(context.Background(), &t)

			// 3. Update the task state in our master memory map when it finishes!
			s.mu.Lock()
			s.tasks[t.ID] = t
			s.mu.Unlock()

			if err != nil {
				fmt.Printf("[Manager] Background task %s failed\n", t.Name)
			}
		}(task)

		// Respond to the API client immediately while the task runs in the background
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	}
}
