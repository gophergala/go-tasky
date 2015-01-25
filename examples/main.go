package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gophergala/go-tasky"
	"github.com/gophergala/go-tasky/examples/workers"
	"github.com/gorilla/mux"
)

var (
	serverAddr string
)

func init() {
	flag.StringVar(&serverAddr, "addr", ":8888", "Set the address of the server. Defaults to :8888")
	flag.Parse()
}

//register all of your custom workers here.
func register() {
	cp := &workers.CopyFile{}
	taskyCopy, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("copyFile Details: %+v\n\n", taskyCopy.Details())

	ifconf := &workers.Ifconfig{}
	taskyifconfig, err := tasky.NewWorker(ifconf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("ifconfig worker: %+v\n\n", taskyifconfig.Details())

	sleeper := &workers.Sleeper{}
	taskySleeper, err := tasky.NewWorker(sleeper)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("sleeper worker: %+v\n\n", taskySleeper.Details())
}

func main() {
	register()

	r := mux.NewRouter()

	tasky.RegisterTaskyHandlers(r)
	log.Println("Starting Tasky server at", serverAddr)

	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		log.Println("Error Starting Tasky Server: ", err.Error())
	}

}
