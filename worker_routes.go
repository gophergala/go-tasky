package tasky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	// "log"
	"net/http"
	// "strings"
)

type WorkerRoutes struct {
}

func (wr *WorkerRoutes) Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index route hit")

	type returnWorker struct {
		Id   string `json:"id"`
		Info string `json:"info"`
	}

	type ReturnWs struct {
		Workers []returnWorker `json:"workers"`
	}
	workersList, _ := listWorkers()
	um := ws{}
	json.Unmarshal(workersList, &um)

	returnWorkers := ReturnWs{Workers: make([]returnWorker, len(um.Workers))}
	for k, _ := range um.Workers {
		workerItem := um.Workers[k]
		buff := new(bytes.Buffer)
		if err := json.Compact(buff, workerItem.Info); err != nil {
			fmt.Println(err)
		}
		returnWorkers.Workers[k].Id = workerItem.Id
		returnWorkers.Workers[k].Info = buff.String()
	}

	w.Header().Set("Content-Type", "application/json")
	RespondJSON(w, req, returnWorkers)
}
