package tasky

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ListWorkers(rw http.ResponseWriter, r *http.Request) {
	jsonStr, _ := listWorkers()

	fmt.Fprintf(rw, "%s\n", jsonStr)
}

func RegisterTaskyHandlers(r *mux.Router) {
	ws := r.Path("/tasky/v1/workers").Subrouter()
	ws.Methods("GET").HandlerFunc(ListWorkers)
}
