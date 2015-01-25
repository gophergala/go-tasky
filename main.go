package main

import (
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/gophergala/go-tasky/tasky"
	"github.com/gophergala/go-tasky/workers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	apiBase     string
	workerStore tasky.Resources
)

func init() {
	apiBase = "/tasky/v1"
	workerStore = tasky.Resources{
		Store: map[string]*tasky.Worker{},
	}
}

func main() {
	println("go-tasky Started!")
	workerRoutes := &tasky.WorkerRoutes{}
	cpWorker := &workers.CopyFile{}
	tw, err := tasky.NewWorker(cpWorker)
	if err != nil {
		log.Panicf("Error Creating new worker: %v", err)
	}

	ifconfWorker := &workers.Ifconfig{}
	tw2, err := tasky.NewWorker(ifconfWorker)
	if err != nil {
		log.Panicf("Error Creating new Worker: %v", err)
	}

	workerStore.RegisterWorker("Copy File", tw)
	workerStore.RegisterWorker("ifconfig", tw2)

	mainRouter := mux.NewRouter()
	// svc := tasky.NewService(apiBase)

	//handles /tasky/v1 routes.  Create new sub-routers off of this
	taskySubRtr := mainRouter.PathPrefix(apiBase).Subrouter()

	//handles /tasky/v1/workers routes
	workersSubRtr := taskySubRtr.PathPrefix("/workers").Subrouter()

	workersSubRtr.HandleFunc("/", workerRoutes.Index).Methods("GET")

	http.ListenAndServe(":4444", mainRouter)
}
