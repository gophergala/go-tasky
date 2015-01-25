package tasky

import (
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

type WorkerRoutes struct {
}

func (wr *WorkerRoutes) Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index route hit")
	workersList, err := listWorkers()
	// workerStore.Lock()
	// workers := make([]Worker, len(workerStore.Store))
	// i := 0
	// for _, worker := range workerStore.Store {
	// 	spew.Dump("Index Worker Parse: ", *worker)
	// 	workers[i] = *worker
	// 	// workers[i].Info = *worker.Info()
	// 	i++
	// }
	// workerStore.Unlock()

	// spew.Dump("Index List store: ", workers[0].Info())
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(workersList)
	if err != nil {
		log.Println("Index Header Write Error: ", err)
	}
}
