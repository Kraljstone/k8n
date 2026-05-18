package manager

import (
	"net/http"
	"sync"

	"github.com/Kraljstone/k8n/internal/types"
	"github.com/google/uuid"
)

type Server struct {
	router  *http.ServeMux
	mu      sync.RWMutex             // Keeps our maps safe from concurrent read/write crashes
	tasks   map[uuid.UUID]types.Task // Our in-memory Task database
	workers map[uuid.UUID]types.Node // Our in-memory Worker database
}
