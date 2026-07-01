package model

import (
	"time"

	"github.com/google/uuid"
)

// NodeRole identifies the role of a node in the cluster.
type NodeRole string

const (
	NodeRoleWorker  NodeRole = "worker"
	NodeRoleManager NodeRole = "manager"
)

// NodeState represents the health state of a node.
type NodeState string

const (
	NodeStateHealthy   NodeState = "healthy"
	NodeStateUnhealthy NodeState = "unhealthy"
)

// Node represents a machine (physical or virtual) in the cluster.
// It consolidates both hardware specs and runtime state into one type.
type Node struct {
	ID       uuid.UUID
	Name     string
	Address  string
	Role     NodeRole
	State    NodeState
	Cores    int
	Memory   int
	Disk     int
	LastPing time.Time
}
