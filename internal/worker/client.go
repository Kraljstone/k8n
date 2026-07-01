package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kraljstone/k8n/internal/api"
)

// Client handles communication with the manager.
type Client struct {
	managerAddr string
	httpClient  *http.Client
}

// NewClient creates a new manager client.
func NewClient(managerAddr string) *Client {
	return &Client{
		managerAddr: managerAddr,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Register sends a registration request to the manager.
func (c *Client) Register(id, address string) error {
	payload := api.RegisterWorkerRequest{
		ID:      id,
		Address: address,
		State:   "healthy",
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal registration: %w", err)
	}

	url := fmt.Sprintf("%s/workers/register", c.managerAddr)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("manager rejected registration with code: %d", resp.StatusCode)
	}

	return nil
}
