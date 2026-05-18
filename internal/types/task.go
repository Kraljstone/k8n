package types

import (
	"time"

	"github.com/google/uuid"
)

// State represents the current lifecycle stage of a task
type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

// Task is the "Container" for the work we want to do
type Task struct {
	ID         uuid.UUID
	Name       string
	State      State
	Image      string // What software to run (e.g., "python-script")
	Memory     int    // Required RAM in MB
	Disk       int    // Required Disk in GB
	StartTime  time.Time
	FinishTime time.Time
}
