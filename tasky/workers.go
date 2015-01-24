package tasky

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

//Worker is an interface for defining a worker type.  All custom workers must implement this interface.
type Workers struct {
	sync.RWMutex
	Store map[string]*Worker
}

type Worker struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (workers *Workers) CreateWorker(worker *Worker) {
	workers.Lock()
	id := fmt.Sprintf("%d", len(workers.Store)) //replace this with database backing id, temp for now
	worker.ID = id
	workers.Store[id] = worker
	workers.Unlock()
}

func (workers *Workers) Index(w http.ResponseWriter, req *http.Request) {
	if len(workers.Store) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	RespondJSON(w, req, workers.Store)

}

func RespondJSON(w http.ResponseWriter, req *http.Request, v interface{}, code ...int) error {
	if code != nil {
		w.WriteHeader(code[0])
	}
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		//encoding failed
		panic(err)
	}

	return err
}
