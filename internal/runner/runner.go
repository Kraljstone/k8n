package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Kraljstone/k8n/internal/model"
)

// ExecuteTask takes a task, runs it on the local OS, and tracks its lifecycle.
func ExecuteTask(ctx context.Context, t *model.Task) error {
	t.State = model.StateRunning
	t.StartTime = time.Now()
	fmt.Printf("[Runner] Starting task: %s\n", t.Name)

	parts := strings.Fields(t.Image)
	if len(parts) == 0 {
		t.State = model.StateFailed
		t.FinishTime = time.Now()
		return fmt.Errorf("empty command provided")
	}

	cmdName := parts[0]
	cmdArgs := parts[1:]

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	t.FinishTime = time.Now()

	if err != nil {
		t.State = model.StateFailed
		fmt.Printf("[Runner] Task %s failed: %v\n", t.Name, err)
		return err
	}

	t.State = model.StateCompleted
	fmt.Printf("[Runner] Task %s completed successfully!\n", t.Name)
	return nil
}
