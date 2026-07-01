package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// State represents the lifecycle stage of a task.
type State int

const (
	StatePending State = iota
	StateScheduled
	StateRunning
	StateCompleted
	StateFailed
)

func (s State) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateScheduled:
		return "scheduled"
	case StateRunning:
		return "running"
	case StateCompleted:
		return "completed"
	case StateFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// Task represents a unit of work to be executed on the cluster.
type Task struct {
	ID         uuid.UUID
	Name       string
	State      State
	Image      string // The command or image to run
	Memory     int    // Required RAM in MB
	Disk       int    // Required disk in GB
	StartTime  time.Time
	FinishTime time.Time
}
