package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophergala/go-tasky/tasky"
	"github.com/gorilla/mux"
)

func register() {
	cp := &CopyFile{}

	tw, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("cp: ", tw)
	fmt.Println("info: ", string(tw.Info()))

	i := &Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("i: ", tw2)

	s := &Sleeper{}

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
