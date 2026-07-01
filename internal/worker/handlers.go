package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kraljstone/k8n/internal/model"
	"github.com/Kraljstone/k8n/internal/runner"
)

func (s *Server) handleTaskRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			fmt.Printf("[Worker] Failed to decode task payload: %v\n", err)
			http.Error(w, "Bad request payload", http.StatusBadRequest)
			return
		}

		fmt.Printf("[Worker] Received assignment: Task %s (ID: %s)\n", task.Name, task.ID)

		// Use context.Background() so the task isn't cancelled when the HTTP response is sent.
		go func(t model.Task) {
			err := runner.ExecuteTask(context.Background(), &t)
			if err != nil {
				fmt.Printf("[Worker] Task execution failed for %s: %v\n", t.Name, err)
				return
			}
			fmt.Printf("[Worker] Task %s finished execution on this node.\n", t.Name)
		}(task)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"running"}`))
	}
}
