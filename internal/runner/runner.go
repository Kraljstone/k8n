package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Kraljstone/k8n/internal/types"
)

// ExecuteTask takes a task, runs it on the local OS, and tracks its lifecycle.
func ExecuteTask(ctx context.Context, t *types.Task) error {
	// 1. Update task state to Running
	t.State = types.Running
	t.StartTime = time.Now()
	fmt.Printf("[Runner] Starting task: %s\n", t.Name)

	// 2. Parse the command string (e.g., "echo hello" becomes ["echo", "hello"])
	// For now, we use the "Image" field to hold our shell command string
	parts := strings.Fields(t.Image)
	if len(parts) == 0 {
		t.State = types.Failed
		t.FinishTime = time.Now()
		return fmt.Errorf("empty command provided")
	}

	cmdName := parts[0]  // e.g., "echo" or "sleep"
	cmdArgs := parts[1:] // e.g., ["hello"] or ["2"]

	// 3. Prepare the system command and link it to our lifecycle context
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)

	// 4. STREAMING: Pipe the command's live output directly to our terminal screen
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 5. Run the command and wait for it to finish
	err := cmd.Run()
	t.FinishTime = time.Now()

	if err != nil {
		t.State = types.Failed
		fmt.Printf("[Runner] Task %s failed: %v\n", t.Name, err)
		return err
	}

	// 6. If it succeeds
	t.State = types.Completed
	fmt.Printf("[Runner] Task %s completed successfully!\n", t.Name)
	return nil
}
