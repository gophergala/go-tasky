package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"

	"github.com/gophergala/go-tasky/examples/workers"
	"github.com/gophergala/go-tasky/tasky"
	"github.com/gorilla/mux"
)

var (
	serverAddr string
)

func init() {
	flag.StringVar(&serverAddr, "addr", ":8888", "Set the address of the server. Defaults to :8888")
	flag.Parse()
}

func register() {
	cp := &workers.CopyFile{}

	tw, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("cp: ", tw)
	fmt.Println("info: ", string(tw.Usage()))
	spew.Dump("Details", tw.Details())

	i := &workers.Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("i: ", tw2)

	// s := &workers.Sleeper{}

	// tw3, err := tasky.NewWorker(s)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Println("s: ", tw3)
	// return
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
