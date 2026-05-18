package types

import "github.com/google/uuid"

// Node represents a physical or virtual machine in our cluster
type Node struct {
	ID     uuid.UUID
	Name   string
	IP     string
	Cores  int
	Memory int
	Disk   int
	Role   string // "Worker" or "Manager"
}
