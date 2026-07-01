package api

// CreateTaskRequest is the JSON payload for creating a new task.
type CreateTaskRequest struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Memory int    `json:"memory"`
	Disk   int    `json:"disk"`
}

// TaskResponse is the JSON representation of a task returned by the API.
type TaskResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Image      string `json:"image"`
	Memory     int    `json:"memory"`
	Disk       int    `json:"disk"`
	StartTime  string `json:"start_time,omitempty"`
	FinishTime string `json:"finish_time,omitempty"`
}
