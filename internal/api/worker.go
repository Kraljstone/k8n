package api

// RegisterWorkerRequest is the JSON payload sent by a worker on startup.
type RegisterWorkerRequest struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	State   string `json:"state"`
}

// WorkerResponse is the JSON representation of a worker returned by the API.
type WorkerResponse struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	State    string `json:"state"`
	LastPing string `json:"last_ping,omitempty"`
}
