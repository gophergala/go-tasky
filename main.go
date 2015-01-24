package main

import (
	"github.com/gophergala/go-tasky/tasky"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	apiBase string
)

func init() {
	apiBase = "/tasky/v1"
}

func main() {
	println("go-tasky Started!")

	workerStore := tasky.Workers{
		Store: map[string]*tasky.Worker{},
	}

	unameWorker := tasky.Worker{
		ID:          "",
		Name:        "uname",
		Description: "fetch operating system name",
	}
	workerStore.CreateWorker(&unameWorker)

	ifconfigWorker := tasky.Worker{
		ID:          "",
		Name:        "ifconfig",
		Description: "fetch network configuration",
	}
	workerStore.CreateWorker(&ifconfigWorker)

	mainRouter := mux.NewRouter()

	//handles /tasky/v1 routes.  Create new sub-routers off of this
	taskySubRtr := mainRouter.PathPrefix(apiBase).Subrouter()

	//handles /tasky/v1/workers routes
	workersSubRtr := taskySubRtr.PathPrefix("/workers").Subrouter()

	workersSubRtr.HandleFunc("/", workerStore.Index).Methods("GET")

	http.ListenAndServe(":4444", mainRouter)

}
