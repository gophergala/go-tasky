package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophergala/go-tasky/examples/workers"
	"github.com/gophergala/go-tasky/tasky"
	"github.com/gorilla/mux"
)

func register() {
	cp := &workers.CopyFile{}

	tw, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("cp: ", tw)
	fmt.Println("info: ", string(tw.Usage()))

	i := &workers.Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("i: ", tw2)

	s := &workers.Sleeper{}

	tw3, err := tasky.NewWorker(s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("s: ", tw3)
}

func main() {
	register()

	r := mux.NewRouter()

	tasky.RegisterTaskyHandlers(r)

	http.ListenAndServe(":12345", r)
}
